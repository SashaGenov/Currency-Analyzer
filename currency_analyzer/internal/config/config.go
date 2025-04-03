package config

import "time"

// Config содержит настройки приложения
type Config struct {
	CBRLink      string        // URL API ЦБ
	DaysToAnalyze int          // Количество дней для анализа
	Timeout      time.Duration // Таймаут HTTP-запросов
}

// Load возвращает конфигурацию по умолчанию
func Load() (*Config, error) {
	return &Config{
		CBRLink:      "http://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=11/11/2020%s",
		DaysToAnalyze: 90,
		Timeout:      10 * time.Second,
	}, nil
}