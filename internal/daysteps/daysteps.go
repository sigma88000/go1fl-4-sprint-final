package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage парсит строку формата: "<steps>,<duration>"
// Примеры: "678,0h50m", "+12345,1h30m", "1000,30.5m"
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid data format")
	}

	stepsStr := parts[0]
	// Важно: пробелы запрещены тестами (например " 12345" и "12345 ")
	if strings.TrimSpace(stepsStr) != stepsStr {
		return 0, 0, fmt.Errorf("invalid steps format")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return 0, 0, fmt.Errorf("invalid steps")
	}

	durationStr := parts[1]
	// пробелы тоже считаем ошибкой
	if strings.TrimSpace(durationStr) != durationStr {
		return 0, 0, fmt.Errorf("invalid duration format")
	}

	dur, err := time.ParseDuration(durationStr)
	if err != nil || dur <= 0 {
		return 0, 0, fmt.Errorf("invalid duration")
	}

	return steps, dur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Дистанция по заданной длине шага 0.65м (именно для daysteps)
	distanceKm := (float64(steps) * stepLength) / mInKm

	// Калории считаем через spentcalories, как в тренажёре спринта:
	// активность тут — "Ходьба"
	cal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		cal,
	)
}