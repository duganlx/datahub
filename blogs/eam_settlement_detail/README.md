# 结算明细报表 设计实现

## 资产概要

该页面所涉及到的数据均保存到`UniverseData`对象中。该对象是由 `Balance` 表延伸而来，包含了很多延伸计算出来的字段。具体的字段如下所示。

| 字段名                         | 中文名                   | 备注                                                                                                                                                                |
| ------------------------------ | ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| auCode                         | 资金账号                 |                                                                                                                                                                     |
| auName                         | 资金账号名称             |                                                                                                                                                                     |
| tradeDate                      | 交易日                   | 日期基准字段，13 位时间戳                                                                                                                                           |
| currency                       | 币种                     | 人民币 CNY                                                                                                                                                          |
| totalAssetInitial              | 日初总资产               |                                                                                                                                                                     |
| totalAsset                     | 总资产                   |                                                                                                                                                                     |
| equityInitial                  | 日初持仓市值             |                                                                                                                                                                     |
| equity                         | 持仓市值                 |                                                                                                                                                                     |
| fundInitial                    | 日初资金                 |                                                                                                                                                                     |
| balance                        | 资金余额                 |                                                                                                                                                                     |
| totalLiabilityInitial          | 日初总负债               |                                                                                                                                                                     |
| totalLiability                 | 总负债                   |                                                                                                                                                                     |
| cashDebtInitial                | 日初资金负债             |                                                                                                                                                                     |
| cashDebt                       | 资金负债                 |                                                                                                                                                                     |
| securityDebtInitial            | 日初证券负债             |                                                                                                                                                                     |
| securityDebt                   | 证券负债                 |                                                                                                                                                                     |
| netEquityTraded                | 净买入市值               |                                                                                                                                                                     |
| equityBuy                      | 买入市值                 |                                                                                                                                                                     |
| equitySell                     | 卖出市值                 |                                                                                                                                                                     |
| fundDepositWithdraw            | 净出入金                 | 资金转入 - 资金转出                                                                                                                                                 |
| fundDeposit                    | 资金转入                 |                                                                                                                                                                     |
| fundWithdraw                   | 资金转出                 |                                                                                                                                                                     |
| equityDeposit                  | 证券转入                 |                                                                                                                                                                     |
| equityWithdraw                 | 证券转出                 |                                                                                                                                                                     |
| commission                     | 手续费                   |                                                                                                                                                                     |
| settleTime                     | 清算时间                 |                                                                                                                                                                     |
| equityInTransit                | 在途市值                 |                                                                                                                                                                     |
| fundAvailable                  | 可用资金                 |                                                                                                                                                                     |
| fundInTransit                  | 在途资金                 |                                                                                                                                                                     |
| fundFrozen                     | 冻结资金                 |                                                                                                                                                                     |
| type                           | 账户类型                 |                                                                                                                                                                     |
| createTime                     | 成交时间                 |                                                                                                                                                                     |
| updateTime                     | 更新时间                 |                                                                                                                                                                     |
| isT0                           | 是否 T+0                 | 根据`ads_eqwads_unit_label_value`表中 label 为 strategy，value 为 T0 和 T1 的记录。如果当天有 T1 的记录，则直接判定为 _非 T0_；否则根据当天是否有 T0 记录进行判定。 |
| isValid                        | 是否有效                 | 头尾如果出现 `[持仓市值, 证券负债, 手续费]` 都为 0，则判定为无效数据，中间部分如果连续三天出现这三个字段为 0 的话，也判定为无效数据                                 |
| totalAssetPnl                  | 当日盈亏                 |                                                                                                                                                                     |
| totalAssetPnlCum               | 累计盈亏                 |                                                                                                                                                                     |
| prevTotalAssetPnlCum           | 昨日累计盈亏             |                                                                                                                                                                     |
| totalAssetPnlPercentage        | 当日盈亏%                |                                                                                                                                                                     |
| totalAssetPnlCumPercentage     | 累计盈亏%                |                                                                                                                                                                     |
| prevTotalAssetPnlCumPercentage | 昨日累计盈亏%            |                                                                                                                                                                     |
| verifyTotalAssetInitial        | 核算字段: 日初总资产     | 日初持仓市值 + 日初资金余额                                                                                                                                         |
| isOkTotalAssetInitial          | 验证字段结果: 日初总资产 | 如果核算的结果和取数回来的结果一致，则为 true；反之为 false                                                                                                         |
| verifyTotalAsset               | 核算字段: 总资产         | 持仓市值 + 在途市值 + 资金余额                                                                                                                                      |
| isOkTotalAsset                 | 验证字段结果: 总资产     | 如果核算的结果和取数回来的结果一致，则为 true；反之为 false                                                                                                         |
| verifyTotalLiability           | 核算字段: 总负债         | 资金负债 + 证券负债                                                                                                                                                 |
| isOkTotalLiability             | 验证字段结果: 总负债     | 如果核算的结果和取数回来的结果一致，则为 true；反之为 false                                                                                                         |
| banchmarkPnlPercentage         | 基准盈亏%                | 数据取 `dm_histdata.bar_day`，按照公式 pnl% = (当日收盘价 - 昨日收盘价) / 昨日收盘价 \* 100% 计算得到                                                               |
| banchmarkPnlCumPercentage      | 基准累计盈亏%            |                                                                                                                                                                     |
| benchmarkPreClose              | 基准指数昨收             | bar_day 表的字段 pre_close                                                                                                                                          |

