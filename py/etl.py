from datahub_pysdk.dataHub import EAMApi

# todo 读取配置文件
eamApi = EAMApi(datahub="eqw.eam.com", user='', password='')

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

print(df)

df.to_csv('tmpfiles/gzmt.csv', index=False)
