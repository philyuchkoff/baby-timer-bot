package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/telebot.v3"

	"github.com/philyuchkoff/baby-timer-bot/internal/calculator"
)

// Простое состояние для хранения типа декрета
var userLeaveType = make(map[int64]string)

// formatMoney форматирует число в российский формат (без пробелов)
func formatMoney(amount float64) string {
	// Округляем до целых рублей для больших сумм
	intAmount := int64(amount)
	return strconv.FormatInt(intAmount, 10)
}

// MainMenu возвращает главное меню бота
func MainMenu() *tb.ReplyMarkup {
	markup := &tb.ReplyMarkup{ResizeKeyboard: true}

	btnPDR := markup.Text("📅 Рассчитать ПДР и срок")
	btnBenefit := markup.Text("💰 Рассчитать декретные")
	btnInfo := markup.Text("ℹ️ О нормативах 2026")

	markup.Reply(
		markup.Row(btnPDR, btnBenefit),
		markup.Row(btnInfo),
	)

	return markup
}

// BenefitMenu возвращает меню выбора типа беременности
func BenefitMenu() *tb.ReplyMarkup {
	markup := &tb.ReplyMarkup{}

	// Создаём кнопки через Data
	btnSingle := markup.Data("👶 Обычная (140 дней)", "single")
	btnCompl := markup.Data("⚠️ Осложнённые (156 дней)", "compl")
	btnMultiple := markup.Data("👶👶 Многоплодная (194 дня)", "multiple")

	markup.Inline(
		markup.Row(btnSingle),
		markup.Row(btnCompl),
		markup.Row(btnMultiple),
	)

	return markup
}

