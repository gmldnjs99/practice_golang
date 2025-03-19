package collector

import (
	"fmt"
	"log"
	"github.com/gocolly/colly"
)

type NewsData struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	URL     string `json:"url"`
	Date    string `json:"date"`
}

func GetNewsData(query string) ([]NewsData, error) {
	var newsList []NewsData

	// 검색할 URL 설정 (여기서는 Naver 뉴스 검색을 예시로 사용)
	searchURL := fmt.Sprintf("https://search.naver.com/search.naver?where=news&query=%s", query)

	c := colly.NewCollector()

	// 뉴스 항목 처리
	c.OnHTML("ul.list_news div.news_area", func(e *colly.HTMLElement) {
		title := e.ChildText("a.news_tit")
		url := e.ChildAttr("a.news_tit", "href")
		summary := e.ChildText("div.dsc_wrap")

		// 뉴스 항목 저장
		news := NewsData{
			Title:   title,
			Summary: summary,
			URL:     url,
			Date:    "", // 날짜 처리 필요 시 추가 가능
		}

		// 리스트에 추가
		newsList = append(newsList, news)
	})

	// 요청 처리
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// URL 방문하여 데이터 크롤링
	err := c.Visit(searchURL)
	if err != nil {
		log.Println("❌ 크롤링 실패:", err)
		return nil, err
	}

	// 최대 10개의 뉴스만 반환 (옵션)
	if len(newsList) > 10 {
		newsList = newsList[:10]
	}

	return newsList, nil
}
