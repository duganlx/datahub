import pandas as pd
import numpy as np

class StockTA(object):
    trade_date: pd.Series = None
    open: pd.Series = None
    close: pd.Series = None
    high: pd.Series = None
    low: pd.Series = None
    pre_close: pd.Series = None
    volume: pd.Series = None  # 成交量（手）

    def __init__(self, df):
        self.trade_date = df['trade_date']
        self.open = df['open']
        self.close = df['close']
        self.high = df['high']
        self.low = df['low']
        self.pre_close = df['pre_close']
        self.volume = df['total_vol'] / 100  # 1手 = 100股

    def ma(self, wins, prod=False) -> pd.Series:
        close = self.close
        ma = close.rolling(window=wins).mean()
        ma.name = f'MA{wins}'
        if prod:
            ma.name = f'ma_MA{wins}'
        return ma

    def expma(self, n, prod=False) -> pd.Series:
        close = self.close
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
        if prod:
            expma = pd.Series(expma, name=f'expma_EMA{n}')
        return expma

    def macd(self, prod=False) -> pd.DataFrame:
        close = self.close
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
        if prod: 
            df = pd.concat([
                pd.Series(dif, name='macd_DIF'),
                pd.Series(dea, name='macd_DEA'),
                pd.Series(macd, name='macd_MACD'),
            ], axis=1)

        return df

    def kdj(self, wins=9, prod=False) -> pd.DataFrame:
        close = self.close
        high = self.high
        low = self.low

        lowest_low = low.rolling(window=wins, min_periods=1).min()
        highest_high = high.rolling(window=wins, min_periods=1).max()
        rsv = (close - lowest_low) / (highest_high - lowest_low) * 100
        # 如果分母为0, 则 rsv为0
        rsv = pd.Series(rsv, name='RSV').fillna(0)
        # lowest_low = pd.Series(lowest_low, name='lowest_low')
        # highest_high = pd.Series(highest_high, name='highest_high')

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

        if prod: 
            k = pd.Series(k, name='kdj_K')
            d = pd.Series(d, name='kdj_D')
            j = pd.Series(j, name='kdj_J')

        df = pd.concat([k, d, j], axis=1)
        return df

    def boll(self, wins=20, k=2, prod=False) -> pd.DataFrame:
        close = self.close
        mid = close.rolling(window=wins).mean()
        std = close.rolling(window=wins).std()
        upper = mid + k * std
        lower = mid - k * std

        mid = pd.Series(mid, name='mid')
        upper = pd.Series(upper, name='upper')
        lower = pd.Series(lower, name='lower')
        if prod: 
            mid = pd.Series(mid, name='boll_MID')
            upper = pd.Series(upper, name='boll_UPPER')
            lower = pd.Series(lower, name='boll_LOWER')

        df = pd.concat([mid, upper, lower], axis=1)
        return df

    def mtm(self, n=6, m=6, prod=False):
        # 动量指标
        close = self.close
        mtm = close - close.shift(n)
        mtmma = mtm.rolling(m).mean()

        df = pd.DataFrame({'MTM': mtm, 'MTMMA': mtmma})
        if prod: 
            df = pd.DataFrame({'mtm_MTM': mtm, 'mtm_MTMMA': mtmma})
        return df
    
    def rsi(self, periods=6, prod=False):
        close = self.close
        length = len(close)
        rsi = [np.nan] * length

        up_avg = 0
        down_avg = 0

        first_t = close[:periods+1]
        for i in range(1, len(first_t)):
            if first_t[i] >= first_t[i-1]:
                up_avg += first_t[i] - first_t[i-1]
            else:
                down_avg += first_t[i-1] - first_t[i]
                
        up_avg = up_avg / periods
        down_avg = down_avg / periods
        rsi[periods] = up_avg / (up_avg + down_avg) * 100

        for j in range(periods+1, length):
            up = 0
            down = 0
            if close[j] >= close[j-1]:
                up = close[j] - close[j-1]
                down = 0
            else:
                up = 0
                down = close[j-1] - close[j]

            up_avg = (up_avg * (periods - 1) + up) / periods
            down_avg = (down_avg * (periods - 1) + down) / periods
            rsi[j] = up_avg / (up_avg + down_avg) * 100

        if prod:
            return pd.Series(rsi, name=f'rsi_RSI{periods}')
        return pd.Series(rsi, name=f'RSI{periods}')

    def dmi(self, n=14, m=6, prod=False):
        close = self.close
        high = self.high
        low = self.low

        # TR: true range
        yesterday_close = close.shift()
        c1 = high - low
        c2 = np.abs(high - yesterday_close)
        c3 = np.abs(low - yesterday_close)
        tr = np.maximum(c1, np.maximum(c2, c3))
        tr = pd.Series(tr).rolling(n).sum()

        # DMP & DMM
        yesterday_high = high.shift()
        yesterday_low = low.shift()
        hd = high - yesterday_high
        ld = yesterday_low - low
        dmp = np.where((hd > 0) & (hd > ld), hd, 0)
        dmm = np.where((ld > 0) & (ld > hd), ld, 0)
        dmp = pd.Series(dmp).rolling(n).sum()
        dmm = pd.Series(dmm).rolling(n).sum()

        # PDI(+DI, 正方向指标) MDI(-DI, 负方向指标)
        pdi = dmp / tr * 100
        mdi = dmm / tr * 100

        # ADX: 平均趋向指数
        a = np.abs(mdi - pdi)
        b = pdi + mdi
        adx = a / b * 100
        adx = pd.Series(adx).rolling(m).mean()
        
        # ADXR: 趋向指数平均数
        a = adx
        b = adx.shift(m)
        adxr = (a + b) / 2

        if prod: 
            return pd.DataFrame({'dmi_PDI': pdi, 'dmi_MDI': mdi, 'dmi_ADX': adx, 'dmi_ADXR': adxr})
        return pd.DataFrame({'PDI': pdi, 'MDI': mdi, 'ADX': adx, 'ADXR': adxr})

    def dma(self, short=10, long=50, m=10, prod=False):
        close = self.close
        short_ma = close.rolling(short, min_periods=1).mean()
        long_ma = close.rolling(long, min_periods=1).mean()
        ddd = short_ma - long_ma
        ama = ddd.rolling(m, min_periods=1).mean()

        if prod:
            return pd.DataFrame({'dma_DDD': ddd, 'dma_AMA': ama})
        return pd.DataFrame({'DDD': ddd, 'AMA': ama})

    def brar(self, n=26, prod=False):
        open = self.open
        close = self.close
        high = self.high
        low = self.low
        # AR 人气指标
        a = (high - open).rolling(n).sum()
        b = (open - low).rolling(n).sum()
        ar = a / b * 100

        # br 意愿指标
        yesterday_close = close.shift()
        a = np.maximum(high - yesterday_close, 0)
        a = pd.Series(a).rolling(n).sum()
        b = np.maximum(yesterday_close - low, 0)
        b = pd.Series(b).rolling(n).sum()
        br = a / b * 100

        if prod:
            return pd.DataFrame({'brar_AR': ar, 'brar_BR': br})
        return pd.DataFrame({'AR': ar, 'BR': br})

    def obv(self, offset, verbose:bool=False, prod=False):
        close = self.close
        pre_close = self.pre_close
        volume = self.volume / 10000 # 单位: 万手

        df = pd.DataFrame({'close': close, 'pre_close': pre_close, "volume": volume})

        obv = []
        for i in range(len(df)):
            series = df.iloc[i, :]

            prev_obv = 0 if len(obv) == 0 else obv[-1]

            if series['close'] > series['pre_close']:
                cur_obv = prev_obv + series['volume']
            elif series['close'] < series['pre_close']:
                cur_obv = prev_obv - series['volume']
            else:
                cur_obv = prev_obv

            obv.append(cur_obv)

        obv = pd.Series(obv, name='OBV')
        obv = obv + offset # revise value

        maobv = obv.rolling(30, min_periods=1).mean()

        if verbose:
            return pd.DataFrame({'trade_date': self.trade_date, 'close': close, 'pre_close': pre_close, 'volume': volume ,'OBV': obv, 'MAOBV': maobv})

        if prod:
            return pd.DataFrame({'obv_OBV': obv, 'obv_MAOBV': maobv})
        return pd.DataFrame({'OBV': obv, 'MAOBV': maobv})

    def wr(self, n, prod=False):
        close = self.close
        high = self.high
        low = self.low

        rhigh = high.rolling(n, min_periods=1).max()
        rlow = low.rolling(n, min_periods=1).min()
        wr = (rhigh - close) / (rhigh - rlow) * 100

        if prod:
            return pd.Series(wr, name=f'w&r_WR{n}')
        return pd.Series(wr, name=f'WR{n}')

    def data_matrix(self):
        date = pd.Series(self.trade_date, name='date')
        open = pd.Series(self.open, name='open')
        close = pd.Series(self.close, name='close')
        high = pd.Series(self.high, name='high')
        low = pd.Series(self.low, name='low')
        volume = pd.Series(self.volume, name='volume')
        ma5 = self.ma(wins=5, prod=True)
        ma10 = self.ma(wins=10, prod=True)
        ma20 = self.ma(wins=20, prod=True)
        ma60 = self.ma(wins=60, prod=True)
        ema5 = self.expma(n=5, prod=True)
        ema10 = self.expma(n=10, prod=True)
        ema20 = self.expma(n=20, prod=True)
        ema60 = self.expma(n=60, prod=True)
        macd = self.macd(prod=True)
        kdj = self.kdj(prod=True)
        boll = self.boll(prod=True)
        rsi = self.rsi(prod=True)
        dmi = self.dmi(prod=True)
        dma = self.dma(prod=True)
        brar = self.brar(prod=True)
        obv = self.obv(offset=32.352-815.769, prod=True)
        wr = self.wr(n=10, prod=True)
        mtm = self.mtm(prod=True)

        bilabel = pd.Series(self.close > self.pre_close, name='bilabel')

        df = pd.concat([
            # basic(6)
            date, open, close, high, low, volume,
            # MA(4)
            ma5, ma10, ma20, ma60,
            # EMA(4)
            ema5, ema10, ema20, ema60,
            # MACD(3), KDJ(3), BOLL(3), RSI(1), DMI(4), DMA(2), BRAR(2)
            macd, kdj, boll, rsi, dmi, dma, brar,
            # OBV(2), WR(1), MTM(2)
            obv, wr, mtm,
            # 二分类标签 涨:True, 跌:False
            bilabel
        ], axis=1)

        return df

    def _rolling_df(df, window, func, name):
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
    
