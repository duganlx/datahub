# 股票技术指标

## 计算

### MTM

动量指标又叫 MTM 指标，其英文全称是“Momentum Index”，其中涉及两个变量 n 和 m，计算如下

$$
mtm_i = close_i - close_{i-n+1}
$$

$$
mtmma_i = \frac{mtm_{i-m+1} + mtm_{i-m+2} + ... + mtm_{i}}{m}
$$

### RSI

相对强弱指数，是根据一定时期内上涨点数和涨跌点数之和的比率制作出的一种技术曲线。能够反映出市场在一定时期内的景气程度。

假设收盘价 `close={p1,p2,...,pn}`, 且计算 RSI 的时间周期为$n$天

第一天的计算：

$$
up\_avg_1 = \frac{\sum_{i=2}^{7} p_{i}-p_{i-1}}{n},\ s.t.\ p_{i}-p_{i-1}>0
$$

$$
down\_avg_1 = \frac{\sum_{i=2}^{7} p_{i-1}-p_{i}}{n},\ s.t.\ p_{i-1}-p_{i}>0
$$

$$
rsi_1 = \frac{up\_avg_1}{up\_avg_1 + down\_avg_1} \times 100
$$

后续 rsi 的计算如下（伪代码描述）：

```txt
up = 0 # 本轮的涨幅
down = 0 # 本轮的跌幅

如果当日收盘价格 > 昨日收盘价:
    up = close[j] - close[j-1]
    down = 0
否则:
    up = 0
    down = close[j-1] - close[j]

up_avg = (上一轮的up_avg * (n - 1) + up) / n
down_avg = (上一轮的down_avg * (n - 1) + down) / n

第i个rsi = up_avg / (up_avg + down_avg) * 100
```

### DMI

动向指数 Directional Movement Index, DMI, 动向指数又叫移动方向指数或趋向指数。是属于趋势判断的技术指标，其基本原理是通过分析股票价格在上升及下跌过程中供需关系的均衡点，即供需关系受价格变动之影响而发生由均衡到失衡的循环过程，从而提供对趋势判断的依据。

计算公式如下

```txt
TR = SUM(MAX(HIGH-LOW, |HIGH-昨日CLOSE|, |LOW-昨日CLOSE|), N);
HD = HIGH - 昨日HIGH;
LD = 昨日LOW - LOW;
DMP = SUM(IF(HD>0 && HD>LD, HD, 0), N);
DMM = SUM(IF(LD>0 && LD>HD, LD, 0), N);
PDI = DMP/TR * 100;
MDI = DMM/TR * 100;
ADX = MA(|MDI-PDI|/(MDI+PDI) * 100, M);
ADXR = (ADX + REF(ADX,M)) / 2;
```

### DMA

DMA（Different of Moving Average，平行线差）是目前股市分析技术指标中的一种中短期指标，它常用于大盘指数和个股的研判。DMA 指标是属于趋向类指标，也是一种趋势分析指标。

有三个参数 short, long, m，计算公式如下

```txt
DDD: MA(CLOSE, short) - MA(CLOSE, long)
AMA: MA(DDD, m)
```

### BRAR

人气指标 AR 和意愿指标 BR 都是以分析历史股价为手段的技术指标，其中人气指标较重视开盘价，从而反映市场买卖的人气，而意愿指标则重视收盘价格，反映的是市场买卖意愿的程度，两项指标分别从不同的角度对股价波动进行分析，达到追踪股价未来动向的目的。

计算公式如下

```text
AR: SUM(HIGH - OPEN, N) / SUM(OPEN - LOW, N) * 100
BR: SUM(MAX(0, HIGH - 昨日CLOSE), N) / SUM(MAX(0, 昨日CLOSE - LOW), N) * 100
```

### OBV

能量潮指标（On Balance Volume, OBV）是葛兰维 Joe Granville 于本世纪 60 年代提出的，并被广泛使用。股市技术分析的四大要素：价、量、时、空。OBV 指标就是从"量" 这个要素作为突破口，来发现热门股票、分析股价运动趋势的一种技术指标。它是将股市的人气——成交量与股价的关系数字化、直观化，以股市的成交量变化来衡量股市的推动力，从而研判股价的走势。关于成交量方面的研究，OBV 能量潮指标是一种相当重要的分析指标之一。

`TOTAL_BARS_COUNT`为自设的基准值，计算公式如下

```text
VA := IF(CLOSE>昨日CLOSE, VOL, IF(CLOSE<昨日CLOSE, -VOL, 0))
OBV: SUM(VA, TOTAL_BARS_COUNT)
```
