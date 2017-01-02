package fe

import (
	"fmt"
	// "html/template"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("yo homie")
	fmt.Fprint(w, "<h1>Homepage</h1><div>o boyo</div>")
}
