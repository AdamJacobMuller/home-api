package redeye

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type API struct {
	Base string
}

func (c *API) GetAndUnmarshal(path string, object interface{}) error {
	var url string
	if path == "" {
		url = c.Base
	} else {
		url = fmt.Sprintf("%s%s", c.Base, path)
	}
	log.WithFields(log.Fields{"url": url}).Info("GetAndUnmarshal")
	resp, err := http.Get(url)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GET failed")
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("read Body failed")
		return err
	}
	resp.Body.Close()

	log.WithFields(log.Fields{"len": len(body)}).Info("Got Body")

	err = xml.Unmarshal(body, object)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("xml.Unmarshal failed")
		return err
	}
	return nil
}

func (c *API) Get(path string) ([]byte, error) {
	var url string
	if path == "" {
		url = c.Base
	} else {
		url = fmt.Sprintf("%s%s", c.Base, path)
	}
	log.WithFields(log.Fields{"url": url}).Info("Get")
	resp, err := http.Get(url)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GET failed")
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("read Body failed")
		return []byte{}, err
	}
	resp.Body.Close()
	return body, nil
}
