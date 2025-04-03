package presenter

import (
	"currency_analyzer/internal/analyzer"
	"currency_analyzer/internal/models"
	"fmt"
)

// Present отображает результаты анализа
func Present(stats analyzer.Stats) {
	fmt.Println("Максимальный курс:")
	вывестиКурс(stats.Max)

	fmt.Println("\nМинимальный курс:")
	вывестиКурс(stats.Min)

	fmt.Printf("\nСредний курс за период: %.4f RUB\n", stats.Avg)
}

// Вывод данных о курсе валюты
func вывестиКурс(курс models.CurrencyRate) {
	fmt.Printf("  Значение:    %.4f RUB\n", курс.Rate)
	fmt.Printf("  Валюта:      %s (%s)\n", курс.CharCode, курс.Name)
	fmt.Printf("  Дата:        %s\n", курс.Date.Format("2006-01-02"))
}
