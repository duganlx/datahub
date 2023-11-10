import pandas as pd
import py.stockta.ta as sta

# == data processing ==
# sma = sta.sma(df['close'], 5)
# macd = sta.macd(df['close'])
# kdj = sta.kdj(df['close'], df['high'], df['low'], verbose=True)
# boll = sta.boll(df['close'])

# df = pd.concat([
#     df[['trade_date']],
#     # sma,
#     # macd,
#     # kdj,
#     boll
# ], axis=1)
# print(df)

# == date range filter ==
# begin_date = pd.to_datetime('2023-09-01')
# end_date = pd.to_datetime('2023-10-20')
# drf_df = df[(df["trade_date"] >= begin_date) & (df["trade_date"] <= end_date)]
# print(drf_df)

# == value location ==
# nan_index = df['K'].index[df['K'].isna()].tolist()[0]
# print(df.loc[nan_index-5:nan_index+5])

# min_index = df['J'].idxmin()
# print("[min]", df.loc[min_index])

# max_index = df['J'].idxmax()
# print("[max]", df.loc[max_index])

# == value  ==
# print('[kdj]', len(df), len(df[df['J'] > 100]), len(df[df['J'] < 0]))
