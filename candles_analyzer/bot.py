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
                stop_loss=0.05):
        # Budget
        self.initial_budget = initial_budget
        # Commission
        self.broker_commission = broker_commission
        # RSI thresholds
        self._rsi_oversold_threshold = 30
        self._rsi_overbought_threshold = 70

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
        self._rsi_signal = RsiSignal.NaN
        # MACD
        self._macd_prev_value = None
        self._macd_prev_diff = None
        self._macd_signal = MacdSignal.NaN
        
        # Profit
        self.profit = 0
        
        self.current_step = 0

    def _set_rsi_signal(self, rsi_value):
        # RSI flag conditions
        if rsi_value > self._rsi_overbought_threshold:
            self._rsi_signal = RsiSignal.OverBought
        elif rsi_value < self._rsi_oversold_threshold:
            self._rsi_signal = RsiSignal.OverSold

    def _set_macd_signal(self, macd_value, macd_signal_value):
        # Check that we have prev value for diff
        if self._macd_prev_value is not None:
            # Calc MACD diff
            macd_diff = macd_value - self._macd_prev_value

            # Check that we have prev diff for condition
            if self._macd_prev_diff is not None:
                # MACD flag conditions
                # GoUp - \_/
                # GoDown - /─\
                if (self._macd_prev_diff <= 0) & (macd_diff > 0) & (macd_signal_value < 0):
                    # GoUp
                    # Buy
                    self._macd_signal = MacdSignal.GoUp
                elif (self._macd_prev_diff >= 0) & (macd_diff < 0) & (macd_signal_value > 0):
                    # GoDown
                    # Sell
                    self._macd_signal = MacdSignal.GoDown

            # Save current diff value as previous
            self._macd_prev_diff = macd_diff

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
        

    def _check_stop_loss(self, current_price):
        # Calc profit
        profit = (current_price * self.bought_count - self.bought_price) / self.bought_price

        if profit <= -self._stop_loss:
            # Sell
            # Set stop loss as profit
            profit = -self._stop_loss
            current_price = (self.bought_price / self.bought_count) * (1 + profit)

            # Get profit
            self.budget += current_price * self.bought_count
            # Calc profit
            self.profit = (self.budget - self.initial_budget) / self.initial_budget

            logging.debug("Step: %d Stop loss %.2f%% from %.2f to %.2f with profit %.2f%% overall %.2f%%" % \
                (self.current_step, -self._stop_loss, self.bought_price / self.bought_count, current_price, profit * 100, self.profit * 100))

            # Reset balance
            self.bought_price = 0
            self.bought_count = 0



    def process(self, data):
        self.current_step += 1
        current_price = random.uniform(data["open_price"],data["close_price"])
        
        # Set RSI signal
        self._set_rsi_signal(data["RSI"])

        # Set MACD signal
        self._set_macd_signal(data["MACD"], data["MACDs"])

        # Try to find buy signal
        if (self._rsi_signal == RsiSignal.OverSold) & (self._macd_signal == MacdSignal.GoUp):
            # Buy
            self._buy(current_price)
        # Try to find sell signal
        elif (self._rsi_signal == RsiSignal.OverBought) & (self._macd_signal == MacdSignal.GoDown):
            # Sell
            self._sell(current_price)
        
        if self.bought_count > 0:
            # Check stop-loss
            self._check_stop_loss(current_price)
            