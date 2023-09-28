# 结算明细报表设计&实现

## 按产品

todo

## 按投资经理

在按投资经理的结算明细报表页面中, 左侧的目录树为三级结构, 各级的关系是 _基金经理-产品-资产单元_. 实现上是先取 `dim_datahub.dim_unit_account_product` 表, 该表有 _资产单元_ 和 _产品_ 之间的映射关系, 接着利用 `ads_eqw.ads_unit_label_value` 表可以取得 _资产单元_ 和 _基金经理_ 之间的映射关系, 通过用户中心的接口`/api/uc/v1/users` 可以请求得到所有用户的信息. 具体如下所示.

- dim_unit_account_product 表: <u>unit_code 资产单元编码</u>, <u>unit_name 资产单元名称</u>, unit_type 资产单元类型, account_code 资金账号编码, account_name 资金账号名称, account_type 资金账号类型, <u>product_inner_code 产品内部编码</u>, fund_record_number 产品协会编号, <u>product_short_name 产品名称简称</u>, product_full_name 产品名称全称, product_type 产品类型, etl_time 数据入库时间. 目前, 仅仅展示 `unit_type=[1, 3]` 的, 跟凯强确认了下 该字段存在三种取值 `1 普通资产单元, 2 默认资产单元, 3 客户资产单元`.

- ads_unit_label_value 表: deal_date 日期, <u>au_code 资产单元</u>, label 标签, <u>value 标签内容</u>. 设置 `label = 'manager'`, `au_code - value` 就是资产单元和基金经理的映射. 需要注意的是因为存在日期的维度, 存在一个资产单元在不同的日期隶属于不同的基金经理的情况, 该情况在展示上就是每个基金经理都会有该资产单元.

- /api/uc/v1/users 接口: <u>id</u>, <u>userName</u>, nickName, email, mobile, avatar, status, ext, roles, sex, depts, qywxId, createAt. `status=0` _应该_ 是属于正常状态.

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
