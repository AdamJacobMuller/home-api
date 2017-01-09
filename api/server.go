package api

import (
	"encoding/json"
	"github.com/AdamJacobMuller/home-api/providers/homeseer"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"
	"io/ioutil"
	"net/http"
)

type ControlRequest struct {
	Match homeseer.HSLookup `json:"match"`
	Value float64           `json:"value"`
}

type APIServer struct {
	Server       *http.Server
	HSController *homeseer.HSController
}

func (a *APIServer) readAndUnmarshalJson(r *http.Request, object interface{}) bool {
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
func (a *APIServer) marshalAndWriteJson(w http.ResponseWriter, object interface{}) bool {
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
func (a *APIServer) SetChildDevicesValue(w http.ResponseWriter, r *http.Request) {
	cr := &ControlRequest{}
	a.readAndUnmarshalJson(r, cr)
	a.HSController.SetChildDevicesValue(cr.Match, cr.Value)
}
func (a *APIServer) SetDevicesValue(w http.ResponseWriter, r *http.Request) {
	cr := &ControlRequest{}
	a.readAndUnmarshalJson(r, cr)
	a.HSController.SetDevicesValue(cr.Match, cr.Value)
}
func (a *APIServer) Serve() {
	err := a.Server.ListenAndServe()
	if err != nil {
		log.Errorln(err)
	}
}

func NewAPIServer() *APIServer {
	apiserver := &APIServer{}
	n := negroni.New()
	recovery := negroni.NewRecovery()
	n.Use(recovery)
	n.Use(negronilogrus.NewMiddleware())
	handler := mux.NewRouter()
	n.UseHandler(handler)

	handler.HandleFunc("/SetDevicesValue", apiserver.SetDevicesValue)
	handler.HandleFunc("/SetChildDevicesValue", apiserver.SetChildDevicesValue)

	apiserver.Server = &http.Server{Handler: n, Addr: "0.0.0.0:8145"}
	return apiserver
}
