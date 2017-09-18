package apiserver

import (
	"encoding/base64"
	"fmt"
	"github.com/gocraft/web"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func (c *Context) Authenticate(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	cookie, err := r.Cookie("auth")
	if err != nil && err != http.ErrNoCookie {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to look up cookie")
		goto CheckBasicAuth
	}
	if cookie != nil {
		data, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("unable to decode cookie")
			goto CheckBasicAuth
		}
		split := strings.SplitN(string(data), ",", 2)
		if len(split) != 2 {
			goto CheckBasicAuth
		}
		if c.apiserver.Controller.AuthorizeUsernamePassword(split[0], split[1]) {
			next(w, r)
			return
		}
	}
CheckBasicAuth:
	username, password, ok := r.BasicAuth()
	if ok {
		if c.apiserver.Controller.AuthorizeUsernamePassword(username, password) {
			cookie := &http.Cookie{
				Name:    "auth",
				Value:   base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))),
				Expires: time.Now().Add(time.Hour * 24 * 365 * 10),
			}
			http.SetCookie(w, cookie)
			next(w, r)
			return
		}
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="API"`)
	w.WriteHeader(401)
}
