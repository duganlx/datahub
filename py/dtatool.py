from datahub_pysdk.dataHub import EAMApi
import os
import pandas as pd
import numpy as np
import yaml

def _getDtabasepath():
   """
   获取保存数据的根路径
   """
   current_path = os.path.dirname(os.path.abspath(__file__))
   project_path = os.path.dirname(current_path)
   dta_path = project_path + '/data'

   return dta_path

def _fromEamApi(universe, where):
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
    universe=universe,
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
    where=where
    # where='trade_date > toDateTime64(\'2023-09-01\', 3, \'Asia/Shanghai\')' # is ok
  )

  df['trade_date'] = pd.to_datetime(df['trade_date'])

  return df

def iodta(src, **kwargs) -> pd.DataFrame:
  """
  :param src:来源 eam, local
  :param save:保存配置(条件 src=eam), 包含三个属性 ctl 是否保存; filename 文件名; dir 保存的文件夹
  :param read:读取文件的配置(条件 src=local), 包含两个属性 filename 文件名; dir 保存的文件夹
  :param universe:股票列表(条件 src=eam)
  :param where:查询过滤条件字符串(条件 src=eam)
  """
  if src == "eam":
     savecfg = kwargs.get("save", {"ctl": False})

     if "universe" not in kwargs:
        raise Exception("universe is not configured.")

     universe = kwargs.get("universe")
     where = kwargs.get("where", "")

     df = _fromEamApi(
        universe=universe, 
        where=where
      )

     if savecfg.get('ctl', False):
        filedir = savecfg.get('dir', 'tmp')

        if "filename" not in savecfg:
           raise Exception(f"filename is not configured")

        filename = savecfg.get('filename')
        bfp = _getDtabasepath()
        
        df.to_csv(f"{bfp}/{filedir}/{filename}", index=False)

  elif src == "local":
     readcfg = kwargs.get("read", {})

     if "dir" not in readcfg or "filename" not in readcfg:
        raise Exception("reading file requires configuring the following properties: [dir, filename]")

     bfp = _getDtabasepath()
     filedir = readcfg.get('dir')
     filename = readcfg.get('filename')

     df = pd.read_csv(f"{bfp}/{filedir}/{filename}")
  else:
     raise Exception(f"This source is currently not supported: {src}")
  
  return df


def rolling_df(df, window, func, name):
    """
    对Dataframe按一行为单位进行滚动, 并且生成结果Series
    """
    res = []
    for i in range(len(df)):
        window_df = df.iloc[(i-window+1):i+1, :]

        if window_df.empty:
            res.append(np.nan)
            continue

        this_res = func(window_df)
        res.append(this_res)
    
    res = pd.Series(res, name=name)
    return res