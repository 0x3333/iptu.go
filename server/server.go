package server

import (
	"net/http"
	"strings"

	"bitbucket.org/terciofilho/iptu.go/api"
	"bitbucket.org/terciofilho/iptu.go/log"
	"bitbucket.org/terciofilho/iptu.go/render"
)

// StartServer starts the webserver to handle the requests from the UI
func StartServer() {
	// handleStatic()
	// handleAPI()
	handlePesquisa()

	log.Info.Println("WebServer started...")
	http.ListenAndServe(":8080", nil)
}

// func handleStatic() {
// 	http.Handle("/", http.FileServer(http.Dir("web")))
// }

// func handleAPI() {
// 	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
// 		if err := r.ParseForm(); err != nil {
// 			log.Error.Printf("Error parsing: %s", err.Error())
// 		}
//
// 		IPTUs, err := api.HandleRequest(r.FormValue("termos"))
// 		if err != nil {
//
// 		}
// 		if IPTUs == nil {
// 			w.Write([]byte("[]"))
// 		} else {
// 			bytes, err := json.Marshal(&IPTUs)
// 			if err != nil {
// 				log.Error.Println(err.Error())
// 			}
// 			w.Write(bytes)
// 		}
// 	})
// }

func handlePesquisa() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.EscapedPath()[1:]
		if strings.Contains(url, "/") {
			http.Error(w, "", http.StatusBadRequest)
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
