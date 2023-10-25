import pandas as pd
import stock_tech as sta
from dtatool import iodta

df = iodta(
  src="eam", 
  
  # CONF: src=eam
  save={'ctl': True, 'dir': 'tmp', 'filename': 'xxx.csv'}, 
  where='trade_date > \'2023-07-01\'',
  universe=['600519.SH'],
  
  # CONF: src=local
  read={'dir': 'tmp', 'filename': 'xxx.csv'},
)

print(df)


# == data processing ==
# sma = sta.sma(df['close'], 5)
# macd = sta.macd(df['close'])
kdj = sta.kdj(df['close'], df['high'], df['low'], verbose=True)

df = pd.concat([
    df[['trade_date']],
    # sma,
    # macd,
    kdj
], axis=1)
print(df)

# == date range filter ==
# begin_date = pd.to_datetime('2023-09-01')
# end_date = pd.to_datetime('2023-10-20')
# drf_df = df[(df["trade_date"] >= begin_date) & (df["trade_date"] <= end_date)]
# print(drf_df)

# == value location ==
# nan_index = df['K'].index[df['K'].isna()].tolist()[0]
# print(df.loc[nan_index-5:nan_index+5])

min_index = df['J'].idxmin()
print("[min]", df.loc[min_index])

max_index = df['J'].idxmax()
print("[max]", df.loc[max_index])

# == value  ==
# print('[kdj]', len(df), len(df[df['J'] > 100]), len(df[df['J'] < 0]))

