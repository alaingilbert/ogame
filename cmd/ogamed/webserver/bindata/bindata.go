package bindata

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path"
)

// content holds our static web server content.
//go:embed image/*
//go:embed template/*
//go:embed assets/*
//go:embed html/index.html
var content embed.FS

func GetContent(path string) http.FileSystem {
	log.Print("using embed mode")

	fsys, err := fs.Sub(content, path)
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func GetDir(d string) ([]fs.DirEntry, error) {
	return content.ReadDir(d)

}

func GetFile(f string) ([]byte, error) {
	return content.ReadFile(f)
}

func GetAllDirectories() {
	// Content
	var res []string
	files, _ := content.ReadDir(".")
	res = getSubDir(files, "", res)
	for _, f := range res {
		fmt.Println(f)
	}
}

func getSubDir(s []fs.DirEntry, dirName string, res []string) []string {
	for _, item := range s {
		if item.IsDir() {
			subDir, _ := content.ReadDir(path.Join(dirName, item.Name()))
			res = getSubDir(subDir, path.Join(dirName, item.Name()), res)
		} else {
			res = append(res, path.Join(dirName, item.Name()))
		}
	}
	return res
}
