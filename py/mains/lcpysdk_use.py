"""
lcpysdk test
"""
import os
import sys

# 解决 ModuleNotFoundError: No module named 'xxx' 问题
this_file_full_path_name = os.path.abspath(__file__)
this_file_folder_path = os.path.dirname(this_file_full_path_name)
parent_folder_path = os.path.dirname(this_file_folder_path)
sys.path.append(parent_folder_path)

from lcpysdk.index import LocalPySdk

if __name__ == '__main__':
  lcsdk = LocalPySdk()

  df = lcsdk.getData(dir='tmp', filename='xxx.csv')
  print(df)
