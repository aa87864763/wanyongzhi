package task

import (
	"errors"
	"sort"
	"strconv"
	"sync"
)

type Task struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	AddedTime string `json:"addedTime"`
}

var (
	tasks  = make(map[string]Task)
	mu     sync.Mutex
	nextID = 1
)

func generateID() string {
	idStr := strconv.Itoa(nextID)
	nextID++
	return idStr
}

func GetAllTasks() []Task {
	mu.Lock()
	defer mu.Unlock()

	all := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		all = append(all, t)
	}

	sort.Slice(all, func(i, j int) bool {
		idI, _ := strconv.Atoi(all[i].ID)
		idJ, _ := strconv.Atoi(all[j].ID)
		return idI < idJ
	})
	return all
}

func GetFilteredTasks(filter string) []Task {
	mu.Lock()
	defer mu.Unlock()

	var filtered []Task
	for _, t := range tasks {
		include := false
		switch filter {
		case "done":
			if t.Completed {
				include = true
			}
		case "undone":
			if !t.Completed {
				include = true
			}
		}
		if include {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func AddTask(newTask Task) (Task, error) {
	mu.Lock()
	defer mu.Unlock()

	if newTask.Name == "" {
		return Task{}, errors.New("任务名称不能为空")
	}
	newTask.ID = generateID()
	tasks[newTask.ID] = newTask
	return newTask, nil
}

func UpdateTaskCompletion(taskID string, completed bool) (Task, error) {
	mu.Lock()
	defer mu.Unlock()

	task, ok := tasks[taskID]
	if !ok {
		return Task{}, errors.New("找不到任务")
	}
	task.Completed = completed
	tasks[taskID] = task
	return task, nil
}

func DeleteTask(taskID string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := tasks[taskID]; !ok {
		return errors.New("找不到任务")
	}
	delete(tasks, taskID)
	return nil
}
