package models

import "time"

// ValCurs соответствует корневому элементу XML от ЦБ
type ValCurs struct {
	Date   string   `xml:"Date,attr"`  
	Valute []Valute `xml:"Valute"`     
}

// Valute описывает данные о валюте
type Valute struct {
	CharCode string `xml:"CharCode"` 
	Name     string `xml:"Name"`     
	Nominal  string `xml:"Nominal"`  
	Value    string `xml:"Value"`    
}

// CurrencyRate содержит обработанные данные
type CurrencyRate struct {
	Date     time.Time 
	CharCode string    
	Name     string    
	Rate     float64   
}