package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/telebot.v3"

	"github.com/philyuchkoff/baby-timer-bot/internal/bot"
	"github.com/philyuchkoff/baby-timer-bot/internal/calculator"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используем системные переменные")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан")
	}

	// Настройка прокси (если нужен)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Пример для прокси (раскомментируйте и укажите ваш прокси)
	// proxyURL, _ := url.Parse("socks5://127.0.0.1:1080") // или http://proxy:8080
	// client.Transport = &http.Transport{
	//     Proxy: http.ProxyURL(proxyURL),
	// }

	// Настройка бота с кастомным клиентом
	prefs := tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: client,
	}

	b, err := tb.NewBot(prefs)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация калькуляторов
	calc := calculator.NewPregnancyCalculator()
	benefitCalc := calculator.NewBenefitCalculator()

	// Настройка обработчиков
	bot.SetupHandlers(b, calc, benefitCalc)

	log.Println("Бот baby-timer-bot запущен!")
	b.Start()
}
