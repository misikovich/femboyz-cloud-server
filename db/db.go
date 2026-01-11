package db

import (
	"database/sql"
	"encoding/json"
	"femboyz/env"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// files table (id, pub_id (unique), meta (json), creation_date (timestamp), issuer, ref_view (integer), ref_dl (integer))
	filesStmt = `CREATE TABLE IF NOT EXISTS files (
				id 				INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
				pub_id 			TEXT NOT NULL UNIQUE, 
				meta 			TEXT NOT NULL, 
				creation_date 	TEXT DEFAULT (strftime('%s', 'now')),
				issuer 			TEXT NOT NULL, 
				ref_view 		INTEGER DEFAULT 0, 
				ref_dl 			INTEGER DEFAULT 0
				);`
	// posts table (id, pub_id (unique), content (json), creation_date (timestamp), issuer, ref_view (integer))
	postsStmt = `CREATE TABLE IF NOT EXISTS posts (
				id 				INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
				pub_id 			TEXT NOT NULL UNIQUE, 
				content 		TEXT NOT NULL, 
				creation_date 	TEXT DEFAULT (strftime('%s', 'now')), 
				issuer 			TEXT NOT NULL, 
				ref_view 		INTEGER DEFAULT 0
				);`
)

var db *sql.DB

func openDB() *sql.DB {
	loclog := "[db.openDB]"
	path := env.DBPath.Get()
	slog.Info(loclog, "pathdatabase", path)
	_db, err := sql.Open("sqlite3", path)
	if err != nil {
		slog.Error(loclog, "FATAL", "failed to open database", "error", err.Error())
		os.Exit(1)
	}
	return _db
}

func InitDB() {
	loclog := "[db.InitDB]"
	db = openDB()

	_, err := db.Exec(filesStmt)
	if err != nil {
		slog.Error(loclog, "FATAL", "failed to create table 'files'", "error", err.Error())
		os.Exit(1)
	}
	slog.Info(loclog, "info", "table 'files' executed")
	_, err = db.Exec(postsStmt)
	if err != nil {
		slog.Error(loclog, "FATAL", "failed to create table 'posts'", "error", err.Error())
		os.Exit(1)
	}
	slog.Info(loclog, "info", "table 'posts' executed")
	slog.Info(loclog, "info", "database initialized")
}

type FileMeta struct {
	OriginalName  string `json:"original_name"`
	Size          int64  `json:"size"`
	Hash          string `json:"hash"`
	LocalFileName string `json:"local_file_name"`
	FileType      string `json:"file_type"`
}

type File struct {
	ID           int64
	PubID        string
	Meta         FileMeta
	CreationDate string
	Issuer       string
	RefView      int
	RefDL        int
}

type Post struct {
	ID           int64
	PubID        string
	Content      string
	CreationDate string
	Issuer       string
	RefView      int
}

func InsertFile(f *File) error {
	loclog := "[db.InsertFile]"
	jsonMeta, err := json.Marshal(f.Meta)
	if err != nil {
		slog.Error(loclog, "SEVERE", "failed to marshal file meta", "error", err.Error(), "pub_id", f.PubID, "meta", f.Meta, "issuer", f.Issuer)
		return err
	}
	result, err := db.Exec("INSERT INTO files (pub_id, meta, issuer) VALUES (?, ?, ?)", f.PubID, jsonMeta, f.Issuer)
	if err != nil {
		slog.Error(loclog, "SEVERE", "failed to insert file in files table", "error", err.Error(), "pub_id", f.PubID, "meta", f.Meta, "issuer", f.Issuer)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		slog.Error(loclog, "SEVERE", "failed to get last insert id", "error", err.Error(), "pub_id", f.PubID, "meta", f.Meta, "issuer", f.Issuer)
		return err
	}
	f.ID = id
	slog.Info(loclog, "info", "file inserted in files table", "pubID", f.PubID)
	return nil
}

func GetFileByPubID(pubID string) (*File, error) {
	loclog := "[db.GetFileByPubID]"
	row := db.QueryRow("SELECT id, pub_id, meta, creation_date, issuer, ref_view, ref_dl FROM files WHERE pub_id = ?", pubID)

	var f File
	var jsonMeta []byte
	err := row.Scan(&f.ID, &f.PubID, &jsonMeta, &f.CreationDate, &f.Issuer, &f.RefView, &f.RefDL)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Debug(loclog, "info", "file not found", "pub_id", pubID)
			return nil, nil // Return nil if not found
		}
		slog.Error(loclog, "SEVERE", "failed to scan file", "error", err.Error(), "pub_id", pubID)
		return nil, err
	}
	json.Unmarshal(jsonMeta, &f.Meta)
	return &f, nil
}

func GetFileByID(id int64) (*File, error) {
	loclog := "[db.GetFileByID]"
	row := db.QueryRow("SELECT id, pub_id, meta, creation_date, issuer, ref_view, ref_dl FROM files WHERE id = ?", id)

	var f File
	var jsonMeta []byte
	err := row.Scan(&f.ID, &f.PubID, &jsonMeta, &f.CreationDate, &f.Issuer, &f.RefView, &f.RefDL)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Debug(loclog, "info", "file not found", "id", id)
			return nil, nil // Return nil if not found
		}
		slog.Error(loclog, "SEVERE", "failed to scan file", "error", err.Error(), "id", id)
		return nil, err
	}
	json.Unmarshal(jsonMeta, &f.Meta)
	return &f, nil
}

func InsertPost(p *Post) error {
	loclog := "[db.InsertPost]"
	_, err := db.Exec("INSERT INTO posts (pub_id, content, issuer) VALUES (?, ?, ?)", p.PubID, p.Content, p.Issuer)
	if err != nil {
		slog.Error(loclog, "SEVERE", "failed to insert post in posts table", "error", err.Error(), "pub_id", p.PubID, "content", p.Content, "issuer", p.Issuer)
		return err
	}

	slog.Info(loclog, "info", "post inserted in posts table", "pubID", p.PubID, "content", p.Content, "issuer", p.Issuer)
	return nil
}

func GetPostByPubID(pubID string) (*Post, error) {
	loclog := "[db.GetPostByPubID]"
	row := db.QueryRow("SELECT id, pub_id, content, creation_date, issuer, ref_view FROM posts WHERE pub_id = ?", pubID)

	var p Post
	var jsonContent []byte
	err := row.Scan(&p.ID, &p.PubID, &jsonContent, &p.CreationDate, &p.Issuer, &p.RefView)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Debug(loclog, "info", "post not found", "pub_id", pubID)
			return nil, nil // Return nil if not found
		}
		slog.Error(loclog, "SEVERE", "failed to scan post", "error", err.Error(), "pub_id", pubID)
		return nil, err
	}
	json.Unmarshal(jsonContent, &p.Content)
	slog.Info(loclog, "info", "post found", "pub_id", pubID, "returning", p)
	return &p, nil
}

func GetFileEntries() (int64, error) {
	loclog := "[db.GetFileEntries]"
	row := db.QueryRow("SELECT COUNT(*) FROM files")
	var count int64
	err := row.Scan(&count)
	if err != nil {
		slog.Error(loclog, "SEVERE", "failed to get file entries", "error", err.Error())
		return 0, err
	}
	return count, nil
}

func GetPostEntries() (int64, error) {
	loclog := "[db.GetPostEntries]"
	row := db.QueryRow("SELECT COUNT(*) FROM posts")
	var count int64
	err := row.Scan(&count)
	if err != nil {
		slog.Error(loclog, "SEVERE", "failed to get post entries", "error", err.Error())
		return 0, err
	}
	return count, nil
}
