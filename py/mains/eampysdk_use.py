"""
eampysdk test
"""
import os
import sys
import pandas as pd
import torch
import numpy as np

# 解决 ModuleNotFoundError: No module named 'xxx' 问题
this_file_full_path_name = os.path.abspath(__file__)
this_file_folder_path = os.path.dirname(this_file_full_path_name)
parent_folder_path = os.path.dirname(this_file_folder_path)
sys.path.append(parent_folder_path)

from eampysdk.index import EamPySdk
from lcpysdk.index import LocalPySdk
from stockta.ta import StockTA

def gsfhelpConf():
    """
    gsf帮助项目 目录树维护
    """
    eamsdk = EamPySdk()
    data = [
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

    eamsdk.uploadData(
        table_name='gsfhelp',
        data=data,
        primary_key=['id'],
        append=True,
        replace=True,
        verbose=True,
        public_table_sign=True
    )

def adsStatementStatus():
    """
    对账单存续状态表
    ads_eqw.ads_statement_status
    """
    eamsdk = EamPySdk()
    data = [
        ['au_code', 'account_name_cn', 'settle_status', 'statement_status'],
        ['270090005318', '达复一安信', 0, 0],
        ['666800007983', '达尔文达复合一号华泰', 1, 1],
        ['902090000445', '达尔文达复合一号华泰信用', 1, 2],
    ]

    eamsdk.uploadData(
        table_name='test_ads_statement_status',
        data=data,
        primary_key=['au_code'],
        append=True,
        replace=True,
        verbose=True,
        public_table_sign=True
    )

def stockTAanalysis(local: bool = False):
    if local:
        lcsdk = LocalPySdk()
        df = lcsdk.getData(dir='tmp', filename='raw.csv')
    else:
        eamsdk = EamPySdk()
        df = eamsdk.getBardayData(
            universe=['600519.SH'],
            # where='trade_date > \'2023-07-01\''
        )
        # eamsdk.savedf(df, dir='tmp', filename='raw.csv')

    stockTa = StockTA(df)

    # ma = stockTa.ma(5)
    # ema = stockTa.expma(5)
    # macd = stockTa.macd()
    # kdj = stockTa.kdj()
    # boll = stockTa.boll()
    # mtm = stockTa.mtm()
    # rsi = stockTa.rsi()
    # dmi = stockTa.dmi()
    # dma = stockTa.dma()
    # brar = stockTa.brar()
    # obv = stockTa.obv(offset=32.352-815.769, verbose=True)
    wr = stockTa.wr(n=10)

    print(wr)

def biclassify(generate=False):
    if generate:
        eamsdk = EamPySdk()
        df = eamsdk.getBardayData(
            universe=['600519.SH'],
            where='trade_date > \'2023-07-01\''
        )
        stockTa = StockTA(df)
        dmatrix = stockTa.data_matrix()
        eamsdk.savedf(df=dmatrix, dir='tmp', filename='ta.csv')
        df = dmatrix
    else:
        lcsdk = LocalPySdk()
        df = lcsdk.getData(dir='tmp', filename='ta.csv')

    # df format
    # == begin ==
    print(df.values)
    # n = df.shape[0]
    # k = 5
    # height = n // k
    # print(n, k, height)
    # small_dfs = np.split(df, [i * height for i in range(1, k)], axis=0)
    # print(small_dfs)
    # arr = np.array(small_dfs)
    # tensor = torch.tensor(small_dfs)
    # == end ==
    # print(arr)

    # print(df)
    # torch.nn.LSTM(input_size=38, hidden_size=20, num_layers=2, bidirectional=False)

if __name__ == '__main__':
    # gsfhelpConf()
    # adsStatementStatus()
    # stockTAanalysis(local=False)

    biclassify(generate=False)
    # rnn = torch.nn.LSTM(input_size=10, hidden_size=20, num_layers=2, bidirectional=True)
    input = torch.randn(5, 3, 10)#(seq_len, batch, input_size)
    # h0 = torch.randn(4, 3, 20) #(num_layers,batch,output_size)
    # c0 = torch.randn(4, 3, 20) #(num_layers,batch,output_size)
    # output, (hn, cn) = rnn(input, (h0, c0))

    # print(output)
    # print(hn, cn)