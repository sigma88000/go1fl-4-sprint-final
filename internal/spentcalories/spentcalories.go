package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага. (в этом пакете не используется, но пусть будет)
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format")
	}

	stepsStr := parts[0]
	// Строгий формат: пробелы вокруг числа запрещены.
	if strings.TrimSpace(stepsStr) != stepsStr {
		return 0, "", 0, fmt.Errorf("invalid steps format")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return 0, "", 0, fmt.Errorf("invalid steps")
	}

	activity := parts[1]

	durationStr := parts[2]
	// Пробелы в длительности тоже запрещены ("1 h30m" — ошибка).
	if strings.TrimSpace(durationStr) != durationStr {
		return 0, "", 0, fmt.Errorf("invalid duration format")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil || duration <= 0 {
		return 0, "", 0, fmt.Errorf("invalid duration")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	return (float64(steps) * stepLen) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("steps must be positive")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be positive")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be positive")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	calories := (weight * speed * minutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("steps must be positive")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be positive")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be positive")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	calories := (weight * speed * minutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}