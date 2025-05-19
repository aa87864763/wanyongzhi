package task

import (
	"encoding/json"
	"net/http"
	"strings"
)

func TasksAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	var filterOrID string
	if len(pathParts) > 2 {
		filterOrID = pathParts[2]
	}

	switch r.Method {
	case http.MethodGet:
		var tasksToReturn []Task
		switch filterOrID {
		case "", "all":
			tasksToReturn = GetAllTasks()
		case "done":
			tasksToReturn = GetFilteredTasks("done")
		case "undone":
			tasksToReturn = GetFilteredTasks("undone")
		default:
			tasksToReturn = GetAllTasks()
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasksToReturn)

	case http.MethodPost:

		if filterOrID != "" {
			http.Error(w, "POST 请求应发送至 /api/tasks/", http.StatusBadRequest)
			return
		}
		var newTask Task

		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "无效的请求体: "+err.Error(), http.StatusBadRequest)
			return
		}

		createdTask, err := AddTask(newTask)
		if err != nil {
			http.Error(w, "添加任务失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdTask)

	case http.MethodPut:
		taskID := filterOrID
		if taskID == "" || taskID == "all" || taskID == "done" || taskID == "undone" {
			http.Error(w, "PUT 请求需要任务ID，且不能是过滤器名称。", http.StatusBadRequest)
			return
		}
		var payload struct {
			Completed bool `json:"completed"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "PUT 请求的请求体无效: "+err.Error(), http.StatusBadRequest)
			return
		}
		updatedTask, err := UpdateTaskCompletion(taskID, payload.Completed)
		if err != nil {
			if err.Error() == "找不到任务" {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, "更新任务失败: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedTask)

	case http.MethodDelete:
		taskID := filterOrID
		if taskID == "" || taskID == "all" || taskID == "done" || taskID == "undone" {
			http.Error(w, "DELETE 请求需要任务ID，且不能是过滤器名称。", http.StatusBadRequest)
			return
		}
		err := DeleteTask(taskID)
		if err != nil {
			if err.Error() == "找不到任务" {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, "删除任务失败: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
	}
}
