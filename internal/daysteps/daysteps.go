package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
	// Коэффициент для расчёта длины шага на основе роста
	stepLengthCoefficient = 0.45
	// Коэффициент для расчёта калорий при ходьбе
	walkingCaloriesCoefficient = 0.5
	// Количество минут в часе
	minInH = 60
)

// parsePackage разбирает строку данных и возвращает количество шагов,
// продолжительность активности и ошибку (если есть).
// Формат входных данных: "шаги,продолжительность" (например, "1000,30m").
func parsePackage(data string) (int, time.Duration, error) {
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf(
			"неверный формат данных: ожидается 2 элемента, получено %d",
			len(parts),
		)
	}

	// Преобразуем первый элемент (количество шагов) в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf(
			"ошибка преобразования количества шагов: %w",
			err,
		)
	}

	// Проверяем, что количество шагов больше 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	// Преобразуем второй элемент в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf(
			"ошибка преобразования продолжительности: %w",
			err,
		)
	}

	return steps, duration, nil
}

// DayActionInfo обрабатывает данные о дневной активности,
// рассчитывает дистанцию, скорость и сожжённые калории,
// возвращает информативную строку с результатами.
func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о количестве шагов и продолжительности
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// Проверяем входные параметры для расчёта калорий
	if duration <= 0 || weight <= 0 || height <= 0 {
		log.Println("некорректные параметры для расчёта калорий")
		return ""
	}

	// Рассчитываем длину шага на основе роста
	stepLen := height * stepLengthCoefficient
	distForCalories := float64(steps) * stepLen / mInKm

	// Рассчитываем среднюю скорость
	durationHours := duration.Hours()
	avgSpeed := distForCalories / durationHours

	// Рассчитываем калории
	durationMinutes := duration.Minutes()
	calories := (weight * avgSpeed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	// Формируем и возвращаем строку
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)
}
