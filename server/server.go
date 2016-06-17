package server

import (
	"net/http"
	"os"
	"strings"

	"bitbucket.org/terciofilho/iptu.go/api"
	"bitbucket.org/terciofilho/iptu.go/log"
	"bitbucket.org/terciofilho/iptu.go/render"
	"github.com/ajays20078/go-http-logger"
	"github.com/nytimes/gziphandler"
)

// StartServer starts the webserver to handle the requests from the UI
func StartServer() {
	handleStatic()
	handlePesquisa()

	log.Info.Println("WebServer started...")
	http.ListenAndServe(":8080", httpLogger.WriteLog(http.DefaultServeMux, os.Stdout))
}

func handleStatic() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == "/" {
			http.Redirect(w, r, "/s/", http.StatusMovedPermanently)
		} else {
			http.ServeFile(w, r, "web/"+r.URL.Path[1:])
		}
	})
}

func handlePesquisa() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redireciona para o domínio sem www(Qualquer dominio diferente de 'consultaiptu.com.br')
		domainParts := strings.Split(r.Host, ".")
		if len(domainParts) != 3 || domainParts[0] != "consultaiptu" {
			http.Redirect(w, r, "http://consultaiptu.com.br"+r.URL.EscapedPath(), http.StatusMovedPermanently)
			return
		}
		url := r.URL.EscapedPath()[3:]
		if strings.Contains(url, "/") {
			http.Error(w, "400 - Bad Request", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		termos := strings.Replace(url, "-", " ", -1)
		if len(termos) == 0 {
			render.Render(nil, true, false, false, w)
		} else {
			IPTUs, err := api.HandleRequest(termos)
			if err != nil {
				log.Error.Println(err.Message)
				render.Render(nil, false, err.Invalid, err.HasError, w)
			} else {
				render.Render(IPTUs, false, false, false, w)
			}
		}
	})
	handlerGz := gziphandler.GzipHandler(handler)
	http.Handle("/s/", handlerGz)
}
