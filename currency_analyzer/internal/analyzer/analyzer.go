package analyzer

import "currency_analyzer/internal/models"

// Stats содержит результаты анализа
type Stats struct {
	Max models.CurrencyRate // Максимальный курс
	Min models.CurrencyRate // Минимальный курс
	Avg float64             // Средний курс
}

// Analyze вычисляет статистику по валютам
func Analyze(rates []models.CurrencyRate) Stats {
	if len(rates) == 0 {
		return Stats{} 
	}

	// Инициализация первым элементом
	max := rates[0]
	min := rates[0]
	sum := 0.0

	// Итерация по всем курсам
	for _, r := range rates {
		sum += r.Rate
		// Обновление максимума
		if r.Rate > max.Rate {
			max = r
		}
		// Обновление минимума
		if r.Rate < min.Rate {
			min = r
		}
	}

	return Stats{
		Max: max,
		Min: min,
		Avg: sum / float64(len(rates)), // Расчет среднего
	}
}