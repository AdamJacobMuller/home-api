package friedrichlink

import (
	"encoding/json"
	"fmt"
	"github.com/AdamJacobMuller/home-api/api/models"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

type FriedrichLink struct {
	username string
	password string
	client   *http.Client
}

func (p *FriedrichLink) SetDevicesValue(apimodels.Match, float64) bool {
	return false
}
func (p *FriedrichLink) SetChildDevicesValue(apimodels.Match, float64) bool {
	return false
}
func (p *FriedrichLink) InvokeDevicesAction(apimodels.Match, string) bool {
	return false
}
func (p *FriedrichLink) InvokeChildDevicesAction(apimodels.Match, string) bool {
	return false
}
func (p *FriedrichLink) GetDevices(apimodels.Match) (apimodels.Devices, bool) {
	return &FriedrichLinkList{}, false
}
func (p *FriedrichLink) GetDevice(apimodels.Match) (apimodels.Device, bool) {
	return &FriedrichLinkDevice{}, false
}
func (p *FriedrichLink) GetChildDevice(apimodels.Match) (apimodels.Device, bool) {
	return &FriedrichLinkDevice{}, false
}
func (p *FriedrichLink) GetChildDevices(apimodels.Match) (apimodels.Devices, bool) {
	return &FriedrichLinkList{}, false
}

type FriedrichLinkConfig struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func (p *FriedrichLink) Create(rawConfig json.RawMessage) bool {
	c := &FriedrichLinkConfig{}
	err := json.Unmarshal(rawConfig, c)
	if err != nil {
		panic(err)
		return false
	}
	p.username = c.Username
	p.password = c.Password

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}

	p.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	p.login()
	p.load()

	return true
}
func (p *FriedrichLink) load() error {
	data, err := p.get("https://friedrichlink.friedrich.com/user_pages.php")
	if err != nil {
		panic(err)
	}
	ndata := regexp.MustCompile("var unit = new Units\\(({.*})\\);").FindStringSubmatch(data)
	fmt.Printf("%s\n", ndata[1])
	fmt.Printf("len(ndata) == %d\n", len(ndata))
	fmt.Printf("p.client.Jar: %+v\n", p.client.Jar)
	return nil
}

func (p *FriedrichLink) get(url string) (string, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return "", err
	}
	return string(body), nil
}

func (p *FriedrichLink) login() error {
	resp, err := p.client.PostForm("http://friedrichlink.friedrich.com/actions/login.php",
		url.Values{
			"login_userName": {p.username},
			"login_password": {p.password},
		})
	fmt.Printf("resp: %+v\n", resp)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return err
	}
	fmt.Printf("Body: %s\n", body)

	return nil
}

func (p *FriedrichLink) TypeString() string {
	return "FriedrichLinkType"
}
func (p *FriedrichLink) IDString() string {
	return "FriedrichLinkInstance"
}
