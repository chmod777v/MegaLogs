package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	metric "main.go/metrics"
)

type Base struct {
	Number int
	String string
}

func Get(w http.ResponseWriter, r *http.Request) {
	req := Base{
		Number: 123,
		String: "Hellow, world",
	}

	bytes, err := json.MarshalIndent(req, "", " ")
	if err != nil {
		slog.Error("Error json.Marshal")
		return
	}
	w.Write(bytes)
	slog.Info("Sucessfull post Request")
}
func Post(w http.ResponseWriter, r *http.Request) {
	var req Base
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Erro while receiving data", "ERROR", err.Error())
		return
	}
	slog.Info("RequestGet", "Data", req)
}
func Handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.ObserveRequest(time.Since(start), 200)
	}()
	switch r.Method {
	case http.MethodGet:
		Get(w, r)
	case http.MethodPost:
		Post(w, r)
	}

}
func main() {
	ipaddr := "0.0.0.0"

	http.HandleFunc("/", Handler)

	go func() {
		_ = metric.Listen(ipaddr + ":8081")

	}()

	slog.Info("Server listening", "Host", ipaddr+":8080")
	if err := http.ListenAndServe(ipaddr+":8080", nil); err != nil {
		slog.Error("Error starting server", "ERROR", err.Error())
	}

}
