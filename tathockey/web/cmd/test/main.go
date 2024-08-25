package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"tat_hockey_pack/internal/configs/logger"
	"tat_hockey_pack/internal/configs/postgre"
	httpH "tat_hockey_pack/internal/handlers/http_handlers"
	"tat_hockey_pack/internal/middleware"
	"tat_hockey_pack/internal/repository"
	"tat_hockey_pack/internal/service/session"
	"tat_hockey_pack/internal/service/user"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	files, err := ioutil.ReadDir("/tmp")
	if err != nil {
		log.Fatal(err)
		return
	}

	for i, f := range files {
		fmt.Fprintf(file, "file %d Name: %s \n", i, f.Name())

	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", "err", err.Error())
		return
	}
	pg, _ := postgre.LoadPgxPool()
	log := logger.InitLogger()
	postRepo := repository.NewPostRepository(pg)
	vidRepo := repository.NewVideoRepository(pg, log)
	serv := user.NewService(log, repository.NewUserRepository(pg, log))
	sesRepo := repository.NewSessionRepository(pg, log)
	sess := session.NewService(log, sesRepo)

	smanager := httpH.NewSessionManager(sess, log)
	manager := httpH.NewUserManager(smanager, serv, log)
	pmanager := httpH.NewPostManager(log, postRepo, vidRepo, sesRepo)

	http.HandleFunc("/gif", pmanager.GifHandler)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	hand := middleware.LoggerMiddleware(log, mux)

	hand = middleware.Auth(smanager, log, hand)
	hand = middleware.RequestIDMiddleware(hand)
	hand = middleware.Panic(hand)

	routes := map[string]string{
		"/feeds":    "/app/tmp/feeds.html",
		"/signup":   "/app/tmp/signup.html",
		"/new_post": "/app/tmp/upload.html",
	}

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Login user", "main")
		if method := r.Method; method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		// Проверка активной сессии
		s, err := smanager.Check(r)
		if err == nil && s != nil {
			log.Error("mux.Login", "err", err.Error(), "msg", "no session")
			// Если сессия существует, перенаправляем на /feeds
			http.Redirect(w, r, "/feeds", http.StatusFound)
			return
		}

		http.ServeFile(w, r, "/app/tmp/login.html")
	})
	mux.HandleFunc("/post", pmanager.PostHandler)

	for path, file := range routes {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, file)
		})
	}

	mux.HandleFunc("/api/v1/login", manager.Login)
	mux.HandleFunc("/api/v1/logout", manager.Logout)

	mux.HandleFunc("/", smanager.Index)

	log.Info("Server start",
		"start on", "localhost:8000")
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: hand,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Error("ListenAndServe", "err", err.Error())
	}
}
