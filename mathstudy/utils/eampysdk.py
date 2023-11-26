from datahub_pysdk.dataHub import EAMApi
import os
import yaml
import pandas as pd

class EamPySdk(object):
  # basic configuration
  # == begin ==
  conf_filepath = ''
  rootfolder_path = ''
  pyfolder_path = ''
  datafolder_path = ''
  url = ''
  user = ''
  password = ''
  # == end ==

  apihandler: EAMApi = None
  utils = None

  def __init__(self) -> None:
    self.initBasicconf()
    self.initApihandler()
    self.utils = Utils()
  
  def initBasicconf(self) -> None:
    current_path = os.path.dirname(os.path.abspath(__file__))
    pyfolder_path = os.path.dirname(current_path)
    rootfolder_path = os.path.dirname(pyfolder_path)
    conf_filepath = f"{pyfolder_path}/config.yaml"
    datafolder_path = f"{rootfolder_path}/data"

    with open(conf_filepath, "r") as file:
        yaml_data = yaml.load(file, Loader=yaml.FullLoader)

    confs = yaml_data['eampysdk']
    del yaml_data

    url = confs['url']
    user = confs['user']
    password = confs['password']

    self.conf_filepath = conf_filepath
    self.rootfolder_path = rootfolder_path
    self.datafolder_path = datafolder_path
    self.url = url
    self.user = user
    self.password = password

    print(f"[Basic config] confpath={conf_filepath}, url={url}, user={user}, password={password}")

  def initApihandler(self):
    url = self.url
    user = self.user
    pwd = self.password

    eamApi = EAMApi(
        datahub=url,
        user=user,
        password=pwd
    )
    self.apihandler = eamApi

  def uploadData(self, table_name: str, data: list, db_name: str = None, primary_key: list = [], append: bool = False, replace: bool = False, universe: str = None, datetime: str = None, verbose: bool = False, partition_by: list = None, public_table_sign: bool = True, optimize: bool = True):
    """
    data 参数传二维数组即可，形如
    [
      ['a', 'b', 'c'],
      ['r1a', 'r1b', 'r1c'],
      ['r2a', 'r2b', 'r2c'],
    ]
    """
    datadf = self.utils.matrix2kv(data)
    datadf = pd.DataFrame(datadf)

    self.apihandler.UploadData(
      table_name=table_name,
      data=datadf,
      db_name=db_name,
      primary_key=primary_key,
      append=append,
      replace=replace,
      universe=universe,
      datetime=datetime,
      verbose=verbose,
      partition_by=partition_by,
      public_table_sign=public_table_sign,
      optimize=optimize
    )

  def getBardayData(self, universe=[], where=''):
    dbname = 'dm_histdata'
    tbname = 'bar_day'
    fields = ['trade_date', 'symbol', 'pre_close', 'open', 'high', 'low', 'close', 'total_vol', 'total_amt', 'upper_limit', 'lower_limit']
    df = self.apihandler.GetData(
      db_name=dbname,
      table_name=tbname,
      verbose=False,
      universe=universe,
      fields=fields,
      orderby='order by trade_date',
      where=where
    )

    df['trade_date'] = pd.to_datetime(df['trade_date'])

    return df

  def savedf(self, df: pd.DataFrame, dir: str, filename: str):
    if self.datafolder_path == '':
      raise Exception("data folder path is not configured")

    df.to_csv(f"{self.datafolder_path}/{dir}/{filename}", index=False)
    print(f"df has save in {dir}/{filename}")

class Utils(object):

  def matrix2kv(self, matrix: list, verbose: bool=False):
    title_index = 0
    title_row = matrix[title_index]
    matrix_row_len = len(matrix)
    matrix_col_len = len(title_row)
    row_index = 1

    if verbose: 
      print('field name: ', title_row)
    kv = {}
    for k in title_row: 
      kv[k] = []

    while row_index < matrix_row_len: 
      data_row = matrix[row_index]
      row_index += 1

      col_index = 0
      while col_index < matrix_col_len:
        k = title_row[col_index]
        v = data_row[col_index]
        col_index += 1
        kv[k].append(v)

    if verbose: 
      print(kv)

    return kv