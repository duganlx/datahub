import pandas as pd
import numpy as np

# 简单移动平均线
def sma(series, window) -> pd.Series:
    ma = series.rolling(window=window).mean()
    ma.name = f'MA{window}'

    return ma

# 指数移动平均线
def expma(close, n):
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

# macd
def macd(close) -> pd.DataFrame:
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

# kdj
def kdj(close, high, low) -> pd.DataFrame:
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