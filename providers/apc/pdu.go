package apc

import (
	"encoding/json"
	"fmt"
	"github.com/AdamJacobMuller/gosnmp"
	"github.com/AdamJacobMuller/home-api/api/models"
	log "github.com/Sirupsen/logrus"
	"net"
	"regexp"
	"strings"
	"time"
)

type SNMPTable struct {
	Values   map[string]map[string]gosnmp.SnmpPDU
	OneStart int
	OneEnd   int
	TwoStart int
	TwoEnd   int
}

func (s *SNMPTable) Print() {
	for k1, l := range s.Values {
		for k2, v := range l {
			fmt.Printf("%10s %10s %+v\n", k1, k2, v)
		}
	}
}

func (s *SNMPTable) Callback(pdu gosnmp.SnmpPDU) error {
	if s.Values == nil {
		s.Values = make(map[string]map[string]gosnmp.SnmpPDU)
	}
	oid := strings.Split(pdu.Name, ".")
	keyOne := strings.Join(oid[s.OneStart:s.OneEnd], ".")
	keyTwo := strings.Join(oid[s.TwoStart:s.TwoEnd], ".")
	if s.Values[keyOne] == nil {
		s.Values[keyOne] = make(map[string]gosnmp.SnmpPDU)
	}
	s.Values[keyOne][keyTwo] = pdu
	return nil
}

type SNMPList struct {
	Values map[string]gosnmp.SnmpPDU
	Start  int
	End    int
}

func (s *SNMPList) Print() {
	for k, v := range s.Values {
		fmt.Printf("%10s %#v\n", k, v.String())
	}
}

func (s *SNMPList) Callback(pdu gosnmp.SnmpPDU) error {
	if s.Values == nil {
		s.Values = make(map[string]gosnmp.SnmpPDU)
	}
	oid := strings.Split(pdu.Name, ".")
	key := strings.Join(oid[s.Start:s.End], ".")
	s.Values[key] = pdu
	return nil
}

type PDU struct {
	IP        net.IP
	SNMP      *gosnmp.GoSNMP
	Community string
	Name      string
	Location  string
	Outlets   []*Outlet
}

type Outlet struct {
	PDU   *PDU
	Index int64
	Name  string
	Phase int64
	State bool
}

func (o *Outlet) Print() {
	fmt.Printf("%s [#%d, Phase = %d, State = %t ]\n", o.Name, o.Index, o.Phase, o.State)
}

func (o *Outlet) GetTypes() []string {
	return []string{"Binary Switch"}
}
func (o *Outlet) ListChildren() []apimodels.Device {
	r := make([]apimodels.Device, 0)
	return r
}

/*
rPDUOutletControlOutletCommand OBJECT-TYPE
   SYNTAX INTEGER {
      immediateOn             (1),
      immediateOff            (2),
      immediateReboot         (3),
      delayedOn               (4),
      delayedOff              (5),
      delayedReboot           (6),
      cancelPendingCommand    (7)
   }
*/

