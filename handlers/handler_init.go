package handlers

import (
	"dcrvnwww/helper"
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
		Cookie:  "dcrvnwww",
		Expires: time.Hour * 24 * 365,
	})
	route.Use(initMiddleware)
	route.Handle("/asset/", http.StripPrefix("/asset", http.FileServer(http.Dir("public"))))
	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmplHelper.Render(w,r,"home", nil)
	})
	//Group for admin
	route.Route("/admin", func(r chi.Router) {
		r.Get("/", adminIndex)
		r.Get("/users", userManage)
		r.Get("/post", postManage)
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

func adminIndex(w http.ResponseWriter, r *http.Request)  {
	tmplHelper.Render(w, r, "main_admin", nil)
}

func userManage(w http.ResponseWriter, r *http.Request)  {
	tmplHelper.Render(w, r, "admin_user_list", nil)
}

func postManage(w http.ResponseWriter, r *http.Request)  {
	tmplHelper.Render(w, r, "admin_post_list", nil)
}
