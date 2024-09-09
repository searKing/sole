#! /bin/python3
# -*- coding: utf-8 -*-

import base64
import ctypes
import json
import random
import requests
import string
import time
from concurrent.futures import ThreadPoolExecutor, ProcessPoolExecutor
from importlib import import_module


def test_http_performance_pool_loop(url, json, headers=None, cookies=None, total_requests_per_worker=1, min_workers=0,
                                    max_workers=1, use_multiprocessing=False,
                                    plot_report=False,
                                    plot_image_path=''):
    reports = []
    plot_workers = []
    plot_qps = []
    plot_costs_in_ms = []
    plot_success_rates_in_percent = []
    step = 1 if min_workers <= max_workers else -1
    print(f'---------------------------------')
    for workers in range(min_workers, max_workers + step, step):
        print(f'#{workers}--------------------------------')
        report = test_http_performance_pool(url, json, headers, cookies, total_requests_per_worker, workers,
                                            use_multiprocessing, prefix=f"  #{workers} ")
        plot_workers.append(workers)
        plot_qps.append(report["qps"])
        plot_costs_in_ms.append(report["cost"] * 1000)
        plot_success_rates_in_percent.append(report["success_rate"] * 100)
        reports.append(f'#{workers} {report["report"]}')

    # print reports

    print(f'#---------------------------------')
    print(
        f'Reports of {image_path}: \nTesting {total_requests_per_worker} per worker with [{min_workers},{max_workers}] {"processes" if use_multiprocessing else "threads"}')
    for report in reports:
        print(f'{report}')
    print(f'#---------------------------------')
    if not plot_image_path:
        plot_image_path = 'report.png'
    if plot_report:
        plot_report_figure(plot_image_path=plot_image_path,
                           plot_workers=plot_workers, plot_qps=plot_qps, plot_costs_in_ms=plot_costs_in_ms,
                           plot_success_rates_in_percent=plot_success_rates_in_percent)


def test_http_performance_pool(url, json, headers=None, cookies=None, total_requests_per_worker=1, workers=1,
                               use_multiprocessing=False, prefix=''):
    if use_multiprocessing:
        with ProcessPoolExecutor(max_workers=workers) as executor:
            return test_http_performance(executor, url, json, headers, cookies,
                                         total_requests=total_requests_per_worker * workers,
                                         prefix=prefix)
    else:
        with ThreadPoolExecutor(max_workers=workers) as executor:
            return test_http_performance(executor, url, json, headers, cookies,
                                         total_requests=total_requests_per_worker * workers,
                                         prefix=prefix)


def test_http_performance(executor, url, json, headers=None, cookies=None, total_requests=1, prefix=''):
    start_time = time.time()
    futures = [executor.submit(send_request, url, json, headers, cookies, prefix) for _ in range(total_requests)]

    failed_requests = 0
    first_resp = None
    total_elapse_time = 0
    max_time = 0
    min_time = -1
    for future in futures:
        elapsed_time, resp = future.result()
        if not resp:
            failed_requests += 1
        if resp:
            first_resp = resp
            if elapsed_time:
                total_elapse_time += elapsed_time
                if elapsed_time > max_time:
                    max_time = elapsed_time
                if elapsed_time < min_time or min_time < 0:
                    min_time = elapsed_time

    if total_requests == 1:
        print(f'{beautify_json(simplify_json(first_resp))}')
    end_time = time.time()
    elapsed_time = end_time - start_time
    if elapsed_time and (total_requests - failed_requests):
        qps = (total_requests - failed_requests) / elapsed_time
        cost = total_elapse_time / (total_requests - failed_requests)
        success_rate = (total_requests - failed_requests) / total_requests
        return {
            "qps": qps,
            "cost": cost,
            "success_rate": success_rate,
            "report": f"QPS:{qps:.2f}, cost: {cost * 1000:.2f}ms,∈[{min_time * 1000.0:.2f}, {max_time * 1000.0:.2f}]ms, {total_requests - failed_requests}/{total_requests} = {success_rate * 100.0:.2f}%"
        }
    else:
        return {
            "qps": 0,
            "cost": 0,
            "success_rate": 0,
            "report": f"No Succeed Requests: {failed_requests}/{total_requests}= {failed_requests / total_requests * 100.0:.2f}%"
        }


