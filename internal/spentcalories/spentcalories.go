package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")
	
	// Проверяем, что длина слайса равна 3
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 3 элемента, получено %d", len(parts))
	}
	
	// Преобразуем первый элемент (количество шагов) в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}
	
	// Проверяем, что количество шагов больше 0
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	
	// Получаем вид активности
	activityType := parts[1]
	
	// Преобразуем третий элемент в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}
	
	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше 0")
	}
	
	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// Рассчитываем длину шага
	stepLen := height * stepLengthCoefficient
	
	// Умножаем количество шагов на длину шага
	distanceMeters := float64(steps) * stepLen
	
	// Переводим в километры
	distanceKm := distanceMeters / mInKm
	
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0
	}
	
	// Вычисляем дистанцию
	dist := distance(steps, height)
	
	// Переводим продолжительность в часы
	durationHours := duration.Hours()
	
	// Вычисляем и возвращаем среднюю скорость
	return dist / durationHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем входные параметры на корректность
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
	
	// Рассчитываем среднюю скорость
	avgSpeed := meanSpeed(steps, height, duration)
	
	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()
	
	// Рассчитываем количество калорий
	calories := (weight * avgSpeed * durationMinutes) / minInH
	
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем входные параметры на корректность
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
	
	// Рассчитываем среднюю скорость
	avgSpeed := meanSpeed(steps, height, duration)
	
	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()
	
	// Рассчитываем количество калорий
	calories := (weight * avgSpeed * durationMinutes) / minInH
	
	// Умножаем на корректирующий коэффициент для ходьбы
	calories *= walkingCaloriesCoefficient
	
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получаем значения из строки данных
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	
	var calories float64
	
	// Проверяем вид тренировки и рассчитываем калории
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
	
	// Рассчитываем дистанцию и среднюю скорость
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()
	
	// Формируем и возвращаем строку
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, durationHours, dist, speed, calories)
	
	return result, nil
}