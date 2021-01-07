import random
import logging
import json
import gym
import math
import itertools
from gym import spaces
import pandas as pd
import numpy as np

MAX_STEPS = 20000

INITIAL_ACCOUNT_BALANCE = 10000

RESET_STEP_RANDOM_LIMIT = 100

NORMALIZER = 1000000

MAX_QTY = 1000

# TODO Add help
class StockTradingEnv(gym.Env):
    """A stock trading environment for OpenAI gym"""
    metadata = {'render.modes': ['human']}

    def __init__(self, df):
        super(StockTradingEnv, self).__init__()

        # Save dataframe
        self.df = df
        # Declare var for stock qty
        self.stock_qty = 0
        # Declare var for stock cost for reward calculation
        self.stock_cost = 0
        # Declare balance
        self.balance = INITIAL_ACCOUNT_BALANCE
        # Overall worth balance + stocks
        self.overall_worth = self.balance

        # Set the current step to a random point within the data frame
        self.current_step = random.randint(0, RESET_STEP_RANDOM_LIMIT)

        # Format: <Qty>
        self.action_space = spaces.Box(
            low=-MAX_QTY, high=MAX_QTY, shape=(1,), dtype=np.float32)

        # Current price + metrics + stock balance + balance (10 + 1 + 1)
        #
        # [
        # "open", "close", "high", "low", "volume",
        # "MACD", "MACDh", "MACDs", "RSI", "MFI",
        # "qty"
        # "available_balance"
        # ]
        self.observation_space = spaces.Box(
            low=0, high=1, shape=(12,), dtype=np.float16)

    def _next_observation(self):
        # Get current data
        data = self.df.iloc[self.current_step]

        # Normalize
        data /= NORMALIZER

        # Append current stock qty
        data = np.append(data, self.stock_qty)

        # Append current balance
        data = np.append(data, self.balance)

        return data

    def _choose_stock(self, index):
        if index < 1:
            return 0
        elif index < 2:
            return 1
        else:
            return 2

    def step(self, action):
        # Parse: <Qty>
        qty = 0 if math.isnan(action) else int(action)

        # logging.debug("STEP action=%s action_type=%d stock_index=%d qty=%d" % (action, action_type, stock_index, qty))

        # Set the current price to a random price within the time step
        # TODO Can be optimized with double index in dataframe
        current_price = random.uniform(
            self.df.iloc[self.current_step]["open"],
            self.df.iloc[self.current_step]["close"])

        # Available actions:
        # < 0 - Sell
        # 0 - Hold
        # > 0 - Buy
        if qty < 0:
            # SELL
            # logging.debug("SELL action")

            # Check that can sell stock with `stock_index` and amount `qty`
            current_qty = self.stock_qty

            if current_qty == 0:
                # logging.debug(
                #     "Try to sell stock #%d, but have nothing" % stock_index)

                # Calc reward
                reward = 0.0

            elif (current_qty - qty) <= 0:
                # Calc sold price based on all current qty
                sold_price = self.stock_qty * current_price
                # Add money to balance
                self.balance += sold_price

                logging.debug("Sold %d stock with price %f. Current balance %f." % (
                    self.stock_qty, current_price, self.balance))

                # Calc reward as diff between prices
                sold_profit = sold_price - self.stock_cost
                reward = sold_profit / self.stock_cost if sold_profit > 0 else 0

                # Set current qty to zero
                self.stock_qty = 0
                # Set cost to zero
                self.stock_cost = 0

            else:
                # Calc sold price based on sold qty
                sold_price = qty * current_price
                # Add money to balance
                self.balance += sold_price

                logging.debug("Sold %d stock with price %f. Current balance %f." % (
                    qty, current_price, self.balance))

                # Calc reward as diff between prices
                sold_profit = sold_price - (self.stock_cost / self.stock_qty * qty)
                reward = sold_profit / self.stock_cost if sold_profit > 0 else 0

                # Decrease current qty
                self.stock_qty -= qty
                # Decrease cost as we can
                self.stock_cost -= qty * current_price

        elif qty > 0:
            # BUY
            # logging.debug("BUY action")

            # Check that we can buy
            if (self.balance < current_price):
                # logging.debug("We can't buy stock #%d, because price is %f and we have %f." % (
                #     stock_index, current_price, self.balance))

                # Calc reward
                reward = 0.0

            elif (self.balance < (current_price * qty)):
                # Calc max amount
                max_qty = self.balance // current_price
                # Decrease balance
                self.balance -= max_qty * current_price
                # Increase amount in balance
                self.stock_qty += max_qty
                # Increase cost
                self.stock_cost += max_qty * current_price

                logging.debug("Bought %d stocks, with price %f. Balance: %f, qty: %d" % (
                    max_qty, current_price, self.balance, self.stock_qty))

                # Calc reward
                reward = 0.0

            else:
                # Decrease balance
                self.balance -= qty * current_price
                # Increase amount in balance
                self.stock_qty += qty
                # Increase cost
                self.stock_cost += qty * current_price

                logging.debug("Bought %d stocks, with price %f. Balance: %f, qty: %d" % (
                    qty, current_price, self.balance, self.stock_qty))

                # Calc reward
                reward = 0.0

        else:
            # HOLD
            # logging.debug("HOLD action")

            # Calc reward
            reward = 0.0

        # Calc overall worth
        self.overall_worth = self.balance + self.stock_qty * current_price
        
        # Increase step counter
        self.current_step += 1

        # Get next observations
        obs = self._next_observation()

        # Checks for done flag
        # Step is last, or overall worth is negative
        done = (self.current_step >= self.df.shape[0]) | (
            self.overall_worth <= 0)
        
        return obs, reward, done, {"step":self.current_step}

    def reset(self):
        # Declare array for stock qty
        self.stock_qty = 0
        self.stock_cost = 0
        # Declare balance
        self.balance = INITIAL_ACCOUNT_BALANCE
        # Overall worth balance + stocks
        self.overall_worth = self.balance

        # Set the current step to a random point within the data frame
        self.current_step = random.randint(0, RESET_STEP_RANDOM_LIMIT)

        return self._next_observation()

    def render(self, mode='human', close=False):
        # Render the environment to the screen
        # TODO Add draw

        if self.current_step % 2000 == 0:
            profit = self.overall_worth - INITIAL_ACCOUNT_BALANCE
            logging.info("\n\nStep: %d Balance: %.2f Overall: %.2f Account: %s Profit: %f\n\n" %\
                (self.current_step, self.balance, self.overall_worth, self.stock_qty, profit))
