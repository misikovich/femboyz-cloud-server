package main

import (
	"femboyz/db"
	"femboyz/env"
	"femboyz/handlers"
	"femboyz/ratelimiter"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

var devMode bool

func init() {
	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "Jan 02 15:04",
	}))
	slog.SetDefault(logger)

	env.LoadEnv()
	devMode = env.DevMode.Get() == "true"
	db.InitDB()
	handlers.Init()

	// test fetching files
	// temporary putting an entry in the database manually
	db.InsertFile(&db.File{
		PubID: "12345ABCDF",
		Meta: db.FileMeta{
			OriginalName:  "YTMusicUltimate_v1.4.1_YTM6.03.1.ipa",
			Size:          1024,
			FileType:      "application/octet-stream",
			Hash:          "test",
			LocalFileName: "SAFHWLEFKSDNMWBEMM",
		},
		Issuer: "test",
	})

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/admin", handlers.Admin)
	mux.HandleFunc("/{id}", handlers.FilePage)
	mux.HandleFunc("/p/{id}", handlers.PostPage)
	mux.HandleFunc("/api/v1/send", handlers.Send)
	mux.HandleFunc("/api/v1/pull/f", handlers.PullFile)
	mux.HandleFunc("/api/v1/pull/p", handlers.PullPost)

	rl, _ := strconv.ParseFloat(env.RateLimit.Get(), 64)
	rb, _ := strconv.Atoi(env.RateBurst.Get())
	limiter := ratelimiter.NewRateLimiter(rate.Limit(rl), rb)
	handler := limiter.Middleware(mux)

	if devMode {
		serve(handler)
	} else {
		serveTLS(handler)
	}

}

func serve(h http.Handler) {
	loclog := "[server.serve]"
	host := env.DevHost.Get()
	port := env.DevPort.Get()
	slog.Info(loclog, "info", "serving on", "host", host, "port", port)
	err := http.ListenAndServe(host+":"+port, h)
	if err != nil {
		slog.Error(loclog, "error", "serving on", "host", host, "port", port, "error", err.Error())
		os.Exit(1)
	}
}

func serveTLS(h http.Handler) {
	loclog := "[server.serveTLS]"
	allowedOrigins := strings.Split(env.AllowedOrigins.Get(), ",")
	allowedMethods := strings.Split(env.AllowedMethods.Get(), ",")
	handler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(h)

	host := env.Host.Get()
	port := env.Port.Get()
	tlsc := env.TLSCertPath.Get()
	tlsk := env.TLSKeyPath.Get()

	slog.Info(loclog, "info", "serving on", "host", host, "port", port)
	err := http.ListenAndServeTLS(host+":"+port, tlsc, tlsk, handler)
	if err != nil {
		slog.Error(loclog, "error", "serving on", "host", host, "port", port, "error", err.Error())
		os.Exit(1)
	}
}