// SetupHandlers настраивает все команды бота
func SetupHandlers(b *tb.Bot, calc *calculator.PregnancyCalculator, benefitCalc *calculator.BenefitCalculator) {

	// ОБРАБОТЧИК ВСЕХ CALLBACK
	b.Handle(tb.OnCallback, func(c tb.Context) error {
		// Очищаем данные от пробелов и переносов строк
		rawData := c.Callback().Data
		data := strings.TrimSpace(rawData)

		// Отвечаем на callback, чтобы убрать индикатор загрузки
		c.Respond(&tb.CallbackResponse{})

		switch data {
		case "single":
			userLeaveType[c.Sender().ID] = "single"
			c.Delete()
			return c.Send(`💰 Введите ваш суммарный заработок за 2024-2025 годы (в рублях):

Например: 1200000`, tb.ModeMarkdown)

		case "compl":
			userLeaveType[c.Sender().ID] = "compl"
			c.Delete()
			return c.Send(`💰 Введите ваш суммарный заработок за 2024-2025 годы (в рублях):

Например: 1200000`, tb.ModeMarkdown)

		case "multiple":
			userLeaveType[c.Sender().ID] = "multiple"
			c.Delete()
			return c.Send(`💰 Введите ваш суммарный заработок за 2024-2025 годы (в рублях):

Например: 1200000`, tb.ModeMarkdown)
		}

		return nil
	})

	// Команда /start
	b.Handle("/start", func(c tb.Context) error {
		welcome := `👋 Привет! Я помогу рассчитать:

📅 Акушерский срок и ПДР
👶 Дату выхода в декрет
💰 Сумму декретных выплат

Если есть вопросы или замечания - напиши @philyuchkoff

Для расчетов используй кнопки ниже 👇`

		return c.Send(welcome, MainMenu())
	})

	// Обработка текстовых сообщений
	b.Handle(tb.OnText, func(c tb.Context) error {
		text := c.Text()

		// Обработка кнопки "📅 Рассчитать ПДР и срок"
		if text == "📅 Рассчитать ПДР и срок" {
			return c.Send(`📅 Введите дату первого дня последних месячных в формате ДД.ММ.ГГГГ
        
Например: 15.03.2026`)
		}

		// Обработка кнопки "💰 Рассчитать декретные"
		if text == "💰 Рассчитать декретные" {
			return c.Send("Выберите тип беременности:", BenefitMenu())
		}

		// Обработка кнопки "ℹ️ О нормативах 2026"
		if text == "ℹ️ О нормативах 2026" {
			info := fmt.Sprintf(`📋 *Нормативы 2026 года для расчёта декретных выплат*

*Что такое декретные выплаты?*
Это пособие по беременности и родам (БиР), которое выплачивается единоразово за весь период отпуска.

*Длительность декретного отпуска:*
• 👶 Обычная беременность (один ребёнок) — *140 дней* (70 дней до родов + 70 дней после)
• ⚠️ Осложнённые роды (кесарево сечение и др.) — *156 дней* (70 дней до + 86 дней после)
• 👶👶 Многоплодная беременность (двойня, тройня) — *194 дня* (84 дня до + 110 дней после)

*Минимальные выплаты в 2026 году (МРОТ = 27093 ₽):*
• За 140 дней: *%s ₽*
• За 156 дней: *%s ₽*
• За 194 дня: *%s ₽*

*Максимальные выплаты в 2026 году (предельная база):*
• За 140 дней: *%s ₽*
• За 156 дней: *%s ₽*
• За 194 дня: *%s ₽*

*Как рассчитывается сумма?*
Пособие = (заработок за 2 года / (730 - исключаемые дни)) × дни декрета

*Важно:* 
• Для расчёта берутся 2024-2025 годы (730 дней)
• Если ваш заработок ниже МРОТ применяется минимальная сумма
• Если выше предельной базы применяется максимальная сумма

_Данные актуальны на 2026 год согласно Постановлению Правительства РФ_`,
				formatMoney(124702),
				formatMoney(138953),
				formatMoney(172801),
				formatMoney(955836),
				formatMoney(1065074),
				formatMoney(1324515))

			return c.Send(info, tb.ModeMarkdown)
		}

		// Пропускаем команды
		if strings.HasPrefix(text, "/") {
			return nil
		}

		// Пробуем распарсить как дату (для ПДР)
		lmp, err := time.Parse("02.01.2006", text)
		if err == nil {
			dueDate := calc.NaegeleDueDate(lmp)
			weeks, days := calc.CurrentGestationalWeek(lmp, time.Now())
			leaveStart := calc.MaternityLeaveStart(dueDate, "single")

			result := fmt.Sprintf(`📊 *Результаты расчёта*

🗓️ Первый день последних месячных: %s
🤰 Акушерский срок: *%d недель %d дней*
👶 ПДР (предполагаемая дата родов): %s
👔 Дата выхода в декрет (30 недель): %s

_*Расчёт по формуле Негеле (LMP + 280 дней)_`,
				"`"+lmp.Format("02.01.2006")+"`",
				weeks, days,
				"`"+dueDate.Format("02.01.2006")+"`",
				"`"+leaveStart.Format("02.01.2006")+"`")

			return c.Send(result, tb.ModeMarkdown)
		}

		// Пробуем распарсить как зарплату
		cleanText := strings.ReplaceAll(text, ",", ".")
		cleanText = strings.ReplaceAll(cleanText, " ", "")
		salary, err := strconv.ParseFloat(cleanText, 64)
		if err == nil {
			leaveType, exists := userLeaveType[c.Sender().ID]
			if !exists {
				return c.Send(`❌ Сначала выберите тип беременности через кнопку "💰 Рассчитать декретные"`)
			}

			var leaveDays int
			var pregnancyType string
			switch leaveType {
			case "compl":
				leaveDays = 156
				pregnancyType = "Осложнённая беременность (156 дней)"
			case "multiple":
				leaveDays = 194
				pregnancyType = "Многоплодная беременность (194 дня)"
			default:
				leaveDays = 140
				pregnancyType = "Обычная беременность (140 дней)"
			}

			input := calculator.BenefitInput{
				LeaveDays:        leaveDays,
				TotalEarnings:    salary,
				ExcludedDays:     0,
				YearsCalculation: "2024-2025",
			}

			result, err := benefitCalc.Calculate(input)
			if err != nil {
				return c.Send("❌ Ошибка: " + err.Error())
			}

			response := fmt.Sprintf(`💰 *Результат расчёта пособия по беременности и родам*

🤰 Тип беременности: *%s*
📅 Длительность декрета: *%d дней*

📊 Средний дневной заработок: *%s ₽*
💵 Итоговое пособие: *%s ₽*

📉 Минимальная граница: *%s ₽*
📈 Максимальная граница: *%s ₽*`,
				pregnancyType,
				leaveDays,
				formatMoney(result.DailyWage),
				formatMoney(result.FinalBenefit),
				formatMoney(result.MinLimit),
				formatMoney(result.MaxLimit))

			if result.AppliedMinLimit {
				response += "\n\n⚠️ *Применён минимальный лимит* (ваш заработок ниже МРОТ)"
			}
			if result.AppliedMaxLimit {
				response += "\n\n⚠️ *Применён максимальный лимит* (ваш заработок выше предельной базы)"
			}

			delete(userLeaveType, c.Sender().ID)

			return c.Send(response, tb.ModeMarkdown)
		}

		return c.Send(`❌ Не распознано. Введите дату в формате ДД.ММ.ГГГГ или сумму зарплаты числом`)
	})
}
