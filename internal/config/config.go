package config

type AgentConfig struct {
	// Ключ для вычисления хеша.
	Key string `json:"-"`
	// Включить поддержку асимметричного шифрования
	CryptoKey string `json:"crypto_key,omitempty"`
	// Адрес хоста в формате 8.8.8.8:8080
	Address string `json:"address,omitempty"`
	// Интервал отправки метрик.
	ReportInterval int `json:"report_interval,omitempty"`
	// Интервал опроса метрик.
	PollInterval int `json:"poll_interval,omitempty"`
	// Лимит
	RateLimit int `json:"-"`
	// Разрешить отправку метрик одним пакетным запросом.
	Batch bool `json:"-"`
}

type ServerConfig struct {
	// Ключ для вычисления хеша.
	Key string `json:"-"`
	// Адрес сервера.
	Address string `json:"address,omitempty"`
	// Путь к файлу хранилищу.
	FileStoragePath string `json:"store_file,omitempty"`
	// Настройки БД в формате dsn.
	DatabaseDsn string `json:"database_dsn,omitempty"`
	// Включить поддержку асимметричного шифрования
	CryptoKey string `json:"crypto_key,omitempty"`
	// Интервал сохранения хранилища.
	StoreInterval int `json:"store_interval,omitempty"`
	// Разрешить загрузку из файла хранилища.
	Restore bool `json:"restore,omitempty"`
	// Разрешить отладку.
	Debug bool `json:"-"`
}
