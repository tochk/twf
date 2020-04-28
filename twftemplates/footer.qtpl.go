// Code generated by qtc from "footer.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line footer.qtpl:1
package twftemplates

//line footer.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line footer.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line footer.qtpl:1
func StreamFooter(qw422016 *qt422016.Writer) {
//line footer.qtpl:1
	qw422016.N().S(`
<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js"
        integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN"
        crossorigin="anonymous"></script>
<script src="https://unpkg.com/popper.js@1.12.6/dist/umd/popper.js"
        integrity="sha384-fA23ZRQ3G/J53mElWqVJEGJzU0sTs+SvzG8fXVWP+kJQ1lwFAOkcUOysnlKJC33U"
        crossorigin="anonymous"></script>
<script src="https://unpkg.com/bootstrap-material-design@4.1.1/dist/js/bootstrap-material-design.js"
        integrity="sha384-CauSuKpEqAFajSpkdjv3z9t8E7RlpJ1UP0lKM/+NdtSarroVKu069AlsRPKkFBz9"
        crossorigin="anonymous"></script>
</body>
</html>
`)
//line footer.qtpl:13
}

//line footer.qtpl:13
func WriteFooter(qq422016 qtio422016.Writer) {
//line footer.qtpl:13
	qw422016 := qt422016.AcquireWriter(qq422016)
//line footer.qtpl:13
	StreamFooter(qw422016)
//line footer.qtpl:13
	qt422016.ReleaseWriter(qw422016)
//line footer.qtpl:13
}

//line footer.qtpl:13
func Footer() string {
//line footer.qtpl:13
	qb422016 := qt422016.AcquireByteBuffer()
//line footer.qtpl:13
	WriteFooter(qb422016)
//line footer.qtpl:13
	qs422016 := string(qb422016.B)
//line footer.qtpl:13
	qt422016.ReleaseByteBuffer(qb422016)
//line footer.qtpl:13
	return qs422016
//line footer.qtpl:13
}
