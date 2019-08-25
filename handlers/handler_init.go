package handlers

import (
	"github.com/dcrvn/dcrvnwww/helper"
	"github.com/kataras/go-sessions"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"reflect"
	"time"
)

type (
	Map map[string]interface{}
)

var (
	sessionManager *sessions.Sessions
	tmplHelper     *helper.TmplHelper
)

func Init(route *chi.Mux)  {
	tmplHelper, _ = helper.NewTPL(helper.TmplConfig{
		Dir:"tpl",
		Suffix: "html",
		NotFound: "error_notFound",
		ProcessData: func(r *http.Request, data map[string]interface{}) map[string]interface{} {
			if data == nil {
				data = make(map[string]interface{})
			}
			ctx := r.Context()
			if ctxData,ok := ctx.Value("data").(map[string]interface{}); ok {
				if _, ok = data["user"]; !ok {
					data["user"] = ctxData["user"]
				}
			}
			return data
		},
	}, funcMap())
	sessionManager = sessions.New(sessions.Config{
		Cookie:  "smgs",
		Expires: time.Hour * 24 * 365,
	})
	route.Use(initMiddleware)
	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmplHelper.Render(w,r,"home", nil)
	})
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"html": func(data string) template.HTML {
			return template.HTML(data)
		},
		"empty": func(x interface{}) bool {
			return x == reflect.Zero(reflect.TypeOf(x)).Interface()
		},
		"minus": func(a, b int) int {
			return a - b
		},
	}
}
