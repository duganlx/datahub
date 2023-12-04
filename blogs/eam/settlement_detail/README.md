# 结算明细报表 设计实现

## 资产概要

该页面所涉及到的数据均保存到`UniverseData`对象中。该对象是由 `Balance` 表延伸而来，包含了很多延伸计算出来的字段。具体的字段如下所示。

计算公式汇总

### 当日盈亏

T0 交易（废弃）

```
当日盈亏 = 卖出市值 + 净出入金 - 买入市值 + 资金转出 - 资金转入
当日盈亏% = 当日盈亏 / 买入市值
```

非 T0 交易，维度: 资产

```text
日末资产 = 总资产 - 总负债 + 资金转出 + 证券转出
日初资产 = 日初总资产 - 日初总负债 + 资金转入 + 证券转入
当日盈亏 = 日末资产 - 日初资产
当日盈亏% = (日末资产 / 日初资产 - 1) * 100%
```

- 非 T0 交易，维度: 市值

```text
日末资产 = 总资产 - 总负债 + 资金转出 + 证券转出
日初资产 = 日初总资产 - 日初总负债 + 资金转入 + 证券转入
当日盈亏 = 日末资产 - 日初资产

日初市值 = 日初持仓市值 - 日初证券负债
# 如果日末资产 <= 0，则当日盈亏% = 0
当日盈亏% = (当日盈亏 / 日初市值) * 100%
```

### 当日盈亏(对冲)

**指数**

```text
基准指数pnl% = dm_histdata.bar_day取close和preclose按照pnl公式计算
当日盈亏(对冲) = (日初持仓市值 + 日初证券负债) * 基准指数pnl%
当日盈亏%(对冲) = 基准指数pnl%
```

历史版本

```
==== v1.0 ====
基准指数pnl% = dm_histdata.bar_day取close和preclose按照pnl公式计算
当日盈亏(对冲) = (日初持仓市值 + 日初证券负债) * 基准指数pnl%
当日盈亏%(对冲) = 当日盈亏(对冲) / 日初资产
-> 日初资产 (资产维度) = 日初总资产 - 日初总负债 + 资金转入 + 证券转入
-> 日初资产 (市值维度) = 日初持仓市值 - 日初证券负债 + 证券转入
```

**虚拟期值**

```text
对冲张数 = 200
对冲票数 = round(日初持仓市值 / (基准指数的昨日收盘价 * 对冲张数))
当日盈亏(对冲) = 对冲票数 * (基准指数的昨日收盘价 * 对冲张数) * 基准指数pnl%
当日盈亏%(对冲) = 基准指数pnl%
```

历史版本

```text
==== 历史版本v1.0 ====
对冲张数 = 200
对冲票数 = round(日初持仓市值 / (基准指数的昨日收盘价 * 对冲张数))
当日盈亏(对冲) = 对冲票数 * (基准指数的昨日收盘价 * 对冲张数) * 基准指数pnl%
当日盈亏%(对冲) = 当日盈亏(对冲) / 日初资产
-> 日初资产 (资产维度) = 日初总资产 - 日初总负债 + 资金转入 + 证券转入
-> 日初资产 (市值维度) = 日初持仓市值 - 日初证券负债 + 证券转入
```

**主力合约**

```text
pnl% = ads_eqw.ads_ic889中的 pnl_close
当日盈亏(对冲) = (日初持仓市值 + 日初证券负债) * pnl%
当日盈亏%(对冲) = pnl%
```

**公司基准**

公司根据实际情况，由一平进行核定每年一个对冲成本的参数

```text
当日盈亏(对冲) = ?
当日盈亏%(对冲) = ?
```

### 当日超额

```text
alpha = 当日盈亏 - 当日盈亏(对冲)
alpha% = 当日盈亏% - 当日盈亏%(对冲)
```

### 累计计算

```text
累计* = 昨日累计* + 当日*
```

### 数据来源

资产单元-标签信息表 ads_eqw.ads_unit_label_value 字段 {deal_date, au_code, label, value}

```sql
# 统计策略在每个资产单元的最早时间
select au_code, `value`, min(deal_date) as deal_date from ads_eqw.ads_unit_label_value where label = 'strategy' group by au_code, `value`
```

### 页面数据流动

闪缩问题定位

begin, end, balance, loading, submitsign, chartDateRange, colDisplay, hiddenColumnFields,

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

---

### UniverseData 对象属性

