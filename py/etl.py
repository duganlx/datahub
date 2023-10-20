from datahub_pysdk.dataHub import EAMApi
import pandas as pd
import yaml
import os
import stock_tech as sta

# == get raw data ==
current_path = os.path.dirname(os.path.abspath(__file__))
with open(current_path + "/config.yaml", "r") as file:
    yaml_data = yaml.load(file, Loader=yaml.FullLoader)

pysdk_conf = yaml_data['pysdk']
del yaml_data

eamApi = EAMApi(
    datahub=pysdk_conf['url'],
    user=pysdk_conf['user'],
    password=pysdk_conf['password']
)

df = eamApi.GetData(
    db_name='dm_histdata',
    table_name='bar_day',
    verbose=False,
    universe=['600519.SH'],
    fields=[
        'trade_date',
        'symbol',
        'pre_close',
        'open',
        'high',
        'low',
        'close',
        'total_vol',
        'total_amt',
        'upper_limit',
        'lower_limit',
    ],
    orderby='order by trade_date',
    # where='trade_date > toDateTime64(\'2023-09-01\', 3, \'Asia/Shanghai\')' # is ok
    where='trade_date > \'2023-07-01\''
)

df['trade_date'] = pd.to_datetime(df['trade_date'])
# raw data
# print(df)

# == data processing ==
# sma
sma = sta.sma(df['close'], 5)
macd = sta.macd(df['close'])
kdj = sta.kdj(df['close'], df['high'], df['low'], verbose=True)

df = pd.concat([
    df[['trade_date']],
    # sma,
    # macd,
    kdj
], axis=1)
# print(df)

# == date range filter ==
begin_date = pd.to_datetime('2023-07-05')
end_date = pd.to_datetime('2023-08-05')
drf_df = df[(df["trade_date"] >= begin_date) & (df["trade_date"] <= end_date)]
print(drf_df)

# df.to_csv('tmpfiles/gzmt.csv', index=False)
