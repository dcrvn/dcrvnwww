package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/dcrvn/dcrvnwww/models"
)

func initMiddleware(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := getUserSession(w, r)
		data := map[string]interface{}{
			"user": user,
		}
		ctx := context.WithValue(r.Context(), "data", data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func needLoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := getUserSession(w, r)
		if err != nil {
			http.Redirect(w, r, "/auth/login", 302)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loginUserSession(w http.ResponseWriter, r *http.Request, user models.User) {
	store := sessionManager.Start(w, r)
	user.LastLogin = time.Now()
	user.LoginTrack()
	store.Set("user", user)
}

func getUserSession(w http.ResponseWriter, r *http.Request) (user models.User, err error) {
	store := sessionManager.Start(w, r)
	var ok bool
	user, ok = store.Get("user").(models.User)
	if !ok {
		err = fmt.Errorf("Can not get user data")
	}
	return
}
