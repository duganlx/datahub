"""
dtest demo
基本使用
"""
import os
import sys
from random import random
import time

# 解决 ModuleNotFoundError: No module named 'dtest' 问题
# 需要将dtest的路径添加到环境变量之中
# this_file_full_path_name = os.path.abspath(__file__)
# this_file_folder_path = os.path.dirname(this_file_full_path_name)
# parent_folder_path = os.path.dirname(this_file_folder_path)
# sys.path.append(parent_folder_path)

from dtest.actuator import Master, Task
from dtest.supporter import Constant
from dtest.supporter import BizError


class MyTask(Task):

    def run(self, **kwargs):
        run_time = random() * 2
        time.sleep(run_time)
        self.logger.debug(f'run time = {run_time}')

        if run_time > 1:
            raise BizError("接口异常")


if __name__ == '__main__':
    # 使用 if __name__ == '__main__' 防止在windows中出现以下错误：
    # An attempt has been made to start a new process before the
    #  current process has finished its bootstrapping phase.
    a = Master(3, 10, Constant.DEBUG_MODE, MyTask)
    a.run()