package fetcher

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"currency_analyzer/internal/models"
	"golang.org/x/text/encoding/charmap" 
)

// Fetcher отвечает за получение данных с API ЦБ
type Fetcher struct {
	baseURL string    
	client  *http.Client 
}

// NewFetcher создает новый экземпляр Fetcher
func NewFetcher(baseURL string) *Fetcher {
	return &Fetcher{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second, // Таймаут запроса
		},
	}
}

// FetchRates получает курсы за период (параллельно)
func (f *Fetcher) FetchRates(startDate, endDate time.Time) ([]models.CurrencyRate, error) {
	var (
		allRates []models.CurrencyRate 
		wg       sync.WaitGroup         
		mu       sync.Mutex             
	)

	// Буферизированный канал для результатов
	results := make(chan []models.CurrencyRate, 10)

	// Параллельные запросы для каждой даты
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		wg.Add(1)
		go func(date time.Time) {
			defer wg.Done()
			rates, err := f.fetchDayRates(date)
			if err != nil {
				fmt.Printf("Ошибка для %s: %v\n", date.Format("2006-01-02"), err)
				return
			}
			results <- rates
		}(d)
	}

	// Закрытие канала после завершения всех горутин
	go func() {
		wg.Wait()
		close(results)
	}()

	// Сбор результатов
	for rates := range results {
		mu.Lock()
		allRates = append(allRates, rates...)
		mu.Unlock()
	}

	return allRates, nil
}

// fetchDayRates получает курсы на конкретную дату
func (f *Fetcher) fetchDayRates(date time.Time) ([]models.CurrencyRate, error) {
	url := fmt.Sprintf(f.baseURL, date.Format("02/01/2006"))

	// Создание запроса
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP-ошибка: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("код статуса %d", resp.StatusCode)
	}

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Декодирование XML с учетом windows-1251
	var valCurs models.ValCurs
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("неподдерживаемая кодировка: %s", charset)
	}

	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("ошибка декодирования XML: %v", err)
	}

	// Парсинг даты из ответа
	dateParsed, err := time.Parse("02.01.2006", valCurs.Date)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга даты: %v", err)
	}

	// Обработка валют
	var rates []models.CurrencyRate
	for _, v := range valCurs.Valute {
		// Конвертация номинала
		nominal, err := strconv.Atoi(v.Nominal)
		if err != nil {
			continue // Пропуск невалидных данных
		}

		// Замена запятой для корректного парсинга
		valueStr := strings.Replace(v.Value, ",", ".", 1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		// Расчет курса за 1 единицу валюты
		rates = append(rates, models.CurrencyRate{
			Date:     dateParsed,
			CharCode: v.CharCode,
			Name:     v.Name,
			Rate:     value / float64(nominal),
		})
	}

	return rates, nil
}