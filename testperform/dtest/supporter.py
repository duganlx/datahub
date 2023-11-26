"""
supporter 中定义支援工具
"""
import logging
import sys


class Report(object):

    def __init__(self, pid):
        self.worker_id = pid
        self.min_time = Constant.INVALID_RESPONSE_MIN_TIME
        self.max_time = Constant.INVALID_RESPONSE_MAX_TIME
        self.error_task_num = 0


class Constant(object):
    """
    常量
    """
    WIN32 = 'win32'  # wins 操作系统
    LINUX = 'linux'  # linux 操作系统

    INVALID_RESPONSE_MAX_TIME = -1  # 无效的最大响应时间
    INVALID_RESPONSE_MIN_TIME = sys.maxsize  # 无效的最小响应时间

    DEBUG_MODE = 'debug'
    INFO_MODE = 'info'
    DEV_MODE_CONFIG1 = {'level': logging.DEBUG,
                        'format': '%(levelname)s | pid[%(process)d]: ' +
                                  '%(message)s'}

    DEV_MODE_CONFIG2 = {'level': logging.DEBUG,
                        'format': '%(levelname)s | %(filename)s %(funcName)s:%(lineno)d' +
                                  '| pid[%(process)d] tid[%(thread)d] ' +
                                  '| %(message)s'}
    INFO_MODE_CONFIG = {'level': logging.INFO, 'format': '%(message)s'}


class BizError(Exception):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return repr(self.value)


class Utils(object):

    @staticmethod
    def detect_os_type():
        os_type = sys.platform

        if os_type == Constant.LINUX:
            return Constant.LINUX
        elif os_type == Constant.WIN32:
            return Constant.WIN32
        else:
            raise Exception(f'不支持的操作系统类型：{os_type}')


class Rule(object):
    @staticmethod
    def average_response_time(run_time, task_sum):
        """
        统计平均响应时间
        :param run_time: 运行时间（秒）
        :param task_sum: 任务数量（个）
        :return: 正常返回数值类型，否则返回字符串 'NAN'
        """
        if task_sum == 0:
            return 'NAN'

        return run_time / task_sum

    @staticmethod
    def response_num_per_second(run_time, task_sum):
        """
        统计每秒请求数 / 吞吐量
        :param run_time: 运行时间（秒）
        :param task_sum: 任务数量（个）
        :return: 正常返回数值类型，否则返回字符串 'NAN'
        """
        val = Rule.average_response_time(run_time, task_sum)

        if isinstance(val, str):
            return 'NAN'

        return 1 / val

    @staticmethod
    def response_min_time(min_time):
        """
        验证最短响应时间
        :param min_time: 时间
        :return: 如果 min_time 为合法值，则返回该值，否则 返回 'NAN'
        """
        if min_time == Constant.INVALID_RESPONSE_MIN_TIME:
            return 'NAN'

        return min_time

    @staticmethod
    def response_max_time(max_time):
        """
        验证最长响应时间
        :param max_time: 时间
        :return: 如果 max_time 为合法值，则返回该值，否则 返回 'NAN'
        """
        if max_time == Constant.INVALID_RESPONSE_MAX_TIME:
            return 'NAN'

        return max_time
