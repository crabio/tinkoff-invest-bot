import random
import logging

from enum import Enum

class RsiSignal(Enum):
    NaN = 1
    OverBought = 2
    OverSold = 3

class MacdSignal(Enum):
    NaN = 1
    GoUp = 2
    GoDown = 3

class TradingBot:
    def __init__(self,
                initial_budget=1000,
                broker_commission=0.0005,
                rsi_oversold_threshold=30,
                rsi_overbought_threshold=70,
                macd_threshold=0.001,
                stop_loss=0.05):
        # Budget
        self.initial_budget = initial_budget
        # Commission
        self.broker_commission = broker_commission
        # RSI thresholds
        self._rsi_oversold_threshold = rsi_oversold_threshold
        self._rsi_overbought_threshold = rsi_overbought_threshold
        # MACD threshold
        self._macd_threshold = macd_threshold

        # Stop loss
        self._stop_loss = stop_loss
        
        self.reset()
    
    def reset(self):
        # Budget
        self.budget = self.initial_budget
        # Bought price
        self.bought_price = 0
        self.bought_count = 0
        # init flags
        # RSI
        self._rsi_k_prev_value = None
        self._rsi_d_prev_value = None
        self._rsi_signal = RsiSignal.NaN
        # MACD
        self._macd_prev_value = None
        self._macd_signal = MacdSignal.NaN
        
        # Profit
        self.profit = 0
        
        self.current_step = 0

    def _set_rsi_signal(self, rsi_k_value, rsi_d_value):

        # Check that we have prev value
        if (self._rsi_k_prev_value is not None) &\
            (self._rsi_d_prev_value is not None):
            # RSI flag conditions
            if (self._rsi_k_prev_value <= self._rsi_oversold_threshold) &\
                (self._rsi_k_prev_value <= self._rsi_d_prev_value) &\
                (rsi_k_value > rsi_d_value):
                self._rsi_signal = RsiSignal.OverSold
            if (self._rsi_k_prev_value >= self._rsi_overbought_threshold) &\
                (self._rsi_k_prev_value >= self._rsi_d_prev_value) &\
                (rsi_k_value < rsi_d_value):
                self._rsi_signal = RsiSignal.OverBought

        # Save value as previous
        self._rsi_k_prev_value = rsi_k_value
        self._rsi_d_prev_value = rsi_d_value


    def _set_macd_signal(self, macd_value):
        # Check that we have prev value for diff
        if self._macd_prev_value is not None:
            # MACD flag conditions
            # GoUp - \_/
            # GoDown - /─\
            if (self._macd_prev_value <= 0) & (macd_value > self._macd_threshold):
                # GoUp
                # Buy
                self._macd_signal = MacdSignal.GoUp
            elif (self._macd_prev_value >= 0) & (macd_value < -self._macd_threshold):
                # GoDown
                # Sell
                self._macd_signal = MacdSignal.GoDown

        # Save current value as previous
        self._macd_prev_value = macd_value


    def _buy(self, current_price):
        # Calc max amount
        max_count = self.budget // current_price
        # Check max available count
        if (max_count > 0):
            logging.debug("Step: %d Buy %d with price=%f" % (self.current_step, max_count, current_price))
        
            # We can buy
            self.budget -= (current_price * max_count) * (1 + self.broker_commission)
            # Increrase count
            self.bought_price += current_price * max_count
            self.bought_count += max_count

            return True

        return False
    
    def _sell(self, current_price):
        # Check that we can sell
        if self.bought_count != 0:
            # Get profit
            self.budget += (current_price * self.bought_count) * (1 - self.broker_commission)
            # Calc profit
            profit = (current_price * self.bought_count - self.bought_price) / self.bought_price
            self.profit = (self.budget - self.initial_budget) / self.initial_budget
            
            logging.debug("Step: %d Sell %d with price=%f with profit %.2f%% and overall %.2f%%" % \
                (self.current_step, self.bought_count, current_price, profit * 100, self.profit * 100))
            
            # Reset balance
            self.bought_price = 0
            self.bought_count = 0

            return True

        return False
        

    def _check_stop_loss(self, current_price):
        # Calc profit
        profit = (current_price * self.bought_count - self.bought_price) / self.bought_price

        if profit <= -self._stop_loss:
            # Sell
            # Get profit
            self.budget += current_price * self.bought_count
            # Calc profit
            self.profit = (self.budget - self.initial_budget) / self.initial_budget

            logging.debug("Step: %d Stop loss %.2f%% from %.2f to %.2f with profit %.2f%% overall %.2f%%" % \
                (self.current_step, -self._stop_loss * 100, self.bought_price / self.bought_count, current_price, profit * 100, self.profit * 100))

            # Reset balance
            self.bought_price = 0
            self.bought_count = 0

            return True

        return False



    def process(self,
                open_price,
                close_price,
                rsi_k_value,
                rsi_d_value,
                macd_h_value):
        """
        Returns
        0 - Hold
        1 - Buy
        2 - Sell
        """

        self.current_step += 1
        # current_price = random.uniform(open_price, close_price)
        current_price = close_price
        
        # Set RSI signal
        self._set_rsi_signal(rsi_k_value, rsi_d_value)

        # Set MACD signal
        self._set_macd_signal(macd_h_value)

        # Buy only if MACD and RSI is ok in this point
        #  & (self._macd_signal == MacdSignal.GoUp)
        if (self._rsi_signal == RsiSignal.OverSold):
            # Buy
            result = self._buy(current_price)

            # Return result of processing
            if (result):
                # Successfull buy
                return 1
            else:
                # Hold
                return 0

        # Buy if MACD or RSI is ok in this point
        #  & (self._macd_signal == MacdSignal.GoDown))
        elif (self._rsi_signal == RsiSignal.OverBought):
            # Sell
            result = self._sell(current_price)

            # Return result of processing
            if (result):
                # Successfull sell
                return 2
            else:
                # Hold
                return 0
        
        if self.bought_count > 0:
            # Check stop-loss
            result = self._check_stop_loss(current_price)

            # Return result of processing
            if (result):
                # Successfull sell
                return 2
            else:
                # Hold
                return 0

        # Hold
        return 0

        
            