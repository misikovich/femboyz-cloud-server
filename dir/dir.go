package dir

import (
	"femboyz/env"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var (
	MainDB = env.DBPath.Get()
)

var dirs = []string{
	"files",
	"logs",
	"db",
}

var files = []string{
	MainDB,
}

var currentExecPath string

func GetExecPath() string {
	return currentExecPath
}

func InitDir() {
	loclog := "[dir.init]"
	p, err := os.Executable()
	if err != nil {
		panic(err)
	}
	currentExecPath = p[:strings.LastIndex(p, "/")]
	slog.Info(loclog, "pwd", currentExecPath)
}

func makeDir(path string) {
	path = filepath.Join(currentExecPath, path)
	loclog := "[dir.makeDir]"
	// check if path exists
	_, err := os.Stat(path)
	if err == nil {
		slog.Info(loclog, "pathexists", path)
		return
	}
	// create path
	err = os.MkdirAll(path, 0755)
	if err != nil {
		slog.Error(loclog, "err", err, "path", path)
		os.Exit(1)
	} else {
		slog.Info(loclog, "pathcreated", path)
	}
}

func makeFile(path string) {
	path = filepath.Join(currentExecPath, path)
	loclog := "[dir.makeFile]"
	// check if path exists
	_, err := os.Stat(path)
	if err == nil {
		slog.Info(loclog, "fileexists", path)
		return
	}
	// create path
	f, err := os.Create(path)
	if err != nil {
		slog.Error(loclog, "err", err, "path", path)
		os.Exit(1)
	} else {
		f.Close()
		slog.Info(loclog, "filecreated", path)
	}
}

func MakeDirs() {
	for _, dir := range dirs {
		makeDir(dir)
	}
}

func MakeFiles() {
	for _, file := range files {
		makeFile(file)
	}
}
