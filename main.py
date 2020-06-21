import tinvest
import pandas as pd

TOKEN = "t.IodbD3G5wmc4_vOlZDLdTKe4N9VVpCoEIAmbTcFhqWVACCiXMTKCLYm0GP7v2SRf-ozJb5Vq6Ix67Bv5-ydMDA"

client = tinvest.SyncClient(TOKEN, use_sandbox=True)
api = tinvest.MarketApi(client)

# Get Stocks
response = api.market_stocks_get()  # requests.Response
if response.status_code == 200:
    # Parse
    markets_instruments_list =  response.parse_json().payload
    
    markets_instruments_data = [ {
            'currency': markets_instrument.currency,
            'figi': markets_instrument.figi,
            'isin': markets_instrument.isin,
            'lot': markets_instrument.lot,
            'min_price_increment': markets_instrument.min_price_increment,
            'name': markets_instrument.name,
            'ticker': markets_instrument.ticker,
            'type': markets_instrument.type
        }
        for markets_instrument
        in markets_instruments_list.instruments]

    markets_instruments_df = pd.DataFrame(markets_instruments_data)

    # print(markets_instruments_df)

    print(markets_instruments_df[markets_instruments_df['name'].str.contains("Morgan")])
else:
    print("ERROR:", response.parse_error())