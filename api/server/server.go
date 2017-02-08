package apiserver

import (
	"encoding/json"
	"github.com/AdamJacobMuller/home-api/api/controller"
	"github.com/AdamJacobMuller/home-api/api/models"
	"github.com/AdamJacobMuller/weblogrus"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"io/ioutil"
	"net/http"
)

type APIServer struct {
	Server     *http.Server
	Controller *apicontroller.APIController
}

func (a *APIServer) readAndUnmarshalJson(r *web.Request, object interface{}) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("unable to read body")
		return false
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("unable to unmarshal json")
		return false
	}
	return true
}
func (a *APIServer) marshalAndWriteJson(w web.ResponseWriter, object interface{}) bool {
	jsonDocument, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(503)
		log.Errorln(err)
		return false
	}
	w.WriteHeader(200)
	w.Header().Add("Content-type", "application/json")
	w.Write(jsonDocument)
	return true
}
func (a *APIServer) ListDevices(w web.ResponseWriter, r *web.Request) {
	cr := a.Controller.ListDevices()
	a.marshalAndWriteJson(w, cr)
}
func (a *APIServer) SetChildDevicesValue(w web.ResponseWriter, r *web.Request) {
	cr := &apimodels.ControlRequest{}
	cr.Value = -1
	ok := a.readAndUnmarshalJson(r, cr)
	if !ok {
		return
	}
	if cr.Value == -1 {
		return
	}
	a.Controller.SetChildDevicesValue(cr.Match, cr.Value)
}
func (a *APIServer) SetDevicesValue(w web.ResponseWriter, r *web.Request) {
	cr := &apimodels.ControlRequest{}
	cr.Value = -1
	ok := a.readAndUnmarshalJson(r, cr)
	if !ok {
		return
	}
	if cr.Value == -1 {
		return
	}
	log.WithFields(log.Fields{"ControlRequest": cr}).Info("SetDevicesValue")
	a.Controller.SetDevicesValue(cr.Match, cr.Value)
}
func (a *APIServer) InvokeChildDevicesAction(w web.ResponseWriter, r *web.Request) {
	cr := &apimodels.ControlRequest{}
	ok := a.readAndUnmarshalJson(r, cr)
	if !ok {
		return
	}
	a.Controller.InvokeChildDevicesAction(cr.Match, cr.Action)
}
func (a *APIServer) InvokeDevicesAction(w web.ResponseWriter, r *web.Request) {
	cr := &apimodels.ControlRequest{}
	ok := a.readAndUnmarshalJson(r, cr)
	if !ok {
		return
	}
	a.Controller.InvokeDevicesAction(cr.Match, cr.Action)
}
func (a *APIServer) LoadAndCreateProviders(filename string) bool {
	return a.Controller.LoadAndCreateProviders(filename)
}
func (a *APIServer) Serve() {
	err := a.Server.ListenAndServe()
	if err != nil {
		log.Errorln(err)
	}
}

func NewAPIServer() *APIServer {
	apiserver := &APIServer{}
	apiserver.Controller = apicontroller.NewAPIController()

	router := web.New(Context{})

	router.NotFound((*Context).NotFound)

	router.Middleware(func(ctx *Context, resp web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {
		ctx.apiserver = apiserver
		next(resp, req)
	})
	x := weblogrus.NewMiddleware()
	router.Middleware(x.ServeHTTP)
	router.Middleware(web.ShowErrorsMiddleware)

	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateFS, rice.LocateEmbedded, rice.LocateAppended},
	}
	var err error
	var box *rice.Box
	var directory http.FileSystem

	if false {
		directory = http.Dir("files/development")
	} else {
		box, err = conf.FindBox("files/production")
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Errorf("unable to find rice box")
			directory = http.Dir("files/development")
		} else {
			directory = box.HTTPBox()
		}
	}
	router.Middleware(web.StaticMiddlewareFromDir(directory, web.StaticOption{Prefix: "/static", IndexFile: "index.html"}))

	admin := router.Subrouter(Context{}, "/")
	admin.Middleware((*Context).DrawLayout)
	admin.Middleware((*Context).AddControllerNavigation)
	admin.Get("/", (*Context).Index)
	admin.Get("/room/:LocationTwo/:LocationOne", (*Context).Room)

	api := router.Subrouter(Context{}, "/api")
	api.Get("/ListDevices", apiserver.ListDevices)
	api.Post("/SetDevicesValue", apiserver.SetDevicesValue)
	api.Post("/SetChildDevicesValue", apiserver.SetChildDevicesValue)
	api.Post("/InvokeDevicesAction", apiserver.InvokeDevicesAction)
	api.Post("/InvokeChildDevicesAction", apiserver.InvokeChildDevicesAction)

	apiserver.Server = &http.Server{Handler: router, Addr: "0.0.0.0:8145"}
	return apiserver
}
