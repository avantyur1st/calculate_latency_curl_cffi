import time
import argparse
from curl_cffi import requests  # Импортируем как requests для совместимости

def calculate_percentile(data, percentile):
    if not data:
        return None
    sorted_data = sorted(data)
    n = len(sorted_data)
    index = (percentile / 100.0) * (n - 1)
    lower = int(index)
    upper = lower + 1
    if upper >= n:
        return sorted_data[lower]
    fraction = index - lower
    return sorted_data[lower] + fraction * (sorted_data[upper] - sorted_data[lower])

def measure_latency(url, num_tests, delay):
    # Создаем сессию один раз
    sess = requests.Session(impersonate="chrome")  # Можно выбрать браузер, например chrome
    
    # Первый запрос (игнорируем его время)
    print("Performing warm-up request...")
    sess.get(url)
    
    latencies = []
    
    for i in range(num_tests):
        start_time = time.time()
        response = sess.get(url)
        end_time = time.time()
        
        # Более точное время из response.elapsed (если нужно общее время запроса)
        # Но для полной задержки включая сеть используем time差
        latency = end_time - start_time
        latencies.append(latency)
        
        print(f"Test {i+1}: Latency = {latency * 1000:.2f} ms")
        
        # Задержка между запросами (кроме последнего)
        if i < num_tests - 1:
            time.sleep(delay)
    
    if latencies:
        avg_latency = sum(latencies) / len(latencies)
        min_latency = min(latencies)
        max_latency = max(latencies)
        percentile_95 = calculate_percentile(latencies, 95)  # 95-й перцентиль
        percentile_99 = calculate_percentile(latencies, 99)  # 99-й перцентиль
        
        print(f"\nAverage latency: {avg_latency * 1000:.2f} ms")
        print(f"Min latency: {min_latency * 1000:.2f} ms")
        print(f"Max latency: {max_latency * 1000:.2f} ms")
        print(f"95th percentile latency: {percentile_95 * 1000:.2f} ms")
        print(f"99th percentile latency: {percentile_99 * 1000:.2f} ms")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Measure latency to a URL using curl_cffi with persistent session.")
    parser.add_argument("url", type=str, help="The URL to test (with parameters if needed)")
    parser.add_argument("--num_tests", type=int, default=5, help="Number of tests to run (after warm-up)")
    parser.add_argument("--delay", type=float, default=1.0, help="Delay in seconds between requests")
    
    args = parser.parse_args()
    
    measure_latency(args.url, args.num_tests, args.delay)