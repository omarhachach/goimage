package main

import (
	"html/template"

	"github.com/omar-h/goimage"
	"github.com/valyala/fasthttp"
)

// indexHandler serves the index at the index route.
func indexHandler(template *template.Template) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		template.Execute(ctx, nil)
	}
}

func viewHandler(template *template.Template) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		id := ctx.UserValue("id")
		exists, fileInfo, err := goimage.GetFileInfo(config.ImageDirectory, id.(string))
		if err != nil {
			ctx.SetStatusCode(500)
			return
		}

		if !exists {
			ctx.SetStatusCode(400)
			return
		}

		template.Execute(ctx, map[string]interface{}{
			"id":  id,
			"ext": fileInfo.Extension,
		})
	}
}
