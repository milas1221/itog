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

	stepLengthCoefficient = 0.45 
	mInKm                 = 1000
	minInH                = 60

	runningCaloriesCoefficient = 2.0   
	walkingCaloriesCoefficient = 1.0   
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

	stepLen := height * stepLengthCoefficient
	return float64(steps) * stepLen / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	return dist / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные данные")
	}

	return weight *
		meanSpeed(steps, height, duration) *
		duration.Minutes() / minInH,
		nil
}


func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные данные")
	}

	return weight *
		meanSpeed(steps, height, duration) *
		duration.Minutes() *
		walkingCaloriesCoefficient / minInH,
		nil
}


func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var spentCalories float64

	switch activity {
	case "Ходьба":
		spentCalories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("WalkingSpentCalories: %w", err)
		}
	case "Бег":
		spentCalories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("RunningSpentCalories: %w", err)
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		distance,
		speed,
		spentCalories,
	), nil
}
