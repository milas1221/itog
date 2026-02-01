package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	stepLength = 0.65
	mInKm      = 1000
	minInH     = 60
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("некорректное количество шагов")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil || duration <= 0 {
		return 0, "", 0, errors.New("некорректная продолжительность")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 {
		return 0
	}
	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || duration <= 0 {
		return 0
	}
	hours := duration.Hours()
	return distance(steps, height) / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные данные")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	return (weight * speed * minutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные данные")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	return (weight * speed * minutes) / minInH, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var (
		name     string
		calories float64
	)

	switch activity {
	case "Бег":
		name = "Бег"
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		name = "Ходьба"
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		name,
		duration.Hours(),
		dist,
		speed,
		calories,
	), nil
}
