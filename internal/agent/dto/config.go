package dto

type Config struct {
	// Ключ для вычисления хеша.
	Key string
	// Включить поддержку асимметричного шифрования
	CryptoKey string `json:"crypto_key,omitempty"`
	// Адрес хоста в формате 8.8.8.8:8080
	Address string `json:"address,omitempty"`
	// Интервал отправки метрик.
	ReportInterval int `json:"report_interval,omitempty"`
	// Интервал опроса метрик.
	PollInterval int `json:"poll_interval,omitempty"`
	// Лимит
	RateLimit int
	// Разрешить отправку метрик одним пакетным запросом.
	Batch bool
}
