package pkg

import (
	"io/ioutil"
	"net/http"
	"text/template"
)

var html, _ = ioutil.ReadFile("./frontend.html")
var frontend = string(html)
var homeTempl = template.Must(template.New("").Parse(frontend))

func ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v = struct {
		Host string
	}{
		r.Host,
	}
	homeTempl.Execute(w, &v)
}
