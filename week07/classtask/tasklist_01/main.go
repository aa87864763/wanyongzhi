package main

import (
	"fmt"
	"log"
	"mytodolist/task"
	"net/http"
)

func main() {
	http.HandleFunc("/api/tasks/", task.TasksAPIHandler)
	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/?filter=all", http.StatusFound)
	})

	http.HandleFunc("/done", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/?filter=done", http.StatusFound)
	})

	http.HandleFunc("/undone", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/?filter=undone", http.StatusFound)
	})

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	port := "8080"
	fmt.Printf("服务器正在启动，访问地址: http://localhost:%s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("服务器启动失败: %s\n", err)
	}
}
