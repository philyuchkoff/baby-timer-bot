package config

import (
	"os"
)

// Config хранит конфигурацию бота
type Config struct {
	BotToken     string
	MinDailyWage float64
	MaxDailyWage float64
	MROT         float64
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		BotToken:     os.Getenv("TELEGRAM_BOT_TOKEN"),
		MinDailyWage: 890.73,  // МРОТ 27 093 / 30,4
		MaxDailyWage: 6827.40, // Предельная база 2026
		MROT:         27093,   // МРОТ 2026
	}
}

// GetMinBenefit возвращает минимальное пособие для заданного количества дней
func (c *Config) GetMinBenefit(days int) float64 {
	return c.MinDailyWage * float64(days)
}

// GetMaxBenefit возвращает максимальное пособие для заданного количества дней
func (c *Config) GetMaxBenefit(days int) float64 {
	return c.MaxDailyWage * float64(days)
}
