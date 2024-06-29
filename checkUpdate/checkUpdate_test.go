package checkupdate

import (
	"testing"
	"time"
)

func TestFindUpdatedArticle(t *testing.T) {
	parsed_time, err := time.Parse(time.RFC3339, "2024-06-11T05:04:11Z")
	if err != nil {
		t.Error("time parsing error", err)
	}

	// boards := FindUpdatedArticle("https://cs.sogang.ac.kr/front/cmsboardlist.do?siteId=cs&bbsConfigFK=1905", true, parsed_time)
	url := "https://cs.sogang.ac.kr/front/cmsboardlist.do?siteId=cs&bbsConfigFK=1746"
	items := FindUpdatedArticle(url, true, parsed_time)

	for _, val := range items {
		t.Log(val.Title, val.Url)
	}

	if len(items) == 0 {
		t.Fail()
	}
}
