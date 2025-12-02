package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчётов.
const (
	// Средняя длина шага в метрах.
	lenStep = 0.65
	// Количество метров в километре.
	mInKm = 1000
	// Количество минут в часе.
	minInH = 60
	// Коэффициент для расчёта длины шага на основе роста.
	stepLengthCoefficient = 0.45
	// Коэффициент для расчёта калорий при ходьбе.
	walkingCaloriesCoefficient = 0.5
)

// parseTraining разбирает строку данных о тренировке и возвращает:
// - количество шагов (int),
// - тип активности (string),
// - продолжительность (time.Duration),
// - ошибку (если возникла).
//
// Формат входных данных: "шаги,тип_активности,продолжительность"
// (например, "1000,Ходьба,30m").
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf(
			"неверный формат данных: ожидается 3 элемента, получено %d",
			len(parts),
		)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf(
			"ошибка преобразования количества шагов: %w",
			err,
		)
	}

	activityType := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf(
			"ошибка преобразования продолжительности: %w",
			err,
		)
	}

	return steps, activityType, duration, nil
}

// distance рассчитывает пройденное расстояние в километрах
// на основе количества шагов и роста человека.
func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLen
	return distanceMeters / mInKm
}

// meanSpeed вычисляет среднюю скорость движения в км/ч.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	durationHours := duration.Hours()
	return dist / durationHours
}

// RunningSpentCalories рассчитывает сожжённые калории при беге.
func RunningSpentCalories(
	steps int,
	weight, height float64,
	duration time.Duration,
) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * avgSpeed * durationMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories рассчитывает сожжённые калории при ходьбе
// с учётом корректирующего коэффициента.
func WalkingSpentCalories(
	steps int,
	weight, height float64,
	duration time.Duration,
) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * avgSpeed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo обрабатывает данные о тренировке, рассчитывает
// и возвращает информативную строку с результатами.
func TrainingInfo(
	data string,
	weight, height float64,
) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64

	switch activityType {
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

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activityType,
		durationHours,
		dist,
		speed,
		calories,
	)

	return result, nil
}
