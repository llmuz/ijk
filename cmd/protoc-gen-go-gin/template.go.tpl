// generate by protoc-gen-go-gin
// powered by xx

type {{ $.InterfaceName }} interface {
{{range .MethodSet}}
    // {{.Path}} {{.Method}}
	{{.Name}}(ctx context.Context, req *{{.Request}}) (resp *{{.Reply}}, err error)
{{end}}
}


func New{{$.Name}}Server(srv IJKGinServer, srvHandler ginsrv.ServiceHandler) (c *{{$.Name}}) {
	c = &{{$.Name}}{
		server:     srv,
		srvHandler: srvHandler,
	}
	return c
}

type {{$.Name}} struct{
	server {{ $.InterfaceName }}
	srvHandler ginsrv.ServiceHandler
}



{{range .Methods}}
func (s *{{$.Name}}) {{ .HandlerName }} (ctx *gin.Context) {
	var in {{.Request}}
{{if .HasPathParams }}
	if err := ctx.ShouldBindUri(&in); err != nil {
		s.srvHandler.ParamsError(ctx, err)
		return
	}
{{end}}
{{if eq .Method "GET" "DELETE" }}
	if err := ctx.ShouldBindQuery(&in); err != nil {
		s.srvHandler.ParamsError(ctx, err)
		return
	}
{{else if eq .Method "POST" "PUT" }}
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

    
	if err := s.srvHandler.Validate(ctx, &in); err != nil {
		s.srvHandler.ParamsError(ctx, err)
		return
	}

	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "request_context", ctx))

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.({{ $.InterfaceName }}).{{.Name}}(newCtx, &in)
	if err != nil {
		s.srvHandler.Error(ctx, err)
		return
	}

	s.srvHandler.Success(ctx, out)
}
{{end}}
