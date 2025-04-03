package main

import (
	"fmt"
	"log"
	"time"

	"currency_analyzer/internal/config"
	"currency_analyzer/internal/fetcher"
	"currency_analyzer/internal/analyzer"
	"currency_analyzer/internal/presenter"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Установка периода анализа (последние 90 дней)
	endDate := time.Now().UTC()
	startDate := endDate.AddDate(0, 0, -cfg.DaysToAnalyze+1) // +1 чтобы включая текущий день

	fmt.Printf("Анализ курсов с %s по %s\n\n", 
		startDate.Format("02.01.2006"), 
		endDate.Format("02.01.2006"))

	// Инициализация сборщика данных
	dataFetcher := fetcher.NewFetcher(cfg.CBRLink)

	// Получение данных
	rates, err := dataFetcher.FetchRates(startDate, endDate)
	if err != nil {
		log.Fatalf("Ошибка получения данных: %v", err)
	}

	if len(rates) == 0 {
		log.Println("Нет данных для анализа")
		return
	}

	// Анализ данных
	stats := analyzer.Analyze(rates)

	// Вывод результатов
	presenter.Present(stats)
}