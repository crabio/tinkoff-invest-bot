import random
import logging
import json
import gym
import itertools
from gym import spaces
import pandas as pd
import numpy as np

MAX_STEPS = 20000

INITIAL_ACCOUNT_BALANCE = 10000

RESET_STEP_RANDOM_LIMIT = 100

INSTRUMENTS_COUNT = 3

MAX_QTY = 1000000

# TODO Add help
class StockTradingEnv(gym.Env):
    """A stock trading environment for OpenAI gym"""
    metadata = {'render.modes': ['human']}

    def __init__(self, df):
        super(StockTradingEnv, self).__init__()

        # Save dataframe
        self.df = df
        # Declare array for stock qty
        self.stock_qty = np.array([0] * INSTRUMENTS_COUNT)
        # Declare array for stock cost for reward calculation
        self.stock_cost = np.array([0] * INSTRUMENTS_COUNT)
        # Declare balance
        self.balance = INITIAL_ACCOUNT_BALANCE
        # Overall worth balance + stocks
        self.overall_worth = self.balance

        # Format: [<Action>, <StockIndex>, <Qty>]
        #
        # Available actions:
        # 0-1 - Sell
        # 1-2 - Hold
        # 2-3 - Buy
        self.action_space = spaces.Box(
            low=np.array([0, 0, 0]), high=np.array([3, 3, MAX_QTY]), dtype=np.float16)

        # Current price + metrics + stock balance + balance (10 X 3 + 3 + 1)
        #
        # [
        # "1_open", "1_close", "1_high", "1_low", "1_volume",
        # "1_MACD", "1_MACDh", "1_MACDs", "1_RSI", "1_MFI",
        # 2 ...
        # 3 ...
        # "1_qty"
        # 2 ...
        # 3 ...
        # "available_balance"
        # ]
        self.observation_space = spaces.Box(
            low=0, high=1, shape=(34,), dtype=np.float16)

    def _next_observation(self):
        # Get current data
        data = self.df.loc[self.current_step]

        # Append current stock qty
        data = np.append(data, self.stock_qty)

        # Append current balance
        data = np.append(data, self.balance)

        return data

    def step(self, action):
        # Parse: [<Action>, <StockIndex>, <Qty>]
        action_type = action[0]
        stock_index = action[1]
        qty = action[2]

        # Set the current price to a random price within the time step
        # TODO Can be optimized with double index in dataframe
        current_price = random.uniform(
            self.df.loc[self.current_step, "%d_open" % stock_index],
            self.df.loc[self.current_step, "%d_close" % stock_index])

        # Available actions:
        # 0-1 - Sell
        # 1-2 - Hold
        # 2-3 - Buy
        if action_type < 1:
            # SELL
            logging.debug("SELL action")

            # Check that can sell stock with `stock_index` and amount `qty`
            current_qty = self.stock_qty[stock_index]

            if current_qty == 0:
                logging.debug(
                    "Try to sell stock #%d, but have nothing" % stock_index)

                # Calc reward
                reward = 0.0

            elif (current_qty - qty) <= 0:
                logging.debug("We have %d stocks #%d, but try to sell %d. Sell all that we have." % (
                    current_qty, stock_index, qty))

                # Calc sold price based on all current qty
                sold_price = self.stock_qty[stock_index] * current_price
                # Add money to balance
                self.balance += sold_price

                logging.debug("Sold %d stock #%d with price %f. Current balance %f." % (
                    self.stock_qty[stock_index], stock_index, current_price, self.balance))

                # Calc reward as diff between prices
                reward = 1.0 if ((self.stock_cost[stock_index] / self.stock_qty[stock_index]) < current_price) else 0.0

                # Set current qty to zero
                self.stock_qty[stock_index] = 0
                # Set cost to zero
                self.stock_cost[stock_index] = 0
            else:
                logging.debug("We have %d stocks #%d and sell %d. Sell all that we have." % (
                    current_qty, stock_index, qty))

                # Calc sold price based on sold qty
                sold_price = qty * current_price
                # Add money to balance
                self.balance += sold_price

                logging.debug("Sold %d stock #%d with price %f. Current balance %f." % (
                    qty, stock_index, current_price, self.balance))

                # Calc reward as diff between prices
                reward = 1.0 if ((self.stock_cost[stock_index] / self.stock_qty[stock_index]) < current_price) else 0.0

                # Decrease current qty
                self.stock_qty[stock_index] -= qty
                # Decrease cost as we can
                self.stock_cost[stock_index] -= qty * current_price

        elif action_type < 2:
            # HOLD
            logging.debug("HOLD action")

            # Calc reward
            reward = 0.0

        else:
            # BUY
            logging.debug("BUY action")

            # Check that we can buy
            if (self.balance < current_price):
                logging.debug("We can't buy stock #%d, because price is %f and we have %f." % (
                    stock_index, current_price, self.balance))
            elif (self.balance < (current_price * qty)):
                logging.debug("We can't buy all %d stocks #%d, because price is %f and we have %f. Buy as we can." % (
                    qty, stock_index, current_price, self.balance))

                # Calc max amount
                max_qty = self.balance // current_price
                # Decrease balance
                self.balance -= max_qty * current_price
                # Increase amount in balance
                self.stock_qty[stock_index] += max_qty
                # Increase cost
                self.stock_cost[stock_index] += max_qty * current_price

                logging.debug("Bought %d stocks #%d, with price %f. Balance: %f, qty: %d" % (
                    max_qty, stock_index, current_price, self.balance, self.stock_qty[stock_index]))

            else:
                logging.debug("Buy %d stocks #%d, with price %f. We have %f." % (
                    qty, stock_index, current_price, self.balance))

                # Decrease balance
                self.balance -= qty * current_price
                # Increase amount in balance
                self.stock_qty[stock_index] += qty
                # Increase cost
                self.stock_cost[stock_index] += qty * current_price

                logging.debug("Bought %d stocks #%d, with price %f. Balance: %f, qty: %d" % (
                    qty, stock_index, current_price, self.balance, self.stock_qty[stock_index]))

            # Calc reward
            reward = 0.0


        # Calc overall worth
        self.overall_worth = self.balance
        # Add stocks cost
        for i in range(INSTRUMENTS_COUNT):
            # Set the current price to a random price within the time step
            # TODO Can be optimized with double index in dataframe
            current_price = random.uniform(
                self.df.loc[self.current_step, "%d_open" % i],
                self.df.loc[self.current_step, "%d_close" % i])
            # Add worth
            self.overall_worth += self.stock_qty[i] * current_price

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
        self.stock_qty = np.array([0] * INSTRUMENTS_COUNT)
        # Declare balance
        self.balance = INITIAL_ACCOUNT_BALANCE
        # Overall worth balance + stocks
        self.overall_worth = self.balance

        # Set the current step to a random point within the data frame
        self.current_step = random.randint(0, RESET_STEP_RANDOM_LIMIT)

        return self._next_observation()

    def render(self, mode='human', close=False):
        # Render the environment to the screen
        profit = self.overall_worth - INITIAL_ACCOUNT_BALANCE

        # TODO Refactor
        logging.debug(f'Step: {self.current_step}')
        logging.debug(f'Balance: {self.balance}')
        logging.debug(f'Shares held: {self.stock_qty}')
        logging.debug(f'Profit: {profit}')
