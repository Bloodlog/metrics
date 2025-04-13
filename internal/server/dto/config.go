package dto

type Config struct {
	// Ключ для вычисления хеша.
	Key string
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
	Debug bool
}
