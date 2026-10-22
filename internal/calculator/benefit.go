package calculator

import (
	"errors"
)

// BenefitCalculator рассчитывает пособие по БиР
type BenefitCalculator struct {
	minDailyWage float64 // Минимальный средний дневной заработок на 2026 год
	maxDailyWage float64 // Максимальный средний дневной заработок на 2026 год
}

// NewBenefitCalculator создаёт новый калькулятор пособий с лимитами 2026 года
func NewBenefitCalculator() *BenefitCalculator {
	return &BenefitCalculator{
		minDailyWage: 890.73,  // МРОТ 27 093 / 30,4 ≈ 890,73 ₽
		maxDailyWage: 6827.40, // Предельная база 2026 года
	}
}

// BenefitInput содержит данные для расчёта пособия
type BenefitInput struct {
	LeaveDays        int     // Количество дней декретного отпуска
	TotalEarnings    float64 // Суммарный заработок за 2 года
	ExcludedDays     int     // Исключаемые дни (больничные, прошлые декреты)
	YearsCalculation string  // "2024-2025" (расчётные годы)
}

// BenefitResult содержит результат расчёта пособия
type BenefitResult struct {
	DailyWage       float64 // Средний дневной заработок
	FinalBenefit    float64 // Итоговое пособие
	MinLimit        float64 // Минимальная граница
	MaxLimit        float64 // Максимальная граница
	AppliedMinLimit bool    // Применён ли минимум
	AppliedMaxLimit bool    // Применён ли максимум
}

// Calculate рассчитывает пособие по беременности и родам
func (c *BenefitCalculator) Calculate(input BenefitInput) (*BenefitResult, error) {
	if input.LeaveDays <= 0 {
		return nil, errors.New("количество дней отпуска должно быть положительным")
	}

	// Расчётные дни: 730 или 731 (високосный год)
	calculationDays := 730 // По умолчанию 2024-2025: 365+365
	if input.YearsCalculation == "2023-2024" {
		calculationDays = 731 // 2024 високосный
	}

	effectiveDays := calculationDays - input.ExcludedDays
	if effectiveDays <= 0 {
		return nil, errors.New("количество исключаемых дней не может превышать общее количество дней")
	}

	// Средний дневной заработок
	dailyWage := input.TotalEarnings / float64(effectiveDays)

	result := &BenefitResult{
		DailyWage: dailyWage,
		MinLimit:  c.minDailyWage * float64(input.LeaveDays),
		MaxLimit:  c.maxDailyWage * float64(input.LeaveDays),
	}

	// Проверка на минимальный и максимальный пределы
	if dailyWage < c.minDailyWage {
		result.FinalBenefit = result.MinLimit
		result.AppliedMinLimit = true
		result.AppliedMaxLimit = false
	} else if dailyWage > c.maxDailyWage {
		result.FinalBenefit = result.MaxLimit
		result.AppliedMinLimit = false
		result.AppliedMaxLimit = true
	} else {
		result.FinalBenefit = dailyWage * float64(input.LeaveDays)
		result.AppliedMinLimit = false
		result.AppliedMaxLimit = false
	}

	return result, nil
}

// GetMinMaxForDisplay возвращает минимальные и максимальные суммы для отображения
func (c *BenefitCalculator) GetMinMaxForDisplay() map[string]float64 {
	return map[string]float64{
		"min_140": c.minDailyWage * 140, // 124 702 ₽
		"max_140": c.maxDailyWage * 140, // 955 836 ₽
		"min_156": c.minDailyWage * 156, // 138 953 ₽
		"max_156": c.maxDailyWage * 156, // 1 065 074 ₽
		"min_194": c.minDailyWage * 194, // 172 801 ₽
		"max_194": c.maxDailyWage * 194, // 1 324 515 ₽
	}
}
