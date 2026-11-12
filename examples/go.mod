module github.com/tochk/twf/examples

go 1.27.0

replace github.com/tochk/twf => ../

tool github.com/valyala/quicktemplate/qtc

require (
	github.com/gorilla/mux v1.8.1
	github.com/tochk/twf v1.1.2
	github.com/valyala/quicktemplate v1.8.0
)

require github.com/valyala/bytebufferpool v1.0.0 // indirect
