package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/0x3333/iptu.go/api"
	"github.com/0x3333/iptu.go/log"
	"github.com/0x3333/iptu.go/render"
	"github.com/ajays20078/go-http-logger"
)

// StartServer starts the webserver to handle the requests from the UI
func StartServer() {
	handlePesquisa()

	log.Info.Println("WebServer started...")
	http.ListenAndServe(":8080", httpLogger.WriteLog(http.DefaultServeMux, os.Stdout))
}

func handlePesquisa() {
	http.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.EscapedPath()[3:] // Remove o /s/ da URL
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
}
