# Currency Analyzer

Анализатор курсов валют ЦБ РФ

## Описание

Проект получает данные о курсах валют с официального сайта Центробанка России и предоставляет:
- Максимальный курс за период (валюта, значение, дата)
- Минимальный курс за период (валюта, значение, дата)
- Средний курс по всем валютам за период

## Установка

Клонируйте репозиторий:

git clone https://github.com/yourusername/currency_analyzer.git

cd currency_analyzer
go run ./cmd/analyzer 

либо 

cd cmd
cd analyzer
go run main.go