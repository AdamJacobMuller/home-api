package friedrichlink

import (
	"encoding/json"
	"regexp"
)

type FLAPI struct {
	UserID  string                 `json:"userId"`
	Serial  string                 `json:"serialNumber"`
	Name    string                 `json:"nickName"`
	Ambient string                 `json:"ambientAir"`
	Model   string                 `json:"modelNumber"`
	Status  map[string]FLAPIStatus `json:"status"`
}

type FLAPIStatus struct {
	DisplayName  string `json:"parameter_display_name"`
	DisplayValue string `json:"display_value"`
	Name         string `json:"parameter_name"`
	Value        string `json:"value"`
}

func (p *FriedrichLink) load() error {
	data, err := p.get("https://friedrichlink.friedrich.com/user_pages.php")
	if err != nil {
		panic(err)
	}
	ndata := regexp.MustCompile("var unit = new Units\\(({.*})\\);").FindStringSubmatch(data)

	fla := &FLAPI{}

	err = json.Unmarshal([]byte(ndata[1]), fla)
	if err != nil {
		panic(err)
	}
	p.UserID = fla.UserID

	device := &FriedrichLinkDevice{
		Name:   fla.Name,
		Serial: fla.Serial,
		Model:  fla.Model,
		fl:     p,
	}

	for _, o := range fla.Status {
		switch o.Name {
		case "ambient_air_temperature":
			device.Ambient = o.Value
		case "auto_cool_set_point":
			device.CoolPoint = o.Value
		case "power":
			device.Power = o.Value
		}
	}

	p.devices = []*FriedrichLinkDevice{device}
	return nil
}
