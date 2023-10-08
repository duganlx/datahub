# 股票技术指标 MACD

平滑异同平滑平均线（Moving Average Convergence Divergence，MACD）Geral Appel 于 1979 年提出的，它是一项利用短期（常用为 12 日）移动平均线与长期（常用为 26 日）移动平均线之间的聚合与分离状况，对买进、卖出时机做出研判的技术指标。在现有的技术分析软件中，MACD 常用参数是快速平滑移动平均线为 12，慢速平滑移动平均线参数为 26. 此外，MACD 还有一个辅助指标柱状图（Bar）。在大多数期货技术分析软件中，柱状线是有颜色的，在低于 0 轴以下是绿色，高于 0 轴以上是红色，前者代表趋势较弱，后者代表趋势较强。

假设收盘价$CLOSE=\{p_1,p_2,...,p_n\}$

$$
ema_i(12) = \frac{2}{13} \times p_i + \frac{11}{13} \times ema_{i-1}(12), ema_1=p_1
$$

$$
ema_i(26) = \frac{2}{27} \times p_i + \frac{25}{27} \times ema_{i-1}(26), ema_1=p_1
$$

$$
dif_i = ema_i(12) - ema_i(26)
$$

$$
dea_i = \frac{2}{10} \times dif_i + \frac{8}{10} \times dea_{i-1}, dea_1=dif_1
$$

$$
macd_i = (dif_i - dea_i) \times 2
$$

代码实现

```py
# MACD 代码实现
import pandas as pd
import numpy as np

def calc_macd(close) -> pd.DataFrame:
    """
    计算平滑异同平均线 MACD
    """
    ema12, ema26 = [], []
    for ele in close.array:
        if len(ema12) == 0:
            ema12.append(ele)
            ema26.append(ele)
            continue

        pre_ema12, pre_ema26 = ema12[-1], ema26[-1]
        cur_ema12 = 2/13 * ele + 11/13 * pre_ema12
        cur_ema26 = 2/27 * ele + 25/27 * pre_ema26

        ema12.append(cur_ema12)
        ema26.append(cur_ema26)

    ema12 = pd.Series(ema12, name='EMA12')
    ema26 = pd.Series(ema26, name='EMA26')
    dif = pd.Series(ema12 - ema26, name='DIF')

    dea = []
    for ele in dif.array:
        if len(dea) == 0:
            dea.append(ele)
            continue

        pre_dif = dea[-1]
        cur_dea = 2/10 * ele + 8/10 * pre_dif
        dea.append(cur_dea)

    dea = pd.Series(dea, name='DEA')
    macd = pd.Series((dif - dea) * 2, name='MACD')

    df = pd.concat([dif, dea, macd], axis=1)
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
```

## 观点

1. MACD 及 DIF 均为正值, 可视为多头市场; MACD 及 DIF 均为负值, 可视为空头市场
2. DIF 向上突破 MACD, 为买进讯号; DIF 向下跌破 MACD, 为卖出讯号
3. DIF 值由负转正, 且穿越 MACD, 为买进讯号; DIF 值由正转负, 且突破 MACD, 为卖出讯号
4. 如果 MACD 及 DIF 皆为正值, 且 DIF 向上突破 MACD, 此为买方市场, 做多较有利
5. 如果 MACD 及 DIF 皆为负值, 且 DIF 向下跌破 MACD, 此为卖方市场, 做空较有利
6. DIF 与大盘指数呈背离走势时, 若股价连续创新低点, 而 DIF 值并未创新低点, 此为"正背离"走势, 为买进时机; 反之,
   若股价连续创新高点, 而 DIF 值并未创新高点, 此为"负背离"走势, 为卖出时机
7. 当 DIF 和 DEA 处于 0 轴以上时，属于多头市场，DIF 线自下而上穿越 DEA 线时是买入信号。
8. 当 DIF 和 DEA 处于 0 轴以下时，属于空头市场。DIF 线自上而下穿越 DEA 线时是卖出信号，DIF 线自下而上穿越 DEA 线时，如果两线值还处于 0 轴以下运行，仅仅只能视为一次短暂的反弹，而不能确定趋势转折，此时是否买入还需要借助其他指标来综合判断。
9. 柱状线收缩和放大。一般地说，柱状线的持续收缩表明趋势运行的强度正在逐渐减弱，当柱状线颜色发生改变时，趋势确定转折。
10. 形态和背离情况。MACD 指标也强调形态和背离现象。当形态上 MACD 指标的 DIF 线与 MACD 线形成高位看跌形态，如头肩顶、双头等，应当保持警惕；而当形态上 MACD 指标 DIF 线与 MACD 线形成低位看涨形态时，应考虑进行买入。在判断形态时以 DIF 线为主，MACD 线为辅。当价格持续升高，而 MACD 指标走出一波比一波低的走势时，意味着顶背离出现，预示着价格将可能在不久之后出现转头下行，当价格持续降低，而 MACD 指标却走出一波高于一波的走势时，意味着底背离现象的出现，预示着价格将很快结束下跌，转头上涨。
11. 黄金交叉: 快线(DIF)上穿慢线(DEA)
12. 死亡交叉: 快线(DIF)向下跌破慢线(DEA)
