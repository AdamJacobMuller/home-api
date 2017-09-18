package apicontroller

import (
	"gopkg.in/hlandau/passlib.v1"
)

type Authentication struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func (c *APIController) AuthorizeUsernamePassword(username string, password string) bool {
	for _, auth := range c.authentication {
		if auth.Username == username {
			_, err := passlib.Verify(password, auth.Password)
			if err != nil {
				return false
			}
			return true
		}
	}
	return false
}
