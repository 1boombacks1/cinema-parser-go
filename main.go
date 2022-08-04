package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
)

func collectData(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find(".site_content .movieItem").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Find("a").Attr("href"); ok {
			g.Get(r.JoinURL(href), func(_g *geziyor.Geziyor, _r *client.Response) {

				genres := make([]string, 0)

				_r.HTMLDoc.Find(".filmInfo_genre a").Each(func(i int, _s *goquery.Selection) {
					genres = append(genres, _s.Text())
				})

				genre := strings.Join(genres, ", ")

				desc := _r.HTMLDoc.Find(".tabs_contentItem:nth-child(3) p").Text()
				duration := _r.HTMLDoc.Find(".filmInfo_info span:nth-child(2)").First().Text()
				allSessions := make(map[string][]string)

				_r.HTMLDoc.Find(".schedule_showtimes > .showtimes_item").Each(func(i int, _s *goquery.Selection) {
					cinema := _s.Find(".showtimesCinema_name").Text()
					sessions := make([]string, 0)

					_s.Find(".showtimes_sessions > a > span:nth-child(1)").Each(func(i int, sl *goquery.Selection) {
						sessions = append(sessions, sl.Text())
					})

					allSessions[cinema] = sessions
				})

				g.Exports <- map[string]interface{}{
					"Название":          strings.TrimSpace(s.Find(".movieItem_title").Text()),
					"Жанр":              strings.TrimSpace(genre),
					"Описание":          strings.TrimSpace(strings.ReplaceAll(desc, "\t", "")),
					"Продолжительность": strings.TrimSpace(duration),
					"Сеансы":            allSessions,
				}
			})
		}
	})
}

func main() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://yakutsk.kinoafisha.info/movies/?date=2022-07-14"},
		ParseFunc: collectData,
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()
}
