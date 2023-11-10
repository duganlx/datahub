"""
定义测试的执行过程
"""
import logging
import multiprocessing
import os
import time

from dtest.supporter import BizError, Constant, Rule, Utils, Report


class Master(object):
    """
    执行器
    """
    os_type = ''
    cpu_core_sum = 0
    task_sum = 0

    worker = None
    logger = None

    # 评价指标
    run_time = 0
    average_response_time = 0
    response_num_pre_second = 0
    min_response_time = Constant.INVALID_RESPONSE_MIN_TIME
    max_response_time = Constant.INVALID_RESPONSE_MAX_TIME
    success_task_num = 0
    error_task_num = 0

    def __init__(self, cpu_core_sum, task_sum, logger_mode, task_cls):

        os_type = Utils.detect_os_type()

        logger_mode_config = None
        if logger_mode == Constant.DEBUG_MODE:
            logger_mode_config = Constant.DEV_MODE_CONFIG1
            logging.basicConfig(**logger_mode_config)
        else:
            logger_mode_config = Constant.INFO_MODE_CONFIG
            logging.basicConfig(**logger_mode_config)

        self.logger = logging.getLogger(__name__)
        self.logger.debug(f'日志模式为: {logger_mode}')

        self.logger.info(f'Actuator开始初始化...')

        self.os_type = os_type
        self.cpu_core_sum = cpu_core_sum
        self.task_sum = task_sum
        self.worker = Worker(self.logger, task_cls, os_type, logger_mode_config=logger_mode_config)
        self.success_task_num = task_sum

        self.logger.debug(f'operating system: {self.os_type}')
        self.logger.debug(f'cpu core sum: {self.cpu_core_sum}')
        self.logger.debug(f'task sum: {self.task_sum}')
        self.logger.debug(f'worker: {self.worker}')
        self.logger.debug(f'success task num: {self.success_task_num}')
        self.logger.info(f'{self} 完成初始化...')

    def run(self):

        mgr = multiprocessing.Manager()
        task_queue = mgr.Queue()
        report_queue = mgr.Queue()
        process_pool = multiprocessing.Pool(self.cpu_core_sum)

        self.worker.prepare_task_data(self.task_sum, task_queue)

        start = time.time()
        for i in range(self.cpu_core_sum):
            process_pool.apply_async(self.worker.job, (task_queue, report_queue))

        process_pool.close()
        process_pool.join()
        self.run_time = time.time() - start

        self.logger.info('所有任务执行完毕...')
        self.evaluate(report_queue)
        self.output_evaluation()

    def evaluate(self, report_queue):
        while not report_queue.empty():
            report = report_queue.get_nowait()

            self.logger.debug(
                f'收到进程[{report.worker_id}]的报告: ' +
                f'最快响应时间为{report.min_time}, ' +
                f'最慢响应时间为{report.max_time}, ' +
                f'错误数量为{report.error_task_num}'
            )

            self.success_task_num -= report.error_task_num
            self.error_task_num += report.error_task_num
            self.max_response_time = max(self.max_response_time, report.max_time)
            self.min_response_time = min(self.min_response_time, report.min_time)

        self.logger.info(f'完成所有进程的报告统计, '
                         f'最快响应时间为{self.min_response_time}, '
                         f'最慢响应时间为{self.max_response_time}, '
                         f'总共错误数量为{self.error_task_num}')

        self.average_response_time = Rule.average_response_time(run_time=self.run_time,
                                                                task_sum=self.success_task_num)

        self.response_num_pre_second = Rule.response_num_per_second(run_time=self.run_time,
                                                                    task_sum=self.success_task_num)
        self.min_response_time = Rule.response_min_time(self.min_response_time)
        self.max_response_time = Rule.response_max_time(self.max_response_time)

    def output_evaluation(self):
        self.logger.info("===============压测结果===================")
        self.logger.info(f'进程数量: {self.cpu_core_sum}')
        self.logger.info(f'任务数量: {self.task_sum}')
        self.logger.info(f'总耗时(秒): {self.run_time}')
        self.logger.info(f'平均响应时间(秒): {self.average_response_time}')
        self.logger.info(f'每秒请求数（吞吐率）: {self.response_num_pre_second}')
        self.logger.info(f'最短响应时间(秒): {self.min_response_time}')
        self.logger.info(f'最长响应时间(秒): {self.max_response_time}')
        self.logger.info(f'任务成功数量：{self.success_task_num}')
        self.logger.info(f'任务失败数量：{self.error_task_num}')
        self.logger.info("========================================")


class Worker(object):
    """
    打工人 / 子进程
    """
    logger = None
    task = None
    os_type = None
    logger_mode_config = None

    def __init__(self, logger, task_cls, os_type, **kwargs):
        self.logger = logger
        self.logger.info('Task开始初始化...')
        self.logger.info(f'{self} 完成初始化...')
        self.logger_mode_config = kwargs['logger_mode_config']

        self.os_type = os_type
        self.task = task_cls(logger)

    def job(self, task_queue, report_queue):
        """
        进程实际执行的方法
        """
        if self.os_type == Constant.WIN32:
            # 在Windows环境下，重新启动一个进程时，logger的日志级别会变成warn（例:<Logger actuator (WARNING)>）
            # 而非从Master中传过来正确的日志级别（例: <Logger actuator (DEBUG)>）
            logging.basicConfig(**self.logger_mode_config)
            self.logger = logging.getLogger(__name__)

        self.logger.debug('Worker 开始工作...')
        self.logger.debug(f'task queue: {task_queue}')
        self.logger.debug(f'report queue: {report_queue}')

        pid = os.getpid()
        report = Report(pid)

        self.task.prepare_worker_status()

        while not task_queue.empty():
            try:
                task_data = task_queue.get_nowait()

                start = time.time()

                self.task.run(**task_data)

                run_time = time.time() - start
                report.min_time = min(run_time, report.min_time)
                report.max_time = max(run_time, report.max_time)

            except BizError as e:
                report.error_task_num += 1
                # 业务错误 等级应该是info, debug为具体参数变量等信息打印
                self.logger.info(f'{e}')

        self.logger.debug('Worker 结束工作...')
        self.logger.debug(report.__dict__)
        report_queue.put(report)

    def prepare_task_data(self, task_sum, task_queue):
        """
        支持自定义放入任务队列(task_queue)中的数据单元，该数据单元可以在worker真正执行的时候取得
        """
        self.task.prepare_task_data(task_sum, task_queue)
        self.logger.debug(f'需要执行的任务数量为 {task_queue.qsize()}')


class Task(object):
    logger = None

    def __init__(self, logger):
        self.logger = logger

    def run(self, **kwargs):
        """
        需要测试的函数方法，如果不符合业务逻辑可以抛出BizError错误(包含其子类)，最终会被记录到报告(Report)中
        """
        pass

    def prepare_task_data(self, task_sum, task_queue):
        """
        默认放入任务队列(task_queue)的数据单元就是自增的task_id
        """
        for i in range(task_sum):
            task_data = {'task_id': i}
            task_queue.put(task_data)
            self.logger.debug(f'将任务数据：[{task_data}] 放入task_queue中')

    def prepare_worker_status(self):
        """
        worker在执行task_queue之前需要做的准备工作
        """
        pass
