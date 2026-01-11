package handlers

import (
	"encoding/json"
	"femboyz/db"
	"femboyz/env"
	"femboyz/uidgenerator"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var timeStarted time.Time

func Init() {
	timeStarted = time.Now()
}

type Health struct {
	Status  string  `json:"status"`
	Uptime  string  `json:"uptime"`
	Entries []int64 `json:"entries"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	loclog := "[handlers.HealthCheck]"
	ip := getRequestIP(r)
	slog.Info(loclog, "info", "health check request", "method", r.Method, "ip", ip)
	// if not GET - drop connection
	if r.Method != http.MethodGet {
		slog.Warn(loclog, "warning", "health check request method not GET", "method", r.Method, "ip", ip)
		return
	}

	// check for token
	token := r.Header.Get("Authorization")
	if token != env.HealthCheckToken.Get() {
		slog.Warn(loclog, "warning", "health check request token not match", "token", token, "ip", ip)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getHealth())
}

func getRequestIP(r *http.Request) string {
	proxyIP := r.Header.Get("X-Forwarded-For")
	if proxyIP != "" {
		return proxyIP
	}
	return r.RemoteAddr
}

func getHealth() Health {
	ut := time.Since(timeStarted)
	fileEntries, _ := db.GetFileEntries()
	postEntries, _ := db.GetPostEntries()
	return Health{
		Status: "ok",
		Uptime: ut.String(),
		Entries: []int64{
			fileEntries,
			postEntries,
		},
	}
}

type File struct {
	Filename string `json:"filename"`
	Filesize int64  `json:"filesize"`
	Filetype string `json:"filetype"`
	Filehash string `json:"filehash"`
	Fileurl  string `json:"fileurl"`
}

func PullFile(w http.ResponseWriter, r *http.Request) {
	loclog := "[handlers.PullFile]"
	ip := getRequestIP(r)
	slog.Info(loclog, "info", "pull file request", "method", r.Method, "ip", ip)
	// if not GET - drop connection
	if r.Method != http.MethodGet {
		slog.Warn(loclog, "warning", "pull file request method not GET", "method", r.Method, "ip", ip)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		slog.Warn(loclog, "warning", "pull file request id not provided", "ip", ip)
		return
	}

	if !uidgenerator.Validate(id) {
		slog.Warn(loclog, "warning", "pull file request id not valid", "id", id, "ip", ip)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := db.GetFileByPubID(id)
	if err != nil {
		slog.Error(loclog, "error", "pull file request failed", "id", id, "ip", ip, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if f == nil {
		slog.Warn(loclog, "warning", "pull file request file not found", "id", id, "ip", ip)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// f.Meta is a json struct that may contain metadata about the file and different types of files may have different metadata
	// for example, image files may have metadata
	fmeta := f.Meta
	sendMeta := map[string]interface{}{
		"creation_date": f.CreationDate,
		"filename":      fmeta.OriginalName,
		"filesize":      fmeta.Size,
		"filetype":      fmeta.FileType,
		"filehash":      fmeta.Hash,
		"file_pub_id":   f.PubID,
		"views":         f.RefView,
		"downloads":     f.RefDL,
	}

	fileData, err := os.ReadFile("testblobs/" + fmeta.LocalFileName)
	if err != nil {
		slog.Error(loclog, "error", "pull file request failed to read file", "id", id, "ip", ip, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mw := multipart.NewWriter(w)
	w.Header().Set("Content-Type", mw.FormDataContentType())

	metaPart, _ := mw.CreateFormFile("metadata", "metadata.json")
	json.NewEncoder(metaPart).Encode(sendMeta)

	filePart, _ := mw.CreateFormFile("file", fmeta.OriginalName)
	filePart.Write(fileData)

	w.WriteHeader(http.StatusOK)
	mw.Close()
}

func Admin(w http.ResponseWriter, r *http.Request) {

}

func FilePage(w http.ResponseWriter, r *http.Request) {

}

func PostPage(w http.ResponseWriter, r *http.Request) {

}

func Send(w http.ResponseWriter, r *http.Request) {

}

func PullPost(w http.ResponseWriter, r *http.Request) {

}
