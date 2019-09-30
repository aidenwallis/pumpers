package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func startServer() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	workDir, _ := os.Getwd()
	fileServer(r, "/", http.Dir(filepath.Join(workDir, "web")))

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebsocket(w, r)
	})

	r.Get("/count", func(w http.ResponseWriter, r *http.Request) {
		count := &CountPayload{
			Total:  currentTotal,
			Day:    currentDay,
			Hour:   currentHour,
			Minute: currentMinute,
		}
		writeJSON(w, http.StatusOK, count)
	})

	http.ListenAndServe(":4338", r)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	bs, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to encode json")
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	write(w, code, bs)
}

func writeString(w http.ResponseWriter, code int, s string) {
	write(w, code, []byte(s))
}

func write(w http.ResponseWriter, code int, bs []byte) {
	w.WriteHeader(code)
	_, err := w.Write(bs)
	if err != nil {
		log.Debug("Failed to write response")
		log.Debug(err)
	}
}
