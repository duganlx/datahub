# 股票技术指标 MA

移动平均线, 可分为简单移动平均线 SMA 和 指数移动平均线 EMA, 假设收盘价 `close = {p1, p2, ..., pn}`, n 日的移动平均线计算公式如下

$$
sma_i = \frac{p_{i-n+1} + p_{i-n+2} + ... + p_{i}}{n}
$$

$$
ema_i = \frac{2}{n+1} \times p_i + \frac{n-1}{n+1} \times ema_{i-1}
$$

代码实现如下

```py
import pandas as pd

def calc_sma(series, window) -> pd.Series:
    """
    计算简单移动平均线 MA
    """
    ma = series.rolling(window=window).mean()
    ma.name = f'MA{window}'

    return ma

def calc_expma(close, n):
    """
    计算指数移动平均线 EXPMA
    """
    expma = []
    for ele in close.array:
        if len(expma) == 0:
            expma.append(ele)
            continue

        prev_expma = expma[-1]
        alpha = 2 / (n+1)
        cur_expma = alpha * ele + (1-alpha) * prev_expma

        expma.append(cur_expma)

    expma = pd.Series(expma, name=f'EXPMA{n}')
    return expma
```
