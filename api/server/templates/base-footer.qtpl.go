// This file is automatically generated by qtc from "base-footer.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line api/server/templates/base-footer.qtpl:1
package templates

//line api/server/templates/base-footer.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line api/server/templates/base-footer.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line api/server/templates/base-footer.qtpl:1
func (p *BasePage) StreamFooter(qw422016 *qt422016.Writer) {
	//line api/server/templates/base-footer.qtpl:1
	qw422016.N().S(`
                `)
	//line api/server/templates/base-footer.qtpl:2
	if false {
		//line api/server/templates/base-footer.qtpl:2
		qw422016.N().S(`
                <div class="footer">
                    <div class="pull-right">
                        10GB of <strong>250GB</strong> Free.
                    </div>
                    <div>
                        <strong>Copyright</strong> Example Company © 2014-2017
                    </div>
                </div>
                `)
		//line api/server/templates/base-footer.qtpl:11
	}
	//line api/server/templates/base-footer.qtpl:11
	qw422016.N().S(`
`)
//line api/server/templates/base-footer.qtpl:12
}

//line api/server/templates/base-footer.qtpl:12
func (p *BasePage) WriteFooter(qq422016 qtio422016.Writer) {
	//line api/server/templates/base-footer.qtpl:12
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line api/server/templates/base-footer.qtpl:12
	p.StreamFooter(qw422016)
	//line api/server/templates/base-footer.qtpl:12
	qt422016.ReleaseWriter(qw422016)
//line api/server/templates/base-footer.qtpl:12
}

//line api/server/templates/base-footer.qtpl:12
func (p *BasePage) Footer() string {
	//line api/server/templates/base-footer.qtpl:12
	qb422016 := qt422016.AcquireByteBuffer()
	//line api/server/templates/base-footer.qtpl:12
	p.WriteFooter(qb422016)
	//line api/server/templates/base-footer.qtpl:12
	qs422016 := string(qb422016.B)
	//line api/server/templates/base-footer.qtpl:12
	qt422016.ReleaseByteBuffer(qb422016)
	//line api/server/templates/base-footer.qtpl:12
	return qs422016
//line api/server/templates/base-footer.qtpl:12
}
