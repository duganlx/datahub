"""
eampysdk test
"""
import os
import sys
import pandas as pd
import torch
import numpy as np

# 解决 ModuleNotFoundError: No module named 'xxx' 问题
this_file_full_path_name = os.path.abspath(__file__)
this_file_folder_path = os.path.dirname(this_file_full_path_name)
parent_folder_path = os.path.dirname(this_file_folder_path)
sys.path.append(parent_folder_path)

from utils.eampysdk import EamPySdk
from utils.localpysdk import LocalPySdk
from utils.ta import StockTA

def stockTAanalysis(local: bool = False):
    if local:
        lcsdk = LocalPySdk()
        df = lcsdk.getData(dir='tmp', filename='raw.csv')
    else:
        eamsdk = EamPySdk()
        df = eamsdk.getBardayData(
            universe=['600519.SH'],
            # where='trade_date > \'2023-07-01\''
        )
        # eamsdk.savedf(df, dir='tmp', filename='raw.csv')

    stockTa = StockTA(df)

    # ma = stockTa.ma(5)
    # ema = stockTa.expma(5)
    # macd = stockTa.macd()
    # kdj = stockTa.kdj()
    # boll = stockTa.boll()
    # mtm = stockTa.mtm()
    # rsi = stockTa.rsi()
    # dmi = stockTa.dmi()
    # dma = stockTa.dma()
    # brar = stockTa.brar()
    # obv = stockTa.obv(offset=32.352-815.769, verbose=True)
    wr = stockTa.wr(n=10)

    print(wr)

def biclassify(generate=False):
    if generate:
        eamsdk = EamPySdk()
        df = eamsdk.getBardayData(
            universe=['600519.SH'],
            where='trade_date > \'2023-07-01\''
        )
        stockTa = StockTA(df)
        dmatrix = stockTa.data_matrix()
        eamsdk.savedf(df=dmatrix, dir='tmp', filename='ta.csv')
        df = dmatrix
    else:
        lcsdk = LocalPySdk()
        df = lcsdk.getData(dir='tmp', filename='ta.csv')

    # df format
    # == begin ==
    label = df.iloc[:, -1]
    date = df.iloc[:, 0]
    df.drop(df.columns[[0, -1]], axis=1, inplace=True)

    n = df.shape[0]
    k = 5
    height = n // k
    # print(n, k, height)
    small_dfs = np.split(df.values, [i * height for i in range(1, k)], axis=0)
    arr = np.array(small_dfs)
    tensor = torch.tensor(arr)
    # == end ==

    # print(df)
    rnn = torch.nn.LSTM(input_size=36, hidden_size=20, num_layers=2, bidirectional=False)
    h0 = torch.randn(4, 5, 20) #(num_layers,batch,output_size)
    c0 = torch.randn(4, 5, 20) #(num_layers,batch,output_size)
    output, (hn, cn) = rnn(tensor, (h0, c0))

    print(output)
    print(hn, cn)


if __name__ == '__main__':
    pass
    # stockTAanalysis(local=False)

    # biclassify(generate=False)
    # rnn = torch.nn.LSTM(input_size=10, hidden_size=20, num_layers=2, bidirectional=True)
    # input = torch.randn(5, 3, 10)#(seq_len, batch, input_size)
    # h0 = torch.randn(4, 3, 20) #(num_layers,batch,output_size)
    # c0 = torch.randn(4, 3, 20) #(num_layers,batch,output_size)
    # output, (hn, cn) = rnn(input, (h0, c0))

    # print(output)
    # print(hn, cn)