def send_request(url, json, headers=None, cookies=None, prefix=''):
    start_time = time.time()
    elapsed_time = 0
    try:
        response = requests.post(url, json=json, headers=headers, cookies=cookies)
        elapsed_time = time.time() - start_time
    except requests.exceptions.RequestException as e:
        print(f"{prefix}Request Exception: {e}, cost={elapsed_time * 1000.0:.2f}ms")
        return elapsed_time, None
    if response.status_code != 200:
        print(f"{prefix}Request Failed: {response.status_code} {response.reason}, cost={elapsed_time * 1000.0:.2f}ms")
        return elapsed_time, None
    print(f"\033[90m{prefix}Request OK, cost={elapsed_time * 1000.0:.2f}ms\033[0m")
    response.encoding = 'utf-8'
    return elapsed_time, response.json()


def file_to_base64(file_path) -> str:
    with open(file_path, "rb") as file:
        return base64.b64encode(file.read()).decode('utf-8')


def beautify_json(obj) -> str:
    return f'{json.dumps(obj, indent=4, ensure_ascii=False)}'


def simplify_json(obj):
    if isinstance(obj, dict):
        return {k: simplify_json(v) for k, v in obj.items()}
    elif isinstance(obj, list):
        return [simplify_json(v) for v in obj]
    elif isinstance(obj, str):
        if len(obj) > 1024:
            return f'{obj[:8]} {len(obj)} bytes'
        return obj
    else:
        return obj


def generate_random_string(length):
    result = ''.join(random.choices(string.ascii_letters + string.digits, k=length))
    return result


def python2_hash(str_info: str) -> int:
    length = len(str_info)

    if length == 0:
        return 0

    x = ord(str_info[0]) << 7
    for c in str_info:
        x = (1000003 * x) ^ ord(c)

    x ^= length
    x &= 0xffffffffffffffff
    if x == -1:
        x = -2

    # Convert to C long type
    v = ctypes.c_long(x).value
    return v


def plot_report_figure(plot_image_path='report.png', plot_workers=None, plot_qps=None, plot_costs_in_ms=None,
                       plot_success_rates_in_percent=None):
    pyplot = import_module("matplotlib.pyplot")
    ticker = import_module("matplotlib.ticker")
    if pyplot is None or ticker is None:
        print("report figure is not plot, please install matplotlib first")
        return
    # 绘制测试报告
    if plot_workers is None:
        plot_workers = []
    if plot_costs_in_ms is None:
        plot_costs_in_ms = []
    if plot_costs_in_ms is None:
        plot_costs_in_ms = []
    if plot_success_rates_in_percent is None:
        plot_success_rates_in_percent = []

    figure, ax_qps = pyplot.subplots()
    ax_qps.set_title('QPS, Latency and Concurrency over Time')
    lines = []
    # 绘制QPS
    line, *_ = ax_qps.plot(plot_workers, plot_qps, label="QPS", color='tab:blue')
    lines.append(line)
    ax_qps.set_xlabel("Workers")
    ax_qps.set_ylabel('QPS', color='tab:blue')
    ax_qps.tick_params(axis='y', labelcolor='tab:blue')
    ax_qps.xaxis.set_major_locator(ticker.MaxNLocator(integer=True))
    ax_qps.grid(True)

    # 绘制平均耗时
    ax_cost = ax_qps.twinx()
    line, *_ = ax_cost.plot(plot_workers, plot_costs_in_ms, label="Cost(ms)", color='tab:red')
    lines.append(line)
    ax_cost.set_ylabel('Cost(ms)', color='tab:red')
    ax_cost.tick_params(axis='y', labelcolor='tab:red')
    # ax_cost.grid(axis='y', linestyle='--')

    # 绘制成功率
    ax_success_rate = ax_qps.twinx()
    # 隐藏第三个Y轴的刻度线
    ax_success_rate.spines["left"].set_position(("axes", -1))
    line, *_ = ax_success_rate.plot(plot_workers, plot_success_rates_in_percent, label="Success Rate(%)",
                                    color='tab:gray')
    lines.append(line)
    ax_success_rate.set_ylabel('Success Rate(%)', color='tab:gray')
    ax_success_rate.yaxis.tick_left()
    ax_success_rate.yaxis.set_label_position("left")

    figure.legend(handles=lines, loc='upper left', bbox_to_anchor=(0, 0.9),
                  bbox_transform=ax_success_rate.transAxes)

    # report_figure_path = f'report.{datetime.now().strftime("%Y%m%d%H%M%S")}.png'
    report_figure_path = f'{plot_image_path}'
    figure.savefig(report_figure_path, dpi=300)
    print(f'report figure is saved {report_figure_path}')
    figure.show()