| 字段名                                          | 中文名                   | 备注                                                        |
| ----------------------------------------------- | ------------------------ | ----------------------------------------------------------- |
| auCode                                          | 资金账号                 |                                                             |
| auName                                          | 资金账号名称             |                                                             |
| tradeDate                                       | 交易日                   | 日期基准字段，13 位时间戳                                   |
| currency                                        | 币种                     | 人民币 CNY                                                  |
| totalAssetInitial                               | 日初总资产               |                                                             |
| totalAsset                                      | 总资产                   |                                                             |
| equityInitial                                   | 日初持仓市值             |                                                             |
| equity                                          | 持仓市值                 |                                                             |
| fundInitial                                     | 日初资金                 |                                                             |
| balance                                         | 资金余额                 |                                                             |
| totalLiabilityInitial                           | 日初总负债               |                                                             |
| totalLiability                                  | 总负债                   | 净资产 = 总资产 - 总负债                                    |
| cashDebtInitial                                 | 日初资金负债             |                                                             |
| cashDebt                                        | 资金负债                 |                                                             |
| securityDebtInitial                             | 日初证券负债             |                                                             |
| securityDebt                                    | 证券负债                 |                                                             |
| netEquityTraded                                 | 净买入市值               |                                                             |
| equityBuy                                       | 买入市值                 |                                                             |
| equitySell                                      | 卖出市值                 |                                                             |
| fundDepositWithdraw                             | 净出入金                 | 资金转入 - 资金转出                                         |
| fundDeposit                                     | 资金转入                 |                                                             |
| fundWithdraw                                    | 资金转出                 |                                                             |
| equityDeposit                                   | 证券转入                 |                                                             |
| equityWithdraw                                  | 证券转出                 |                                                             |
| commission                                      | 手续费                   |                                                             |
| settleTime                                      | 清算时间                 |                                                             |
| equityInTransit                                 | 在途市值                 |                                                             |
| fundAvailable                                   | 可用资金                 |                                                             |
| fundInTransit                                   | 在途资金                 |                                                             |
| fundFrozen                                      | 冻结资金                 |                                                             |
| type                                            | 账户类型                 |                                                             |
| createTime                                      | 成交时间                 |                                                             |
| updateTime                                      | 更新时间                 |                                                             |
| isT0                                            | 是否 T+0                 |                                                             |
| isValid                                         | 是否有效                 |                                                             |
| totalAssetPnl                                   | 当日盈亏                 | 共有三种情况，参看下方公式 _当日盈亏_                       |
| totalAssetPnlCum                                | 累计盈亏                 | 参看下方公式汇总 _累计盈亏_                                 |
| prevTotalAssetPnlCum                            | 昨日累计盈亏             | 计算 _累计盈亏_ 时使用                                      |
| totalAssetPnlPercentage                         | 当日盈亏%                | 共有三种情况，参看下方公式 _当日盈亏%_                      |
| totalAssetPnlCumPercentage                      | 累计盈亏%                | 参看下方公式汇总 _累计盈亏_                                 |
| prevTotalAssetPnlCumPercentage                  | 昨日累计盈亏%            | 计算 _累计盈亏_ 时使用                                      |
| verifyTotalAssetInitial                         | 核算字段: 日初总资产     | 日初持仓市值 + 日初资金余额                                 |
| isOkTotalAssetInitial                           | 验证字段结果: 日初总资产 | `verifyTotalAssetInitial == totalAssetInitial`              |
| verifyTotalAsset                                | 核算字段: 总资产         | 持仓市值 + 在途市值 + 资金余额                              |
| isOkTotalAsset                                  | 验证字段结果: 总资产     | 如果核算的结果和取数回来的结果一致，则为 true；反之为 false |
| verifyTotalLiability                            | 核算字段: 总负债         | 资金负债 + 证券负债                                         |
| isOkTotalLiability                              | 验证字段结果: 总负债     | 如果核算的结果和取数回来的结果一致，则为 true；反之为 false |
| banchmarkPnlPercentage                          | 基准盈亏%                |                                                             |
| banchmarkPnlCumPercentage                       | 基准累计盈亏%            | 参看下方公式汇总 _累计盈亏_                                 |
| benchmarkPreClose                               | 基准指数昨收             | bar_day 表的字段 pre_close                                  |
| zs_totalAssetPnlHedge                           | 当日盈亏(对冲)           | `zs` 开头表示对冲类型为 指数                                |
| zs_totalAssetPnlHedgeCum                        | 累计盈亏(对冲)           | 累计盈亏(对冲) = 昨日累计盈亏(对冲) + 当日盈亏(对冲)        |
| zs_prevTotalAssetPnlHedgeCum                    | 昨日累计盈亏(对冲)       | 计算 _累计盈亏(对冲)_ 时使用                                |
| zs_alpha                                        | 当日超额                 | 当日盈亏 - 当日盈亏(对冲, 对冲类型: 指数)                   |
| zs_alphaCum                                     | 累计超额                 | 累计超额 = 昨日累计超额 + 当日超额\*                        |
| zs_prevAlphaCum                                 | 昨日累计超额             | 计算累计超额时使用                                          |
| zs_totalAssetPnlHedgePercentage                 | 当日盈亏%(对冲)          | 所除分母维度是总资产，参看下方公式汇总 _当日盈亏%(对冲)_    |
| zs_totalAssetPnlHedgeCumPercentage              | 累计盈亏%(对冲)          | 累计盈亏%(对冲) = 昨日累计盈亏%(对冲) + 当日盈亏%(对冲)     |
| zs_prevTotalAssetPnlHedgeCumPercentage          | 昨日累计盈亏%(对冲)      | 计算累计盈亏%(对冲)时使用                                   |
| zs_totalAssetPnlHedgePercentage_rcccsz          | 当日盈亏%(对冲)          | 所除分母维度是总市值，参看下方公式汇总 _当日盈亏%(对冲)_    |
| zs_totalAssetPnlHedgeCumPercentage_rcccsz       | 累计盈亏%(对冲)          |                                                             |
| zs_prevTotalAssetPnlHedgeCumPercentage_rcccsz   | 昨日累计盈亏%(对冲)      |                                                             |
| zs_alphaPercentage                              | 当日超额                 |                                                             |
| zs_alphaCumPercentage                           | 累计超额                 |                                                             |
| zs_prevAlphaCumPercentage                       | 昨日累计超额             |                                                             |
| zs_alphaPercentage_rcccsz                       | 当日超额                 |                                                             |
| zs_alphaCumPercentage_rcccsz                    | 累计超额                 |                                                             |
| zs_prevAlphaCumPercentage_rcccsz                | 昨日累计超额             |                                                             |
| xnqz_ticket                                     | 张数                     | `xnqz` 开头表示对冲类型为 虚拟期指                          |
| xnqz_totalAssetPnlHedge                         | 当日盈亏(对冲)           |                                                             |
| xnqz_totalAssetPnlHedgeCum                      | 累计盈亏(对冲)           |                                                             |
| xnqz_prevTotalAssetPnlHedgeCum                  | 昨日累计盈亏(对冲)       |                                                             |
| xnqz_alpha                                      | 当日超额                 |                                                             |
| xnqz_alphaCum                                   | 累计超额                 |                                                             |
| xnqz_prevAlphaCum                               | 昨日累计超额             |                                                             |
| xnqz_totalAssetPnlHedgePercentage               | 当日盈亏%(对冲)          |                                                             |
| xnqz_totalAssetPnlHedgeCumPercentage            | 累计盈亏%(对冲)          |                                                             |
| xnqz_prevTotalAssetPnlHedgeCumPercentage        | 昨日累计盈亏%(对冲)      |                                                             |
| xnqz_totalAssetPnlHedgePercentage_rcccsz        | 当日盈亏%(对冲)          |                                                             |
| xnqz_totalAssetPnlHedgeCumPercentage_rcccsz     | 累计盈亏%(对冲)          |                                                             |
| xnqz_prevTotalAssetPnlHedgeCumPercentage_rcccsz | 昨日累计盈亏%(对冲)      |                                                             |
| xnqz_alphaPercentage                            | 当日超额%                |                                                             |
| xnqz_alphaCumPercentage                         | 累计超额%                |                                                             |
| xnqz_prevAlphaCumPercentage                     | 昨日累计超额%            |                                                             |
| xnqz_alphaPercentage_rcccsz                     | 当日超额%                |                                                             |
| xnqz_alphaCumPercentage_rcccsz                  | 累计超额%                |                                                             |
| xnqz_prevAlphaCumPercentage_rcccsz              | 昨日累计超额%            |                                                             |

说明

- isT0: 根据`ads_eqwads_unit_label_value`表中 label 为 strategy，value 为 T0 和 T1 的记录。如果当天有 T1 的记录，则直接判定为 _非 T0_；否则根据当天是否有 T0 记录进行判定。
- banchmarkPnlPercentage: 数据取 `dm_histdata.bar_day`，按照公式 pnl% = (当日收盘价 - 昨日收盘价) / 昨日收盘价 \* 100% 计算得到
- isValid: 头尾如果出现 `[持仓市值, 证券负债, 手续费]` 都为 0，则判定为无效数据，中间部分如果连续三天出现这三个字段为 0 的话，也判定为无效数据
