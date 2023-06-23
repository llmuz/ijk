// generate by protoc-gen-go-gin
// powered by xx

{{$.Comment}}
type {{ $.InterfaceName }} interface {
{{range .MethodSet}}
    // {{.Path}} {{.Method}}
    {{.Comment}}
	{{.Name}}(ctx context.Context, req *{{.Request}}) (resp *{{.Reply}}, err error)
{{end}}
}




// {{ $.InterfaceName }} service
type {{$.ImplHandlerName}} struct{
	server {{ $.InterfaceName }}
	srvHandler ginsrv.ServiceHandler
}

// New{{$.ImplHandlerName}} build {{$.ImplHandlerName}}
func New{{$.ImplHandlerName}}(srv {{ $.InterfaceName }}, srvHandler ginsrv.ServiceHandler) (c *{{$.Name}}Handler) {
	c = &{{$.Name}}Handler{
		server:     srv,
		srvHandler: srvHandler,
	}
	return c
}

{{range .Methods}}
 {{.Comment}}
func (s *{{$.ImplHandlerName}}) {{ .HandlerName }} (ctx *gin.Context) {
	var in {{.Request}}
{{if .HasPathParams }}
    // bind uri param, etc: /api/:name bind :name
    // type FooBar struct {
    //      Name string `uri:"name"` // binding uri name
    // }
	if err := ctx.ShouldBindUri(&in); err != nil {
		s.srvHandler.ParamsError(ctx, err)
		return
	}
{{end}}
{{if eq .Method "GET" "DELETE" }}
    // bind query param, etc: /api/name?query=data, bind query
    // type FooBar struct {
    //      Name string `form:"name"` // binding form name
    // }
	if err := ctx.ShouldBindQuery(&in); err != nil {
		s.srvHandler.ParamsError(ctx, err)
		return
	}
{{else if eq .Method "POST" "PUT" }}
    // bind body data, etc: PUT, POST
    // type FooBar struct {
    //      Name string `json:"name"  // binding json name
    // }
	if err := ctx.ShouldBindJSON(&in); err != nil {
		s.srvHandler.ParamsError(ctx, err)
		return
	}
{{else}}
	if err := ctx.ShouldBind(&in); err != nil {
		s.srvHandler.ParamsError(ctx, err)
		return
	}
{{end}}

    // validate param ok or not?
    // type FooBar struct {
    //      Name string `json:"name" validate:"min=10,max=32"` // 表示 name的长度只能在 [10, 32]
    // }
	if err := s.srvHandler.Validate(ctx, &in); err != nil {
		s.srvHandler.Error(ctx, err)
		return
	}

    // pull *gin.Context
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "request_context", ctx))

    // pull http request header
	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx.Request.Context(), md)

	// call HandlerImpl
	out, err := s.server.({{ $.InterfaceName }}).{{.Name}}(newCtx, &in)
	// handler error
	if err != nil {
		s.srvHandler.Error(ctx, err)
		return
	}

    // handler success
	s.srvHandler.Success(ctx, out)
}
{{end}}
