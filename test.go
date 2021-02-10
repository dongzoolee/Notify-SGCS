package main

import (
	"fmt"

	"github.com/anaskhan96/soup"
)

type BoardItems struct {
	Title string
	Url   string
}

func main() {
	resp, err := soup.Get("https://cs.sogang.ac.kr/front/cmsboardlist.do?siteId=cs&bbsConfigFK=1747")
	if err != nil {
		panic(err)
	}
	doc := soup.HTMLParse(resp)
	lis := doc.Find("div", "class", "list_box").FindAll("li")

	for _, li := range lis {
		strong := li.Find("div").Find("a").Find("strong")
		title := li.Find("div").Find("a").Text()
		fmt.Println(strong.Error == nil)
		fmt.Println(title)
	}
}