计算公式汇总

**当日盈亏** 的计算目前存在三种情况

- T0 交易

```text
当日盈亏 = 卖出市值 + 净出入金 - 买入市值 + 资金转出 - 资金转入
```

- 非 T0 交易，按日初总资产

```text
日末资产 = 总资产 - 总负债 + 资金转出 + 证券转出
日初资产 = 日初总资产 - 日初总负债 + 资金转入 + 证券转入
当日盈亏 = 日末资产 - 日初资产
当日盈亏% = (日末资产 / 日初资产 - 1) * 100%
```

- 非 T0 交易，按日初持仓市值。如果日末资产 ≤0，则当日盈亏%为 0

```text
日末资产 = 持仓市值 - 证券负债 + 证券转出
日初资产 = 日初持仓市值 - 日初证券负债 + 证券转入
当日盈亏 = 日末资产 - 日初资产
# 如果日末资产 <= 0，则当日盈亏% = 0
当日盈亏% = (日末资产 / 日初资产 - 1) * 100%
```

**当日盈亏(对冲)** 的计算目前存在二种情况，

- 指数

```
当日盈亏(对冲) = (日初持仓市值 - 日初证券负债) * 基准指数pnl%
当日盈亏%(对冲) = 当日盈亏(对冲) / 日初资产

# 日初资产的计算公式存在 按日初总资产 和 按日初持仓市值 两种
# 情况1：按日初总资产
日初资产 = 日初总资产 - 日初总负债 + 资金转入 + 证券转入
# 情况2：按日初持仓市值
日初资产 = 日初持仓市值 - 日初证券负债 + 证券转入
```

- 虚拟期值

```text
对冲张数 = 200
对冲票数 = round(日初持仓市值 / (基准指数的昨日收盘价 * 对冲张数))
当日盈亏(对冲) = 对冲票数 * (基准指数的昨日收盘价 * 对冲张数) * 基准指数pnl%
当日盈亏%(对冲) = 当日盈亏(对冲) / 日初资产

# 日初资产的计算公式存在 按日初总资产 和 按日初持仓市值 两种
# 情况1：按日初总资产
日初资产 = 日初总资产 - 日初总负债 + 资金转入 + 证券转入
# 情况2：按日初持仓市值
日初资产 = 日初持仓市值 - 日初证券负债 + 证券转入
```

- 主力合约 (目前暂未实现)

Alpha 的计算为

