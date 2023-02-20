package main

import (
	"html/template"
	"log"
	"net/http"
)

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", index)

	log.Println("Servidor corriendo")

	http.ListenAndServe(":8080", nil)

}

func index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hola developer")
	plantillas.ExecuteTemplate(w, "index", nil)
}
