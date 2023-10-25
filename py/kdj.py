# 分析 kdj
import pandas as pd
import stock_tech as sta
from dtatool import iodta, rolling_df

df = iodta(
  src="eam", 
  
  # CONF: src=eam
  save={'ctl': False, 'dir': 'tmp', 'filename': 'xxx.csv'}, 
  where='trade_date > \'2023-07-01\'',
  universe=['600519.SH'],
  
  # CONF: src=local
  read={'dir': 'tmp', 'filename': 'xxx.csv'},
)

kdj = sta.kdj(df['close'], df['high'], df['low'], verbose=False, wins=9)
# kdj_3d = sta.kdj(df['close'], df['high'], df['low'], verbose=False, wins=27)

df = pd.concat([df[['trade_date']], kdj], axis=1)
# print(df)

def bullish(rdf):
  """
  kdj 金叉
  j 从下突破 k 和 d
  """
  first, second = rdf.iloc[0], rdf.iloc[1]

  c1 = False
  if first.at['J'] <= first.at['K'] and first.at['J'] <= first.at['D']:
    c1 = True
  
  c2 = False
  if second.at['J'] >= second.at['K'] and second.at['J'] >= second.at['D']:
    c2 = True

  return c1 and c2

gs = rolling_df(df, 2, bullish, "good").fillna(False)
df = pd.concat([df, gs], axis=1)

print(df[df['good'] == True])