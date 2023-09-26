import pandas as pd


def sma(series, window) -> pd.Series:
    ma = series.rolling(window=window).mean()
    ma.name = f'MA{window}'

    return ma
