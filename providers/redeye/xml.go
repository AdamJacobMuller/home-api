package redeye

import (
	"fmt"
	"github.com/AdamJacobMuller/home-api/api/models"
	log "github.com/Sirupsen/logrus"
	"regexp"
)

// <redeye serialNumber="E0102-80639" wifiMacAddress="00:23:87:00:5c:31" lanMacAddress="00:23:87:00:5c:32" hardwareType="RedEye Pro" description="" name="RedEye Pro">
type RedEyeRoot struct {
	SerialNumber   string `xml:"serialNumber,attr"`
	WifiMacAddress string `xml:"wifiMacAddress,attr"`
	LanMacAddress  string `xml:"lanMacAddress,attr"`
	HardwareType   string `xml:"hardwareType,attr"`
	Description    string `xml:"description,attr"`
	Name           string `xml:"name,attr"`
}

type RedEyeRooms struct {
	Rooms []*RedEyeRoom `xml:"room"`
}

// <room name="RedEye Pro" currentActivityId="-1" roomId="-1" description=""/>
type RedEyeRoom struct {
	Name              string `xml:"name,attr"`
	CurrentActivityID string `xml:"currentActivityId,attr"`
	RoomID            string `xml:"roomId,attr"`
	Description       string `xml:"description,attr"`
	Activities        []*RedEyeActivity
	RedEye            *RedEye
}

type RedEyeActivities struct {
	Activities []*RedEyeActivity `xml:"activity"`
}

// <activity name="Apple TV" activityId="365" description="" activityType="11"/>
type RedEyeActivity struct {
	Name         string `xml:"name,attr"`
	ActivityID   string `xml:"activityId,attr"`
	Description  string `xml:"description,attr"`
	ActivityType string `xml:"activityType,attr"`
	Room         *RedEyeRoom
}

func (p *RedEye) LoadData() {
	root := &RedEyeRoot{}
	p.API.GetAndUnmarshal("", root)
	fmt.Printf("ROOT: %+v\n", root)

	rooms := &RedEyeRooms{}
	p.API.GetAndUnmarshal("rooms", rooms)
	fmt.Printf("ROOMS: %+v\n", rooms)
	for _, room := range rooms.Rooms {
		room.RedEye = p
		fmt.Printf("ROOM: %+v\n", room)
		activities := &RedEyeActivities{}
		p.API.GetAndUnmarshal(fmt.Sprintf("rooms/%s/activities", room.RoomID), activities)
		room.Activities = activities.Activities
		fmt.Printf("ACTIVITIES: %+v\n", activities)
		for _, activity := range activities.Activities {
			fmt.Printf("ACTIVITY: %+v\n", activity)
			activity.Room = room
		}
	}
	p.Rooms = rooms.Rooms
}

func (room *RedEyeRoom) Matches(find apimodels.Match) bool {
	var rVal *regexp.Regexp
	var sVal string
	var ok bool
	var err error
	log.WithFields(log.Fields{"room": room, "find": find}).Info("Matches called")
	sVal, ok = find["RegexpName"].(string)
	if ok {
		rVal, err = regexp.Compile(sVal)
		if err != nil {
			log.WithFields(log.Fields{"regexp": sVal, "error": err}).Error("unable to compile regexp")
			return false
		}
		ok = rVal.MatchString(room.Name)
		if !ok {
			return false
		}
	}
	/*
		sVal, ok = find["LocationOne"].(string)
		if ok {
			if sVal != o.LocationOne {
				return false
			}
		}
	*/
	sVal, ok = find["LocationOne"].(string)
	if ok {
		if sVal != room.Name {
			return false
		}
	}
	sVal, ok = find["DeviceID"].(string)
	if ok {
		if sVal != room.RoomID {
			return false
		}
	}
	sVal, ok = find["Name"].(string)
	if ok {
		if sVal != "RedEye" {
			return false
		}
	}
	/*
		var lVal apimodels.Match
		lVal, ok = find["Child"].(apimodels.Match)
		if ok {
			var matched bool
			matched = false
			for _, child := range o.Children {
				if child.Matches(lVal) {
					matched = true
				}
			}
			if !matched {
				return false
			}
		}
	*/
	log.WithFields(log.Fields{"room": room, "find": find}).Info("Matches matches")
	return true
}

func (room *RedEyeRoom) HasChildDevice(find apimodels.Match) bool {
	return false
}
func (room *RedEyeRoom) GetChildDevice(find apimodels.Match) apimodels.Device {
	return &Device{}
}
func (room *RedEyeRoom) IsHidden() bool {
	return false
}
func (room *RedEyeRoom) ListChildren() []apimodels.Device {
	return []apimodels.Device{}
}
func (a *RedEyeActivity) GetName() string {
	return a.Name
}
func (room *RedEyeRoom) ListActions() []apimodels.Action {
	var actions []apimodels.Action

	if len(room.Activities) == 0 {
		return actions
	}

	actions = append(actions, &RedEyeActivity{Name: "Off"})

	for _, roomi := range room.Activities {
		actions = append(actions, roomi)
	}

	return actions
}
func (room *RedEyeRoom) ListTypes() []apimodels.Type {
	var t []apimodels.Type
	t = append(t, &RedEyeType{Name: "RedEye"})
	return t
}
func (room *RedEyeRoom) SetValue(float64) bool {
	return false
}
func (room *RedEyeRoom) InvokeAction(action string) bool {
	log.WithFields(log.Fields{"Action": action}).Info("InvokeAction")
	if action == "Off" {
		room.RedEye.API.Get(fmt.Sprintf("rooms/%s/activities/launch?activityId=-1", room.RoomID))
		return true
	}
	for _, activity := range room.Activities {
		if activity.Name == action {
			room.RedEye.API.Get(fmt.Sprintf("rooms/%s/activities/launch?activityId=%s", room.RoomID, activity.ActivityID))
		}
	}
	return true
}
func (room *RedEyeRoom) GetName() string {
	return "RedEye"
}
func (room *RedEyeRoom) GetLocationOne() string {
	return room.Name
}
func (room *RedEyeRoom) ProviderIDString() string {
	return room.RedEye.API.Base
}
func (room *RedEyeRoom) GetLocationTwo() string {
	return "1st Floor"
}
func (room *RedEyeRoom) IDString() string {
	return room.RoomID
}

type RedEyeType struct {
	Name string
}

func (t *RedEyeType) GetName() string {
	return t.Name
}
