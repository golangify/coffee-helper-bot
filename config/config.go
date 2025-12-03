package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	IsRelease            bool          `json:"is_release" koanf:"is_release"`
	TelegramAPIToken     string        `json:"telegram_api_token" koanf:"telegram_api_token"`
	Database             string        `json:"database" koanf:"database"`
	Whitelist            []string      `json:"whitelist" koanf:"whitelist"`
	DelayBetweenActivity time.Duration `json:"delay_between_activity" koanf:"delay_between_activity"`
	ItemsPerPage         int           `json:"items_per_page"`
}

// LoadJSON загружает конфигурацию из JSON файла с помощью koanf.
// Если файл не существует, создается новый с настройками по умолчанию.
func LoadJSON(path string) (*Config, error) {
	k := koanf.New(".")

	// Пытаемся загрузить конфигурацию из файла
	err := k.Load(file.Provider(path), kjson.Parser())
	if err != nil {
		// Если файл не существует, создаем новый
		if os.IsNotExist(err) {
			return createDefaultConfig(path)
		}
		// Для других ошибок возвращаем их как есть
		return nil, err
	}

	// Успешно загрузили конфигурацию, парсим ее
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// createDefaultConfig создает конфигурационный файл с настройками по умолчанию
func createDefaultConfig(path string) (*Config, error) {
	defaultConfig := &Config{
		IsRelease:            false,
		TelegramAPIToken:     "your_telegram_bot_token_here",
		Database:             "database.db",
		Whitelist:            []string{},
		DelayBetweenActivity: time.Second / 3,
		ItemsPerPage:         10,
	}

	// Создаем директорию, если она не существует
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	// Сохраняем конфиг по умолчанию в файл
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return nil, err
	}

	return defaultConfig, nil
}
