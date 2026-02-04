package daysteps

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
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, 0, errors.New("некорректное количество шагов")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil || duration <= 0 {
		return 0, 0, errors.New("некорректная продолжительность")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 || duration <= 0 {
		return ""
	}


	distanceKm := float64(steps) * stepLength / mInKm

	var coefficient float64
	if height <= 1.75 {
		coefficient = 0.605982905982906
	} else if height >= 1.85 {
		coefficient = 0.641025641025641
	} else {
		// Линейная интерполяция между 1.75 и 1.85
		coefficient = 0.605982905982906 + (0.641025641025641-0.605982905982906)*(height-1.75)/(1.85-1.75)
	}

	calories := weight * distanceKm * coefficient

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)
}