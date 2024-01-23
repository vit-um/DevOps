package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func ml5(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Print(fmt.Sprintf("GET: %s", r.URL.Path))
		if r.URL.RawQuery == "" {
			r.URL.Path = "index.html"
		} else if r.URL.RawQuery != "" {
			log.Print(fmt.Sprintf("Q: %s", r.URL.RawQuery))
			r.URL.Path = "ml5.html"
		}
		lp := filepath.Join("ml5/templates", "layout.html")
		fp := filepath.Join("ml5/templates", filepath.Clean(r.URL.Path))

		// Return a 404 if the template doesn't exist
		info, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				log.Print(fmt.Sprintf("no template %s", err))
				http.NotFound(w, r)
				return
			}
		}

		// Return a 404 if the request is for a directory
		if info.IsDir() {
			log.Print("is dir")
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles(lp, fp)
		tmpl.New("img").Parse(`{{define "img"}}` + `static/img/img.` + r.URL.RawQuery + `{{end}}`)

		if err != nil {
			// Log the detailed error
			log.Println(err.Error())
			// Return a generic "Internal Server Error" message
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}

	case "POST":
		log.Print(fmt.Sprintf("POST: %s", r.URL.Path))

		tmpfile, err := ioutil.TempFile("ml5/img", "img.")
		if err != nil {
			log.Print(err)
		}
		f, _, _ := r.FormFile("image")

		defer f.Close() // clean up

		io.Copy(tmpfile, f)
		log.Print(tmpfile.Name())
		w.Write([]byte(fmt.Sprintf(`{"uploadUrl":"?%s"}`, strings.Split(tmpfile.Name(), ".")[1])))
	}
}
