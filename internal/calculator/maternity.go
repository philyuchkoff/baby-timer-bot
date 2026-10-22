package calculator

import (
	"time"
)

// PregnancyCalculator рассчитывает сроки беременности и ПДР
type PregnancyCalculator struct{}

// NewPregnancyCalculator создаёт новый калькулятор
func NewPregnancyCalculator() *PregnancyCalculator {
	return &PregnancyCalculator{}
}

// NaegeleDueDate рассчитывает ПДР по формуле Негеле (LMP + 280 дней)
func (c *PregnancyCalculator) NaegeleDueDate(lmp time.Time) time.Time {
	return lmp.AddDate(0, 0, 280) // 40 недель = 280 дней
}

// CurrentGestationalWeek рассчитывает текущий акушерский срок в неделях и днях
func (c *PregnancyCalculator) CurrentGestationalWeek(lmp time.Time, currentDate time.Time) (weeks, days int) {
	diff := currentDate.Sub(lmp)
	totalDays := int(diff.Hours() / 24)
	weeks = totalDays / 7
	days = totalDays % 7
	return
}

// MaternityLeaveStart рассчитывает дату начала декретного отпуска
// birthType: "single" (обычная), "multiple" (многоплодная)
func (c *PregnancyCalculator) MaternityLeaveStart(dueDate time.Time, birthType string) time.Time {
	switch birthType {
	case "multiple":
		// При многоплодной: 84 дня до ПДР
		return dueDate.AddDate(0, 0, -84)
	default:
		// При одноплодной: 70 дней до ПДР
		return dueDate.AddDate(0, 0, -70)
	}
}

// LeaveDuration возвращает длительность декретного отпуска в днях
func (c *PregnancyCalculator) LeaveDuration(complications bool, multipleBirth bool) int {
	if multipleBirth {
		return 194 // 84 до + 110 после
	}
	if complications {
		return 156 // 70 до + 86 после
	}
	return 140 // 70 до + 70 после
}
