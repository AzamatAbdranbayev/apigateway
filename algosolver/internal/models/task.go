package models

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

const TaskWater_level = "water_level"

type Task struct {
	Id          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskSolutionRequest struct {
	Id     string `json:"id" validate:"required"`
	UserId string `json:"user_id" validate:"required"`
}
type UserChangeBalanceRequest struct {
	Id      string  `json:"id" validate:"required"`
	Balance float64 `json:"balance" validate:"required"`
}
type TaskSolutionResponse struct {
	Description string `json:"description"`
	Answer      string `json:"answer" `
	Input       string `json:"input"`
}

type TaskPriceRequest struct {
	Id   string  `json:"id" validate:"required"`
	Cost float64 `json:"cost" validate:"required"`
}

func (t *Task) GetSolution() (TaskSolutionResponse, error) {
	var res TaskSolutionResponse
	switch t.Type {
	case TaskWater_level:
		height := []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}
		result := calculateWaterLevels(height)
		log.Println(result)
		bts, err := json.Marshal(height)
		if err != nil {
			return res, err
		}
		res.Description = t.Description
		res.Input = string(bts)
		res.Answer = strconv.Itoa(result)
	}
	return res, nil
}

func calculateWaterLevels(height []int) int {
	if len(height) < 3 {
		return 0 // Не хватает высот для образования озер
	}

	left := 0                // Левый указатель
	right := len(height) - 1 // Правый указатель
	leftMax := 0             // Максимальная высота слева
	rightMax := 0            // Максимальная высота справа
	result := 0              // Результат (количество уровней, которые можно заполнить водой)

	for left < right {
		if height[left] < height[right] {
			if height[left] > leftMax {
				leftMax = height[left]
			} else {
				result += leftMax - height[left]
			}
			left++
		} else {
			if height[right] > rightMax {
				rightMax = height[right]
			} else {
				result += rightMax - height[right]
			}
			right--
		}
	}

	return result
}

type UserTaskHistory struct {
	UserId    string    `json:"user_id"`
	TaskId    string    `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTaskHistoryRequest struct {
	UserId string  `json:"user_id" validate:"required"`
	Limit  float64 `json:"limit" validate:"required"`
	Page   float64 `json:"page"`
}
