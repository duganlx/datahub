# 股票技术指标 KDJ

随机指标（KDJ）由 George C. Lane 创建. 它综合了动量观念、强弱指标及移动平均线的优点, 用来度量股价脱离价格正常范围的变异程度. 该指标总共包含三根线, 分别是 K, D, J 线

假设收盘价 `CLOSE = {close_1, close_2,..., close_n}`，最高价 `HIGH = {high_1, high_2, ..., high_n}`，最低价 `LOW = {low_1, low_2, ..., low_n}`，则 kdj 计算如下

$$
lowest\_low_i = min\{low_{i-8},\ low_{i-7},\ ...,\ low_{i}\}
$$

$$
highest\_high_i = max\{high_{i-8},\ high_{i-7},\ ...,\ high_{i}\}
$$

$$
rsv_i = \frac{close_i - lowest\_low_i}{highest\_high_i - lowest\_low_i} \times 100
$$

$$
k_i = \frac{2}{3} \times k_{i-1} + \frac{1}{3} \times rsv_i,\ k_1=50
$$

$$
d_i = \frac{2}{3} \times d_{i-1} + \frac{1}{3} \times k_i,\ d_1=50
$$

$$
j_i = 3 \times k_i - 2 \times d_i
$$

代码实现

```py
# KDJ 代码实现
import pandas as pd

def calc_kdj(close, high, low) -> pd.DataFrame:
    """
    计算 KDJ
    """
    lowest_low = low.rolling(window=9, min_periods=1).min()
    highest_low = high.rolling(window=9, min_periods=1).max()
    rsv = (close - lowest_low) / (highest_low - lowest_low) * 100
    rsv = pd.Series(rsv, name='RSV')

    k, d, j = [], [], []
    for ele in rsv.array:
        if len(k) == 0:
            # 初始值50
            k.append(50)
            d.append(50)
            j.append(50)
            continue

        pre_k = k[-1]
        cur_k = (2 * pre_k + ele) / 3

        pre_d = d[-1]
        cur_d = (2 * pre_d + cur_k) / 3

        cur_j = 3 * cur_k - 2 * cur_d

        k.append(cur_k)
        d.append(cur_d)
        j.append(cur_j)

    k = pd.Series(k, name='K')
    d = pd.Series(d, name='D')
    j = pd.Series(j, name='J')

    df = pd.concat([k, d, j], axis=1)

    return df
```

## 研究

rsv 计算的分子是 _收盘价_ 减去 _窗口内的最低价_，分母是 _窗口内的最高价_ 与 _窗口内的最低价_ 相减。分母代表窗口内最大波动范围，必然大于 0。分子收盘价 ≥ 窗口内的最低价。所以 rsv 的取值范围为 [0, 100]，k 利用 _昨日 k_ 平滑 rsv，d 利用 _昨日 d_ 平滑 k。平滑则会更加稳定，也让其更加钝感。

在 k1=50 的前提下，k2 要想大于 50，需要 rsv 大于 50，而 rsv 取值范围是[0, 100]，所以 k、d 必定是一个大于 0 的值，但波动范围也是 [0, 100]。

因为 d 的计算是 k 之上再做一次平滑，所以可以理解为 "慢线"，k 相对即为 "快线"。j 是 三倍的快线减去两倍的慢线，当快速下跌时，j 可能会小于 0。

不过 rsv、k、d 必定是 > 0 的。由于 d 的计算是 k 之上再做一次平滑，所以可以理解为 "慢线"，k 相对即为 "快线"。而 j 是 三倍的快线减去两倍的慢线，考虑最极端的情况，j 最小取值为 -200（`k=0, d=100`），最大取值为 300（`k=100, d=0`）。

当股价快速拉升时，快线 k 的值会变得很大，而慢线 d 会迟钝一点，那么此时 j 值可能会超过 100。拿贵州茅台 2023-07-25 举例子，数据如下所示。_窗口内 min(最低价)_ 为 `1713.80`， _窗口内 max(最高价)_ 为 `1828.88`。而当天收盘价为 `1828.55`，与 _窗口内 max(最高价)_ 仅仅有 3.3 毛钱的缺口。在这种情况下，rsv 就会非常大（`99.71`，`114.75/115.08`）。快线 k 在 _昨日 k_ 的基础上 加上 `1/3` 的 rsv，慢线 d 在 _昨日 d_ 的基础上 加上 `1/3` 的 k（即 对于 rsv 而言，仅仅取了 `1/9` ）。 j 用 3 倍快线（`263.71`）减去 2 倍慢线（`142.9`）等于 `110.80`。

```
   trade_date          K          D           J    close  lowest_low  highest_high        RSV
14 2023-07-21  66.595025  58.846389   82.092295  1771.30     1702.10       1772.50  98.295455
15 2023-07-24  76.994362  64.895713  101.191658  1771.15     1711.33       1772.50  97.793036
16 2023-07-25  84.567322  71.452916  110.796134  1828.55     1713.80       1828.88  99.713243
17 2023-07-26  87.681922  76.862585  109.320596  1828.55     1713.80       1835.99  93.911122
18 2023-07-27  87.825492  80.516887  102.442702  1838.03     1713.80       1854.79  88.112632
```

todo
当股价快速

如何正确理解 KDJ

知乎话题：https://www.zhihu.com/question/27652388

于渊观点 共三篇：

- [kdj 指标详解（一）](https://zhuanlan.zhihu.com/p/340129258)
- [kdj 指标详解（二）](https://zhuanlan.zhihu.com/p/340277068)
- [kdj 指标详解（三）](https://zhuanlan.zhihu.com/p/341339283)

移动速度（快 -> 慢）/敏感性（强 -> 弱）/稳定性（差 -> 好）：J 线 > K 线 > D 线

当 k、d、j 三个值都在 50 附近时，表示多空力量比较均衡

当 k、d、j 三个值都大于 50 时，表示多方力量占优

当 k、d、j 三个值都小于 50 时，表示空方力量占优。
