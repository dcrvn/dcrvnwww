package handlers

import (
	"dcrvnwww/helper"
	"github.com/kataras/go-sessions"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

	// Init asset
	workDir, _ := os.Getwd()
	node_modules := filepath.Join(workDir, "node_modules")
	asset := filepath.Join(workDir, "public")
	FileServer(route, "/node_modules/", http.Dir(node_modules))
	FileServer(route, "/asset/", http.Dir(asset))

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
	tmplHelper.Render(w, r, "admin_dashboard", nil)
}

func userManage(w http.ResponseWriter, r *http.Request)  {
	tmplHelper.Render(w, r, "admin_user_list", nil)
}

func postManage(w http.ResponseWriter, r *http.Request)  {
	tmplHelper.Render(w, r, "admin_post_list", nil)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
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
