// This file is automatically generated by qtc from "row_thermostat.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line api/server/templates/row_thermostat.qtpl:1
package templates

//line api/server/templates/row_thermostat.qtpl:1
import "github.com/AdamJacobMuller/home-api/api/models"

//line api/server/templates/row_thermostat.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line api/server/templates/row_thermostat.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line api/server/templates/row_thermostat.qtpl:3
type Thermostat struct {
	Title      string
	ProviderID string
	DeviceID   string
	Device     *apimodels.Device
	Actions    []*Action
}

func (b *Thermostat) AddAction(action *Action) {
	b.Actions = append(b.Actions, action)
}

//line api/server/templates/row_thermostat.qtpl:15
func (b Thermostat) StreamHTML(qw422016 *qt422016.Writer) {
	//line api/server/templates/row_thermostat.qtpl:15
	qw422016.N().S(`
                    <div class="ibox-title">
                        <h5>`)
	//line api/server/templates/row_thermostat.qtpl:17
	qw422016.E().S(b.Title)
	//line api/server/templates/row_thermostat.qtpl:17
	qw422016.N().S(`</h5>
                        <div class="ibox-tools">
                            <a class="collapse-link">
                                <i class="fa fa-chevron-up"></i>
                            </a>
                        </div>
                    </div>
                    <div class="ibox-content">
                        <p>
                            `)
	//line api/server/templates/row_thermostat.qtpl:26
	for _, action := range b.Actions {
		//line api/server/templates/row_thermostat.qtpl:26
		qw422016.N().S(`
                                <button type="button" class="btn btn-w-m btn-default" `)
		//line api/server/templates/row_thermostat.qtpl:27
		StreamOnClickInvokeDeviceAction(qw422016, b.ProviderID, b.DeviceID, action.Title)
		//line api/server/templates/row_thermostat.qtpl:27
		qw422016.N().S(`>`)
		//line api/server/templates/row_thermostat.qtpl:27
		qw422016.E().S(action.Title)
		//line api/server/templates/row_thermostat.qtpl:27
		qw422016.N().S(`</button>
                            `)
		//line api/server/templates/row_thermostat.qtpl:28
	}
	//line api/server/templates/row_thermostat.qtpl:28
	qw422016.N().S(`
                        </p>
                    </div>
`)
//line api/server/templates/row_thermostat.qtpl:31
}

//line api/server/templates/row_thermostat.qtpl:31
func (b Thermostat) WriteHTML(qq422016 qtio422016.Writer) {
	//line api/server/templates/row_thermostat.qtpl:31
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line api/server/templates/row_thermostat.qtpl:31
	b.StreamHTML(qw422016)
	//line api/server/templates/row_thermostat.qtpl:31
	qt422016.ReleaseWriter(qw422016)
//line api/server/templates/row_thermostat.qtpl:31
}

//line api/server/templates/row_thermostat.qtpl:31
func (b Thermostat) HTML() string {
	//line api/server/templates/row_thermostat.qtpl:31
	qb422016 := qt422016.AcquireByteBuffer()
	//line api/server/templates/row_thermostat.qtpl:31
	b.WriteHTML(qb422016)
	//line api/server/templates/row_thermostat.qtpl:31
	qs422016 := string(qb422016.B)
	//line api/server/templates/row_thermostat.qtpl:31
	qt422016.ReleaseByteBuffer(qb422016)
	//line api/server/templates/row_thermostat.qtpl:31
	return qs422016
//line api/server/templates/row_thermostat.qtpl:31
}
