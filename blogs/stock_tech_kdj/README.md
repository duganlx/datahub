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
j_i = 3 \times d_i - 2 \times k_i
$$

公式解读

`rsv` 计算的分子是 _收盘价_ 减去 _窗口内的最低价_，分母是 _窗口内的最高价_ 与 _窗口内的最低价_ 相减。分母代表窗口内最大波动范围，必然大于 0。分子收盘价 ≥ 窗口内的最低价。所以 `rsv` 的取值范围为 [0, 100]，`k` 利用 _昨日 k_ 平滑 `rsv`，`d` 利用 _昨日 d_ 平滑 `k`。平滑则会更加稳定，也让其更加钝感。

Q1: 如果是这样理解，那么 rsv 始终都为正值，k、d 也会一直都是 > 50 的值?

代码实现如下

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

### 【知乎】如何正确理解 KDJ

知乎话题：https://www.zhihu.com/question/27652388

于渊观点 共三篇：

- [kdj 指标详解（一）](https://zhuanlan.zhihu.com/p/340129258)
- [kdj 指标详解（二）](https://zhuanlan.zhihu.com/p/340277068)
- [kdj 指标详解（三）](https://zhuanlan.zhihu.com/p/341339283)

移动速度（快 -> 慢）/敏感性（强 -> 弱）/稳定性（差 -> 好）：J 线 > K 线 > D 线

当 k、d、j 三个值都在 50 附近时，表示多空力量比较均衡

当 k、d、j 三个值都大于 50 时，表示多方力量占优

当 k、d、j 三个值都小于 50 时，表示空方力量占优。
