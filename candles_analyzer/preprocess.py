def preprocess_data(df):
    # Sort values with ts
    df = df.sort_values("ts")
    # Reset index
    df = df.reset_index(drop=True)

    # Add MACD
    df.ta.macd(fast=12,slow=26,signal=9,append=True)

    # Add Stochastic RSI
    df.ta.stochrsi(append=True)

    # Rename columns
    df = df.rename(columns={
        "MACD_12_26_9": "MACD",
        "MACDs_12_26_9": "MACDs",
        "MACDh_12_26_9": "MACDh",
        "STOCHRSIk_14_14_3_3": "RSI"})

    # Shift MACD signals
    # TEST!!!
    shift_count = 4
    df["MACD"] = df["MACD"].shift(-shift_count)
    df["MACDs"] = df["MACDs"].shift(-shift_count)
    df["MACDh"] = df["MACDh"].shift(-shift_count)

    # Calc RSI signals
    RSI_oversold_threshold = 20
    RSI_overbouht_threshold = 80

    df["RSI_oversold"] = df["RSI"] <= RSI_oversold_threshold
    df["RSI_overbought"] = df["RSI"] >= RSI_overbouht_threshold

    df["MACD_buy"] = (df["MACDh"] > 0) & (df["MACDh"].shift() <= 0)
    df["MACD_sell"] = (df["MACDh"] <= 0) & (df["MACDh"].shift() > 0)

    # Fill na
    df = df.fillna(0)

    return df