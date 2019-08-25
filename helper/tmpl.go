package helper

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type (
	TmplHelper struct {
		tmpl   *template.Template
		config TmplConfig
	}
	TmplConfig struct {
		Dir      string
		Suffix   string
		NotFound string
		ProcessData func(*http.Request, map[string]interface{},) map[string]interface{}
	}
)

var (
	cacheKey = ""
)

func init() {
	cacheKey = RandStringBytes(10)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewTPL(config TmplConfig, funcMap template.FuncMap) (tmpl *TmplHelper, err error) {
	tmpl = &TmplHelper{}
	validFiles := []string{}
	err = filepath.Walk(config.Dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, "."+config.Suffix) {
			validFiles = append(validFiles, path)
		}
		return nil
	})
	if err != nil {
		return
	}
	funcMap["uuid"] = func() string {
		return cacheKey
	}
	tmpl.tmpl = template.Must(template.New("tmpl").Funcs(funcMap).ParseFiles(validFiles...))
	tmpl.config = config
	return
}

func (t TmplHelper) Render(wr io.Writer, r *http.Request, name string, data map[string]interface{}) {
	if t.config.ProcessData != nil {
		data = t.config.ProcessData(r,data)
	}
	data["uuid"] = cacheKey
	err := t.tmpl.ExecuteTemplate(wr, name, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (t TmplHelper) NotFound(wr io.Writer, r *http.Request) {
	t.Render(wr, r, t.config.NotFound, nil)
}
