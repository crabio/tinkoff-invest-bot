from preprocess import preprocess_data
from bot import TradingBot

def test_bot_profit(df):
    # Preprocess data
    bot_df = preprocess_data(df)
    # Select columns
    bot_df = bot_df[["ts","open_price","close_price","high_price","low_price","volume","MACD","MACDs","RSI"]]

    # Create bot instance
    bot = TradingBot(rsi_oversold_threshold=20,
                    rsi_overbought_threshold=80,
                    stop_loss=0.005)
    # Reset bot
    bot.reset()
    # Run bot
    for _, data in bot_df.iterrows():
        bot.process(data)

    return bot.profit * 100