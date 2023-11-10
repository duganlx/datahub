"""
dtest demo
模拟场景1：对股票进行下单
"""
from random import random
import time
import os
import sys
from random import random
import time

this_file_full_path_name = os.path.abspath(__file__)
this_file_folder_path = os.path.dirname(this_file_full_path_name)
parent_folder_path = os.path.dirname(this_file_folder_path)
sys.path.append(parent_folder_path)

from dtest.supporter import Constant, BizError
from dtest.actuator import Master, Task

# stocks = [('股票代码', '市场', '最低点', '最高点'), ...]
stocks = [('000056', 'SZ', 4.46, 5.45), ('000062', 'SZ', 12.09, 14.77), ('000078', 'SZ', 2.99, 3.65),
          ('000096', 'SZ', 8.0, 9.78), ('000099', 'SZ', 6.8, 8.31), ('000100', 'SZ', 3.88, 4.74),
          ('000151', 'SZ', 10.73, 13.11), ('000153', 'SZ', 7.62, 9.32), ('000416', 'SZ', 3.59, 4.39),
          ('000422', 'SZ', 17.75, 21.69), ('000506', 'SZ', 2.99, 3.65), ('000532', 'SZ', 10.07, 12.31),
          ('000543', 'SZ', 4.69, 5.73), ('000595', 'SZ', 6.5, 7.94), ('000636', 'SZ', 15.52, 18.96),
          ('000688', 'SZ', 18.05, 22.06), ('000691', 'SZ', 4.3, 5.26), ('000712', 'SZ', 13.0, 15.88),
          ('000726', 'SZ', 7.01, 8.57), ('000757', 'SZ', 4.02, 4.92), ('000798', 'SZ', 9.53, 11.65),
          ('000800', 'SZ', 7.63, 9.33), ('000819', 'SZ', 20.43, 24.97), ('000825', 'SZ', 4.39, 5.37),
          ('000868', 'SZ', 5.23, 6.39), ('000893', 'SZ', 31.19, 38.13), ('000905', 'SZ', 7.58, 9.26),
          ('000909', 'SZ', 8.98, 10.98), ('000910', 'SZ', 7.99, 9.77), ('000928', 'SZ', 5.61, 6.85)]
username = 'username'
password = 'password'
host = '127.0.0.1'
port = '10000'

glo_cpu_core_sum = 3
glo_task_sum = 10
glo_logger_mode = Constant.DEBUG_MODE


class MyTask(Task):

    def run(self, **kwargs):
        self.logger.info(f'取得任务数据为' +
                         f'task_id={kwargs["task_id"]}, ' +
                         f'stock_id={kwargs["stock_id"]}, ' +
                         f'market={kwargs["market"]}, '
                         f'down_limit={kwargs["down_limit"]}, '
                         f'up_limit={kwargs["up_limit"]}')

        # 调用待测试的方法 并将参数传递
        # 通过返回值来判断是否出错. 如果出错，则通过抛出BizError及其子类
        # 这里用随机函数模拟是否下单成功
        run_time = random() * 2
        time.sleep(run_time)
        self.logger.debug(f'place order success, order id = {round(random() * 100)}')

        if run_time > 1:
            raise BizError("place order failed: time out")

    def prepare_task_data(self, task_sum, task_queue):
        # 总任务需要分为两半，一半为买，一半为卖
        two_phase_task_sum = round(task_sum / 2)

        for i in range(two_phase_task_sum):
            stock = stocks[i % len(stocks)]
            # 推荐放字典(dict)格式
            task_data = {
                'task_id': i,
                "stock_id": stock[0],
                'market': stock[1],
                'down_limit': stock[2],
                'up_limit': stock[3],
            }
            task_queue.put(task_data)
            self.logger.debug(f'将任务数据：[{task_data}] 放入task_queue中')

    def prepare_worker_status(self):
        """
        worker在执行task_queue之前需要做的准备工作
        """
        self.logger.debug(f'连接服务器{host}:{port}...')
        self.logger.debug(f'使用用户名[{username}], 密码[{password}] 登陆到服务器...')


if __name__ == '__main__':
    a = Master(glo_cpu_core_sum, glo_task_sum, glo_logger_mode, MyTask)
    a.run()
