package storage

import (
	"io"
	"os"
	"path/filepath"

	beego "github.com/beego/beego/v2/server/web"
)

func BaseDir() string {
	dir, _ := beego.AppConfig.String("storage_dir")
	if dir == "" {
		dir = "storage"
	}
	return dir
}

func EnsureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}

func SaveReader(relPath string, r io.Reader) (string, error) {
	full := filepath.Join(BaseDir(), relPath)
	if err := EnsureDir(full); err != nil {
		return "", err
	}
	f, err := os.Create(full)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	return relPath, nil
}

func FullPath(rel string) string {
	return filepath.Join(BaseDir(), rel)
}

func URLPath(rel string) string {
	if rel == "" {
		return ""
	}
	return "/static/" + rel
}
