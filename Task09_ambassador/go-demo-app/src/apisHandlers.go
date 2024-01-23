package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

func ascii(w http.ResponseWriter, r *http.Request) {

	var b []byte
	b = append([]byte(""), Environment...)

	w.Write(b)

}

func data(w http.ResponseWriter, r *http.Request) {

	var b []byte
	b = append([]byte(""), Environment...)

	w.Write(b)

}

func version(ctx *fasthttp.RequestCtx) {
	var b []byte
	b = append([]byte(""), Environment...)

	ctx.Write(b)
}

func healthz(ctx *fasthttp.RequestCtx) {

	ctx.Write([]byte("Healthz: alive!"))
}

func readinez(w http.ResponseWriter, r *http.Request) {

	flag.Parse()
	switch Role {

	case "api":
		w.Write([]byte("READY"))

	case "img", "ml5":
		w.Write([]byte("READY"))

	case "ascii":
		_, err := CACHE.Ping().Result()
		if err != nil {
			log.Print(err)
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		} else {

			w.Write([]byte("READY"))
		}

	case "data":
		_, err := CACHE.Set("readiness_probe", 0, 0).Result()
		if err != nil {
			log.Print(err)
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		}

		err = DB.Ping()

		if err != nil {
			log.Print(err)
			http.Error(w, "Not Ready", http.StatusServiceUnavailable)
		} else {

			w.Write([]byte("READY"))

		}

	default:
		http.Error(w, "Not Ready", http.StatusServiceUnavailable)

	}

}