```
alpha = 当日盈亏 - 当日盈亏(对冲)
alpha% = 当日盈亏% - 当日盈亏(对冲)%
```

_累计 xxx_ 的计算可以抽象为 `累计* = 昨日累计* + 当日*`

## 按产品

todo

## 按投资经理

在按投资经理的结算明细报表页面中, 左侧的目录树为三级结构, 各级的关系是 _基金经理-产品-资产单元_. 实现上是先取 `dim_datahub.dim_unit_account_product` 表, 该表有 _资产单元_ 和 _产品_ 之间的映射关系, 接着利用 `ads_eqw.ads_unit_label_value` 表可以取得 _资产单元_ 和 _基金经理_ 之间的映射关系, 通过用户中心的接口`/api/uc/v1/users` 可以请求得到所有用户的信息. 具体如下所示.

- dim*unit_account_product 表: \_unit_code* 资产单元编码, _unit_name_ 资产单元名称, unit*type 资产单元类型, account_code 资金账号编码, account_name 资金账号名称, account_type 资金账号类型, \_product_inner_code* 产品内部编码, fund*record_number 产品协会编号, \_product_short_name* 产品名称简称, product_full_name 产品名称全称, product_type 产品类型, etl_time 数据入库时间. 目前, 仅仅展示 `unit_type=[1, 3]` 的, 跟凯强确认了下 该字段存在三种取值 `1 普通资产单元, 2 默认资产单元, 3 客户资产单元`.

- ads*unit_label_value 表: deal_date 日期, \_au_code* 资产单元, label 标签, _value_ 标签内容. 设置 `label = 'manager'`, `au_code - value` 就是资产单元和基金经理的映射. 需要注意的是因为存在日期的维度, 存在一个资产单元在不同的日期隶属于不同的基金经理的情况, 该情况在展示上就是每个基金经理都会有该资产单元.

- /api/uc/v1/users 接口: _id_, _userName_, nickName, email, mobile, avatar, status, ext, roles, sex, depts, qywxId, createAt. `status=0` _应该_ 是属于正常状态.

具体实现逻辑是当取得这三张表的数据, 根据 ads_unit_label_value 表生成一个 `Map<基金经理, 资产单元[]>`, 使用 dim_unit_account_product 表生成一个 `Map<资产单元, 产品>`, 根据前面两个 Map 可以产生 `Map<基金经理, Map<产品, 资产单元[]>>`, 利用用户中心 userName 取得基金经理信息进行绑定.

```js
// products 使用 product_inner_code 作为 key
InvestMgrItem{userId: string; userName: string; nickName: string; category: 'investMgr'; name: string; products: Record<string, InvestMgrProduct>}
// units 使用 unit_code 作为 key
InvestMgrProduct{category: 'investProduct'; type: string; name: string; fullName: string; code: string; units: Record<string, StlUnitItem>}
StlUnitItem{category: 'unit'; type: string; name: string; fullName: string; code: string; isDefault: boolean}

// 关联关系
// ads_unit_label_value.au_code --- dim_unit_account_product.
```

### 资产概要 - 具体某日指标

按投资经理的资产概要页面中，有一个统计了标定时间范围内的业绩表现走势图，以及一个查看具体某一天表现的 _块块_，如下图红色框起来的部分。在这个 _块块_ 中有具体三个指标，分别是 `基准盈亏(万)`、`区间盈亏(万)`、`区间超额(万)`。基准盈亏是根据下拉选择的基准指数而不同，具体有 A 股的指数、港股的指数、美股的指数（数据保存在`ads_eqw.ads_eqw_benchmark`）。结算的数据取自 `ads_eqw.ads_unit_balance_pending`，而基准指数的数据取自表 `dm_histdata.bar_day`。

由于不同市场开市情况大不相同，所以存在一种情况，结算有日期 a 的数据，而基准指数没有日期 a 的数据，在这种情况下，处理方式为取距离日期 a 最近的一次有数据的记录作为日期 a 的基准指数的数据进行接下来的计算。

![demo](aSummaryIndex.png)
