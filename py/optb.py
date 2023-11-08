# 用于操作数据集市表
import pandas as pd
from dtatool import eamapiHandler

eamApi = eamapiHandler()

def matrix2kv(matrix, verbose=False):
  title_index = 0
  title_row = odata[title_index]
  matrix_row_len = len(odata)
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


# == Project: 金葵花帮助 ==
odata = [
  ['id', 'title', 'filepath', 'createTime', 'updateTime', 'publish', 'archive'],
  [1, '进化论代码管理规范及使用指引', '/examples/code-control/CodeControl.md', '2023-05-01', '2023-05-01', 1, '1. 信息安全专题'],
  [2, '进化论文件类数据管理规范', '/examples/dataopts/data_file_opts.md', '2023-06-01', '2023-06-01', 1, '1. 信息安全专题'],
  [3, 'Remote-SSH的配置方式', '/examples/ssh_vscode/ConfigRemoteSSH.md', '2022-11-01', '2022-11-01', 1, '2. 环境配置专题'],
  [4, 'Dev Container的配置方式', '/examples/ssh_vscode/DevContainerGsf2.md', '2022-12-01', '2022-12-01', 1, '2. 环境配置专题'],
  [5, '在服务器上安装配置Docker环境', '/examples/config-docker/configure_docker_env.md', '2022-12-01', '2022-12-01', 1, '2. 环境配置专题'],
  [6, '将gsf服务托管至非内网环境', '/examples/hosting-dataservice/hosting_dataservice.md', '2023-08-01', '2023-08-01', 1, '2. 环境配置专题'],
  [7, 'gsfctl工具的使用', '/examples/gsfctl/GsfCtlHelp.md', '2023-05-01', '2023-05-01', 1, '3. 金葵花2.0开发专题'],
  [8, 'gsf2从服务器构建环境和快速使用', '/examples/gsfctl/GsfQuickGuide.md', '2023-01-01', '2023-01-01', 1, '3. 金葵花2.0开发专题'],
  [9, '如何开发用户的工具库，且在多项目中共享使用', '/examples/gsfctl/develop_shared_conan_package.md', '2023-01-01', '2023-01-01', 1, '3. 金葵花2.0开发专题'],
  [10, 'pygsf使用 - PyGsfRpc', '/examples/pygsf/pygsf_rpc.md', '2023-11-08', '2023-11-08', 1, '3. 金葵花2.0开发专题'],
  [11, 'pygsf使用 - Oms', '/examples/pygsf/pygsf_oms.md', '2023-11-08', '2023-11-08', 1, '3. 金葵花2.0开发专题'],
  [12, 'pygsf使用 - Model', '/examples/pygsf/pygsf_model.md', '2023-11-08', '2023-11-08', 1, '3. 金葵花2.0开发专题'],
  [13, 'pygsf使用 - Backtest', '/examples/pygsf/pygsf_backtest.md', '2023-11-08', '2023-11-08', 1, '3. 金葵花2.0开发专题'],
  [14, '高效哈希表phmap', '/examples/phmap/phmap.md', '2023-08-01', '2023-08-01', 1, '4. 开发笔记专题'],
  [15, 'C++ memory order', '/examples/memory-model/memory-order.md', '2023-09-01', '2023-09-01', 1, '4. 开发笔记专题']
]

eamApi.UploadData(
  table_name='gsfhelp',
  data=pd.DataFrame(matrix2kv(odata, verbose=True)),
  primary_key=['id'],
  append=True,
  replace=True,
  verbose=True,
  public_table_sign=True
)