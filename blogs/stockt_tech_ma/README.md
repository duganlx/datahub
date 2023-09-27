# MA

移动平均线, 可分为简单移动平均线 SMA 和 指数移动平均线 EMA, 假设收盘价 `close = {p1, p2, ..., pn}`, n 日的移动平均线计算公式如下

$$
sma_i = \frac{p_{i-n+1} + p_{i-n+2} + ... + p_{i}}{n}
$$

$$
ema_i = \frac{2}{n+1} \times p_i + \frac{n-1}{n+1} \times ema_{i-1}
$$


