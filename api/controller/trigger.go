package apicontroller

import (
	"github.com/AdamJacobMuller/home-api/api/models"

	log "github.com/sirupsen/logrus"
)

type Trigger struct {
	Matches []TriggerMatch `json:"Matches"`
	Actions []apimodels.ControlRequest
}

type TriggerMatch struct {
	Type       string `json:"Type"`
	ProviderID string `json:"ProviderID"`
	DeviceID   string `json:"DeviceID"`
	Action     string `json:"Action"`
	Stage      string `json:"Stage"`
}

type Triggers struct {
	Triggers []Trigger
}

func (t *Triggers) InvokeDevicesAction(controller *APIController, stage string, match apimodels.Match, action string) {
	providerID, haveProviderID := match["ProviderID"]
	deviceID, haveDeviceID := match["DeviceID"]
	var result bool
	for _, trigger := range t.Triggers {
		for _, triggerMatch := range trigger.Matches {
			result = true
			if triggerMatch.Stage != "" {
				if triggerMatch.Stage != stage {
					result = false
				}
			}
			if triggerMatch.Type != "" {
				if triggerMatch.Type != "InvokeDevicesAction" {
					result = false
				}
			}
			if triggerMatch.ProviderID != "" && haveProviderID {
				if triggerMatch.ProviderID != providerID {
					result = false
				}
			}
			if triggerMatch.DeviceID != "" && haveDeviceID {
				if triggerMatch.DeviceID != deviceID {
					result = false
				}
			}
			if triggerMatch.Action != "" {
				if triggerMatch.Action != action {
					result = false
				}
			}
			log.WithFields(log.Fields{
				"stage":                    stage,
				"triggerMatch.Stage":       triggerMatch.Stage,
				"haveProviderID":           haveProviderID,
				"haveDeviceID":             haveDeviceID,
				"trigger.Match.DeviceID":   triggerMatch.DeviceID,
				"trigger.Match.ProviderID": triggerMatch.ProviderID,
				"providerID":               providerID,
				"deviceID":                 deviceID,
				"result":                   result,
			}).Info("checked trigger")
			if result {
				for _, action := range trigger.Actions {
					ar := controller.ControlRequest(action)
					log.WithFields(log.Fields{"action": action, "ar": ar}).Info("executing action")
				}
			}
		}
	}
}
