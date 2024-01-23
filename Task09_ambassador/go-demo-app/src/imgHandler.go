package main

import (
	"bytes"
	_ "image/jpeg"
	"io"
	"log"

	"github.com/den-vasyliev/image2ascii/convert"
	"github.com/valyala/fasthttp"
)

func img(ctx *fasthttp.RequestCtx) {

	switch string(ctx.Method()) {
	case "GET":
		var b []byte
		b = append([]byte(""), Environment...)
		ctx.Write(b)
	case "POST":
		var Buf bytes.Buffer

		f, _ := ctx.FormFile("image")
		ff, _ := f.Open()
		//defer f.Close()
		io.Copy(&Buf, ff)
		b := Buf.Bytes()

		Buf.Reset()
		//ctx.Header().Set("Content-Type", "text/plain")

		if convertOptions, err := parseOptions(); err == nil {
			converter := convert.NewImageConverter()

			ctx.Write([]byte(converter.ImageBuf2ASCIIString(b, convertOptions)))
		} else {
			log.Print("No opt")
		}

	}
}
