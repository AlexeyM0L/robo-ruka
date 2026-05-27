package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"

	"robo-ruka/internal/config"
	"robo-ruka/internal/handler"
	"robo-ruka/internal/repository"
	"robo-ruka/internal/service"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles(cfg.TemplatePath)
	if err != nil {
		log.Fatalf("parse template: %v", err)
	}

	repo := repository.NewRepository(cfg.StatePath)
	svc := service.NewService(repo)
	h := handler.New(svc, tmpl)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Index)
	
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	log.Printf("listening on %s, state=%s", addr, cfg.StatePath)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
