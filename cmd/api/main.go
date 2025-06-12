package main

import (
	"log/slog"
	"net/http"
	"task-api/internal/executor"
	"task-api/internal/factory"
	"task-api/internal/gateway"
	"task-api/internal/operator"
	"task-api/internal/repository"
	"task-api/pkg/webservice"
)

func main() {
	repo := repository.New()
	oper := operator.New(repo, executor.New())
	gat := gateway.New(repo, oper, factory.New())

	s := webservice.New()
	webservice.Register(s, "Tasks.CreateTask", gat.CreateTask)
	webservice.Register(s, "Tasks.ListTasks", gat.ListTasks)
	webservice.Register(s, "Tasks.CancelTask", gat.CancelTask)
	webservice.Register(s, "Tasks.DeleteTask", gat.DeleteTask)
	webservice.Register(s, "Tasks.GetTaskResult", gat.GetTaskResult)
	webservice.Register(s, "Tasks.GetTaskDetails", gat.GetTaskDetails)

	s.WithErrorMapper(mapError)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api", s.Handle)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error(err.Error())
	}
}
