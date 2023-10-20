from datahub_pysdk.dataHub import EAMApi
import pandas as pd
import yaml
import os
import stock_tech as sta

current_path = os.path.dirname(os.path.abspath(__file__))
with open(current_path +"/config.yaml", "r") as file:
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
    # where='trade_date > \'2023-09-01\''
)

# raw data
# print(df)

# sma
sma = sta.sma(df['close'], 5)
macd = sta.macd(df['close'])

df = pd.concat([
    df[['trade_date', 'close']], 
    # sma,
    macd,
], axis=1)
print(df)

# df.to_csv('tmpfiles/gzmt.csv', index=False)
