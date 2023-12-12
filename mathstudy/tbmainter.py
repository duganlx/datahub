from utils.eampysdk import EamPySdk

def settlementGSJZ():
    """
    对冲类型: 公司基准
    ads_eqw
    """
    eamsdk = EamPySdk()
    data = [
        ['year', 'pnl'],
        [2021, 10],
        [2022, 4],
        [2023, 4],
        [2024, 3],
    ]

    eamsdk.uploadData(
        table_name="ads_settlement_hedge_gsjz",
        data=data,
        primary_key=['year'],
        append=True,
        replace=True,
        verbose=True,
        public_table_sign=True,
        db_name="ads_eqw"
    )


def gsfhelpConf():
    """
    gsf帮助项目 目录树维护
    """
    eamsdk = EamPySdk()
    data = [
        ['id', 'title', 'filepath', 'createTime', 'updateTime', 'publish', 'archive'],
        # [1, '进化论代码管理规范及使用指引', '/examples/code-control/CodeControl.md', '2023-05-01', '2023-05-01', 1, '1. 信息安全专题'],
        # [2, '进化论文件类数据管理规范', '/examples/dataopts/data_file_opts.md', '2023-06-01', '2023-06-01', 1, '1. 信息安全专题'],
        # [3, 'Remote-SSH的配置方式', '/examples/ssh_vscode/ConfigRemoteSSH.md', '2022-11-01', '2022-11-01', 1, '2. 环境配置专题'],
        # [4, 'Dev Container的配置方式', '/examples/ssh_vscode/DevContainerGsf2.md', '2022-12-01', '2022-12-01', 1, '2. 环境配置专题'],
        # [5, '在服务器上安装配置Docker环境', '/examples/config-docker/configure_docker_env.md', '2022-12-01', '2022-12-01', 1, '2. 环境配置专题'],
        # [6, '将gsf服务托管至非内网环境', '/examples/hosting-dataservice/hosting_dataservice.md', '2023-08-01', '2023-08-01', 1, '2. 环境配置专题'],

        # 1. gsfctl 系列
        [7, 'gsfctl 从服务器构建环境和快速使用', '/examples/gsfctl/GsfQuickGuide.md', '2023-01-01', '2023-01-01', 1, '1. gsfctl 系列'],
        [8, 'gsfctl 工具的使用', '/examples/gsfctl/GsfCtlHelp.md', '2023-05-01', '2023-05-01', 1, '1. gsfctl 系列'],
        [9, 'gsfctl 如何开发用户的工具库，且在多项目中共享使用', '/examples/gsfctl/develop_shared_conan_package.md', '2023-01-01', '2023-01-01', 1, '1. gsfctl 系列'],

        # 2. pygsf系列
        [10, 'pygsf使用 - PyGsfRpc', '/examples/pygsf/pygsf_rpc.md', '2023-11-08', '2023-11-08', 1, '2. pygsf系列'],
        [11, 'pygsf使用 - Oms', '/examples/pygsf/pygsf_oms.md', '2023-11-08', '2023-11-08', 1, '2. pygsf系列'],
        [12, 'pygsf使用 - Model', '/examples/pygsf/pygsf_model.md', '2023-11-08', '2023-11-08', 1, '2. pygsf系列'],
        [13, 'pygsf使用 - Backtest', '/examples/pygsf/pygsf_backtest.md', '2023-11-08', '2023-11-08', 1, '2. pygsf系列'],

        # 3. gsfctl 系列
        [20, '高频数据服务（DataService) 框架：WWDS', '/examples/WWDS/WWDS.md', '2023/10/20', '2023/10/20', 1, '3. gsfctl 系列'],
        [21, '高频Touch类算法（Algo）框架：WWAlgo', '/examples/WWAlgo/WWAlgo.md', '2023/10/21', '2023/10/21', 1, '3. gsfctl 系列'],

        # 4. DataService系列
        [22, 'BarDs 使用', '/examples/bards/bards.md', '2023/11/16', '2023/11/16', 1, '4. DataService系列'],
        [23, 'KLineDs 使用', '/examples/klineds/klineds.md', '2023/11/16', '2023/11/16', 1, '4. DataService系列'],

        # 6. 开发笔记专题
        # [14, '高效哈希表phmap', '/examples/phmap/phmap.md', '2023-08-01', '2023-08-01', 1, '6. 开发笔记专题'],
        # [15, 'C++ memory order', '/examples/memory-model/memory-order.md', '2023-09-01', '2023-09-01', 1, '6. 开发笔记专题']
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


if __name__ == '__main__':
    settlementGSJZ()