func (o *Outlet) InvokeAction(action string) bool {
	switch action {
	case "immediateOn":
		return o.SetValue(1)
	case "on":
		return o.SetValue(1)
	case "immediateOff":
		return o.SetValue(2)
	case "off":
		return o.SetValue(2)
	case "immediateReboot":
		return o.SetValue(3)
	case "reboot":
		return o.SetValue(3)
	default:
		log.WithFields(log.Fields{"action": action}).Error("invalid action")
		return false
	}
	return false
}
func (o *Outlet) GetName() string {
	return o.Name
}
func (o *Outlet) SetValue(value float64) bool {
	oid := fmt.Sprintf(".1.3.6.1.4.1.318.1.1.12.3.3.1.1.4.16.%d", o.Index)
	pdu := gosnmp.SnmpPDU{Name: oid, Value: int(value), Type: gosnmp.Integer}
	result, err := o.PDU.SNMP.Set([]gosnmp.SnmpPDU{pdu})
	log.WithFields(log.Fields{"result": result, "error": err, "value": value}).Error("o.PDU.SNMP.Set")
	if result.Error > 0 {
		return false
	} else {
		return true
	}
}
func (o *Outlet) GetLocationOne() string {
	return "Basement"
}
func (o *Outlet) GetLocationTwo() string {
	return o.PDU.Location
}
func (o *Outlet) ProviderIDString() string {
	return o.PDU.Name
}
func (o *Outlet) IDString() string {
	return fmt.Sprintf("%d", o.Index)
}
func (o *Outlet) Matches(find apimodels.Match) bool {
	var rVal *regexp.Regexp
	var sVal string
	var ok bool
	var err error
	sVal, ok = find["RegexpName"].(string)
	if ok {
		rVal, err = regexp.Compile(sVal)
		if err != nil {
			log.WithFields(log.Fields{"regexp": sVal, "error": err}).Error("unable to compile regexp")
			return false
		}
		ok = rVal.MatchString(o.Name)
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
	sVal, ok = find["LocationTwo"].(string)
	if ok {
		if sVal != o.PDU.Location {
			return false
		}
	}
	sVal, ok = find["Name"].(string)
	if ok {
		if sVal != o.Name {
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
	return true
}

func (p *PDU) SetDevicesValue(find apimodels.Match, value float64) bool {
	devices, ok := p.GetDevices(find)
	if ok {
		devices.SetValue(value)
		return true
	}
	return false
}

/*
	case "immediateOn":
	case "on":
	case "immediateOff":
	case "off":
	case "immediateReboot":
	case "reboot":
*/

type OutletAction struct {
	Name string
}

func (a *OutletAction) GetName() string {
	return a.Name
}

func (o *Outlet) ListActions() []apimodels.Action {
	return []apimodels.Action{
		&OutletAction{Name: "on"},
		&OutletAction{Name: "off"},
		&OutletAction{Name: "reboot"},
	}
}

func (p *PDU) InvokeDevicesAction(find apimodels.Match, action string) bool {
	devices, ok := p.GetDevices(find)
	if ok {
		devices.InvokeAction(action)
		return true
	}
	return false
}
func (p *PDU) SetChildDevicesValue(apimodels.Match, float64) bool {
	return false
}
func (p *PDU) InvokeChildDevicesAction(apimodels.Match, string) bool {
	return false
}

type OutletList struct {
	Devices []*Outlet
}

func (l *OutletList) InvokeAction(action string) bool {
	for _, device := range l.Devices {
		device.InvokeAction(action)
	}
	return true
}
func (l *OutletList) Add(device *Outlet) {
	l.Devices = append(l.Devices, device)
}

func (l *OutletList) SetValue(value float64) bool {
	for _, device := range l.Devices {
		device.SetValue(value)
	}
	return true
}
func (l *OutletList) List() []apimodels.Device {
	r := make([]apimodels.Device, 0)
	for _, c := range l.Devices {
		r = append(r, c)
	}
	return r
}

func (p *PDU) GetDevices(find apimodels.Match) (apimodels.Devices, bool) {
	result := &OutletList{}
	for _, device := range p.Outlets {
		if device.Matches(find) {
			result.Add(device)
		}
	}
	return result, true
}
func (p *PDU) GetDevice(find apimodels.Match) (apimodels.Device, bool) {
	for _, device := range p.Outlets {
		if device.Matches(find) {
			return device, true
		}
	}
	return &Outlet{}, false
}

func (p *PDU) GetChildDevice(find apimodels.Match) (apimodels.Device, bool) {
	return &Outlet{}, false
}

func (p *PDU) GetChildDevices(find apimodels.Match) (apimodels.Devices, bool) {
	return &OutletList{}, false
}

type PDUConfiguration struct {
	IP        string `json:"IP"`
	Community string `json:"Community"`
}

func (p *PDU) Create(configurationRaw json.RawMessage) bool {
	var configuration *PDUConfiguration
	configuration = &PDUConfiguration{}
	json.Unmarshal(configurationRaw, configuration)

	if configuration.IP == "" {
		log.WithFields(log.Fields{}).Error("IP is a required configuration option")
		return false
	} else {
		p.IP = net.ParseIP(configuration.IP)
	}

	if configuration.Community == "" {
		log.WithFields(log.Fields{}).Error("Community is a required configuration option")
		return false
	} else {
		p.Community = configuration.Community
	}
	p.Init()
	return true
}

func (p *PDU) Print() {
	for _, outlet := range p.Outlets {
		outlet.Print()
	}
}

func (p *PDU) Load() {
	var err error
	sl := &SNMPTable{OneStart: 16, OneEnd: 17, TwoStart: 15, TwoEnd: 16}

	err = p.SNMP.BulkWalk(".1.3.6.1.4.1.318.1.1.12.3.5.1.1", sl.Callback)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetPDU p.SNMP.BulkWalk failed")
		return
	}
	var outlets []*Outlet
	var ok bool
	var pdu gosnmp.SnmpPDU
	var name string
	var phase int64
	var index int64
	var state bool
	var state_i int64
	for k, l := range sl.Values {
		// The name of the outlet. The maximum string size is device dependent. This OID is provided for informational purposes only.
		pdu, ok = l["1"]
		if !ok {
			log.WithFields(log.Fields{"k": k, "l": l}).Error("rPDUOutletStatusIndex missing")
			continue
		}

		index = pdu.BigInt().Int64()
		// The name of the outlet. The maximum string size is device dependent. This OID is provided for informational purposes only.
		pdu, ok = l["2"]
		if !ok {
			log.WithFields(log.Fields{"k": k, "l": l}).Error("rPDUOutletStatusOutletName missing")
			continue
		}
		name = pdu.String()

		// The phase/s associated with this outlet.
		// For single phase devices, this object will always return phase1(1).
		// For 3-phase devices, this object will return phase1 (1), phase2 (2), or phase3 (3) for outlets tied to a single phase.
		// For outlets tied to two phases, this object will return phase1-2 (4) for phases 1 and 2, phase2-3 (5) for phases 2 and 3, and phase3-1 (6) for phases 3 and 1.
		pdu, ok = l["3"]
		if !ok {
			log.WithFields(log.Fields{"k": k, "l": l}).Error("rPDUOutletStatusOutletPhase missing")
			continue
		}
		phase = pdu.BigInt().Int64()

		pdu, ok = l["4"]
		if !ok {
			log.WithFields(log.Fields{"k": k, "l": l}).Error("rPDUOutletStatusOutletState missing")
			continue
		}
		// Getting this variable will return the outlet state.
		// If the outlet is on, the outletStatusOn (1) value will be returned.
		// If the outlet is off, the outletStatusOff (2) value will be returned.

		switch state_i = pdu.BigInt().Int64(); state_i {
		case 1:
			state = true
		case 2:
			state = false
		default:
			log.WithFields(log.Fields{"state": state_i, "l": l}).Error("rPDUOutletStatusOutletState	invalid")
			continue
		}

		outlets = append(outlets, &Outlet{
			Index: index,
			Name:  name,
			Phase: phase,
			State: state,
			PDU:   p,
		})
	}
	p.Outlets = outlets
}
func (p *PDU) Init() {
	p.SNMP = &(*gosnmp.Default)
	p.SNMP.Target = p.IP.String()
	p.SNMP.Version = gosnmp.Version2c
	p.SNMP.Timeout = time.Duration(10 * time.Second)
	p.SNMP.Community = p.Community
	err := p.SNMP.Connect()
	if err != nil {
		log.WithFields(log.Fields{"ip": p.IP.String(), "community": p.Community, "error": err}).Error("unable to connect to SNMP target")
		return
	}

	result, err := p.SNMP.Get([]string{".1.3.6.1.2.1.1.5.0", ".1.3.6.1.2.1.1.6.0"})
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("unable to request")
		return
	}
	name, err := result.GetPDU(".1.3.6.1.2.1.1.5.0")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetPDU sysName failed")
		return
	}
	p.Name = name.String()

	location, err := result.GetPDU(".1.3.6.1.2.1.1.6.0")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("GetPDU sysLocation failed")
		return
	}
	p.Location = location.String()
	p.Load()

}
func (p *PDU) TypeString() string {
	return "APC PDU"
}
func (p *PDU) IDString() string {
	return p.Name
}

func NewAPCPDU(IP net.IP, community string) *PDU {
	apci := &PDU{
		IP:        IP,
		Community: community,
	}
	apci.Init()
	return apci
}
