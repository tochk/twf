// Code generated by qtc from "custom_menu.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line custom_menu.qtpl:1
package example_templates

//line custom_menu.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line custom_menu.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line custom_menu.qtpl:1
func StreamAdminMenu(qw422016 *qt422016.Writer) {
//line custom_menu.qtpl:1
	qw422016.N().S(`
<nav class="navbar navbar-expand-lg navbar-light" style="background-color: #e3f2fd;">
    <a class="navbar-brand" href="/">TWF admin menu</a>
    <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent"
            aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
    </button>

    <div class="collapse navbar-collapse" id="navbarSupportedContent">
        <ul class="navbar-nav mr-auto">
            <li class="nav-item">
                <a class="nav-link" href="#">Example 1</a>
            </li>
            <li class="nav-item dropdown">
                <a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                    Dropdown example 1
                </a>
                <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                    <a class="dropdown-item" href="/users/">Users</a>
                    <div class="dropdown-divider"></div>
                    <a class="dropdown-item" href="#">Example 2</a>
                </div>
            </li>
        </ul>
        <ul class="navbar-nav justify-content-end">
            <li class="nav-item">
                <a class="nav-link active" href="/">Logout</a>
            </li>
        </ul>
    </div>
</nav>
`)
//line custom_menu.qtpl:32
}

//line custom_menu.qtpl:32
func WriteAdminMenu(qq422016 qtio422016.Writer) {
//line custom_menu.qtpl:32
	qw422016 := qt422016.AcquireWriter(qq422016)
//line custom_menu.qtpl:32
	StreamAdminMenu(qw422016)
//line custom_menu.qtpl:32
	qt422016.ReleaseWriter(qw422016)
//line custom_menu.qtpl:32
}

//line custom_menu.qtpl:32
func AdminMenu() string {
//line custom_menu.qtpl:32
	qb422016 := qt422016.AcquireByteBuffer()
//line custom_menu.qtpl:32
	WriteAdminMenu(qb422016)
//line custom_menu.qtpl:32
	qs422016 := string(qb422016.B)
//line custom_menu.qtpl:32
	qt422016.ReleaseByteBuffer(qb422016)
//line custom_menu.qtpl:32
	return qs422016
//line custom_menu.qtpl:32
}
