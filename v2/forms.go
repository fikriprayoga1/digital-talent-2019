// forms.go
package main

import (
    "html/template"
	"net/http"
	"log"
)

type inputData struct {
    Temprature   string
}

func main() {
    tmpl := template.Must(template.ParseFiles("forms.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            tmpl.Execute(w, nil)
            return
        }

        details := inputData{
            Temprature:   r.FormValue("temperature"),
        }

		// do something with details
		log.Println(details.Temeprature)
        _ = details.Temeprature

        tmpl.Execute(w, struct{ Success bool }{true})
    })

    http.ListenAndServe(":8080", nil)
}