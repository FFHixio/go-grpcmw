package template

import "text/template"

// Code template keys
const (
	serviceKey     = "service"
	serviceTypeKey = "serviceType"
)

// Code templates
const (
	serviceTypeCode = `Interceptor_{{.Package}}{{.Service}}`

	serviceCode = `
type server{{template "serviceType" .}} struct {
	grpcmw.ServerInterceptor
}

type client{{template "serviceType" .}} struct {
	grpcmw.ClientInterceptor
}

func (i *server{{template "pkgType" .}}) Register{{.Service}}() *server{{template "serviceType" .}} {
	service, ok := i.ServerInterceptor.(grpcmw.ServerInterceptorRegister).Get("{{.Service}}")
	if !ok {
		ret := &server{{template "serviceType" .}}{
			ServerInterceptor: grpcmw.NewServerInterceptorRegister("{{.Service}}"),
		}
		i.ServerInterceptor.(grpcmw.ServerInterceptorRegister).Register(ret.ServerInterceptor)
		{{with .Interceptors}}ret.ServerInterceptor.Merge({{range .Indexes}}
			registry.GetServerInterceptor("{{.}}"),{{end}}
		){{end}}
		{{range .Methods}}{{if .Interceptors}}
		ret.{{.Method}}().AddInterceptor({{$method := .}}{{range .Interceptors.Indexes}}
			registry.GetServerInterceptor("{{.}}").{{template "methodType" $method.Stream}}ServerInterceptor(),{{end}}
		){{end}}{{end}}
		return ret
	}
	return &server{{template "serviceType" .}}{
		ServerInterceptor: service,
	}
}

func (i *client{{template "pkgType" .}}) Register{{.Service}}() *client{{template "serviceType" .}} {
	service, ok := i.ClientInterceptor.(grpcmw.ClientInterceptorRegister).Get("{{.Service}}")
	if !ok {
		ret := &client{{template "serviceType" .}}{
			ClientInterceptor: grpcmw.NewClientInterceptorRegister("{{.Service}}"),
		}
		i.ClientInterceptor.(grpcmw.ClientInterceptorRegister).Register(ret.ClientInterceptor)
		{{with .Interceptors}}ret.ClientInterceptor.Merge({{range .Indexes}}
			registry.GetClientInterceptor("{{.}}"),{{end}}
		){{end}}
		{{range .Methods}}{{if .Interceptors}}
		ret.{{.Method}}().AddInterceptor({{$method := .}}{{range .Interceptors.Indexes}}
			registry.GetClientInterceptor("{{.}}").{{template "methodType" $method.Stream}}ClientInterceptor(),{{end}}
		){{end}}{{end}}
		return ret
	}
	return &client{{template "serviceType" .}}{
		ClientInterceptor: service,
	}
}

{{range .Methods}}{{template "method" .}}{{end}}
`
)

func init() {
	template.Must(initCodeTpl.New(serviceKey).Parse(serviceCode))
	template.Must(initCodeTpl.New(serviceTypeKey).Parse(serviceTypeCode))
}
