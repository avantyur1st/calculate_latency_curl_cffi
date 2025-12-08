package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

// calculatePercentile вычисляет перцентиль из набора данных
func calculatePercentile(data []float64, percentile float64) float64 {
	if len(data) == 0 {
		return 0
	}

	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)

	n := len(sorted)
	index := (percentile / 100.0) * float64(n-1)
	lower := int(index)
	upper := lower + 1

	if upper >= n {
		return sorted[lower]
	}

	fraction := index - float64(lower)
	return sorted[lower] + fraction*(sorted[upper]-sorted[lower])
}

// measureLatency измеряет задержку для указанного URL
func measureLatency(url string, numTests int, delay time.Duration) {
	// Создаём HTTP клиент с разумными таймаутами
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// Прогревочный запрос (игнорируем его время)
	fmt.Println("Performing warm-up request...")
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Warm-up request failed: %v", err)
	} else {
		resp.Body.Close()
	}

	latencies := make([]float64, 0, numTests)

	for i := 0; i < numTests; i++ {
		startTime := time.Now()
		resp, err := client.Get(url)
		endTime := time.Now()

		if err != nil {
			log.Printf("Test %d failed: %v", i+1, err)
			continue
		}
		resp.Body.Close()

		latency := endTime.Sub(startTime).Seconds()
		latencies = append(latencies, latency)

		fmt.Printf("Test %d: Latency = %.2f ms\n", i+1, latency*1000)

		// Задержка между запросами (кроме последнего)
		if i < numTests-1 {
			time.Sleep(delay)
		}
	}

	if len(latencies) > 0 {
		// Вычисляем статистику
		var sum float64
		minLatency := latencies[0]
		maxLatency := latencies[0]

		for _, lat := range latencies {
			sum += lat
			if lat < minLatency {
				minLatency = lat
			}
			if lat > maxLatency {
				maxLatency = lat
			}
		}

		avgLatency := sum / float64(len(latencies))
		percentile95 := calculatePercentile(latencies, 95)
		percentile99 := calculatePercentile(latencies, 99)

		// Выводим результаты
		fmt.Printf("\nAverage latency: %.2f ms\n", avgLatency*1000)
		fmt.Printf("Min latency: %.2f ms\n", minLatency*1000)
		fmt.Printf("Max latency: %.2f ms\n", maxLatency*1000)
		fmt.Printf("95th percentile latency: %.2f ms\n", percentile95*1000)
		fmt.Printf("99th percentile latency: %.2f ms\n", percentile99*1000)
	} else {
		fmt.Println("No successful tests completed")
	}
}

func main() {
	// Парсим аргументы командной строки
	url := flag.String("url", "", "The URL to test (with parameters if needed)")
	numTests := flag.Int("num_tests", 5, "Number of tests to run (after warm-up)")
	delaySeconds := flag.Float64("delay", 1.0, "Delay in seconds between requests")

	flag.Parse()

	// Проверяем обязательный параметр URL
	if *url == "" {
		fmt.Println("Error: URL is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		return
	}

	delay := time.Duration(*delaySeconds * float64(time.Second))
	measureLatency(*url, *numTests, delay)
}
