package main

import (
	"fmt"
	"net/http"

	"github.com/anaskhan96/soup"
	"github.com/gorilla/mux"
)

func ErrCheck(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {
	h := mux.NewRouter()

	resp, err := soup.Get("https://cs.sogang.ac.kr/front/cmsboardlist.do?siteId=cs&bbsConfigFK=1905")
	ErrCheck(err)
	doc := soup.HTMLParse(resp)
	lis := doc.Find("div", "class", "list_box").FindAll("li")
	for _, li := range lis {
		fmt.Println(li.Find("div").Find("a").Text())
	}

	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to Main Page")
	})
	http.Handle("/", h)
	http.ListenAndServe(":4567", nil)

}
