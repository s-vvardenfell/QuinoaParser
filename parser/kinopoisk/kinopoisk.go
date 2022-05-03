package kinopoisk

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Kinopoisk struct {
}

// returns a list of releases by the given month
func (kp *Kinopoisk) ReleasesByMonth(month int) []Film {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(filepath.Dir(wd), "parser/resources/Fiddler_19-57-29.htm")
	// path := filepath.Join("../", "resources/Fiddler_19-57-29.htm")
	f, _ := os.Open(path)

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	sel := doc.Find(".prem_list")
	items := sel.Eq(1).Find(".premier_item")

	var res []Film

	for i := 0; i < items.Size(); i++ {
		rd, ok := items.Eq(i).Find(`[itemprop="startDate"]`).Attr("content")
		if !ok {
			rd = ""
		}

		mu, ok := items.Eq(i).Find(`[itemprop="image"]`).Attr("content")
		if !ok {
			mu = ""
		}

		u, ok := items.Eq(i).Find(".text").Find(".name").Find("a:nth-child(1)").Attr("href")
		if !ok {
			u = ""
		}

		cntr := clearCoutry(items.Eq(i).Find(`.text > div > span:nth-child(3)`).Text())
		genre := clearGenre(items.Eq(i).Find(`.text > div > span:nth-child(4)`).Text())

		res = append(res, Film{
			Title:         items.Eq(i).Find(".text").Find(".name").Text(),
			Country:       cntr,
			ReleaseDate:   rd,
			MiniPosterUrl: mu,
			FilmUrl:       u,
			Genre:         genre,
			Description:   "",
		})
	}
	return res
}

func clearCoutry(in string) string {
	in = strings.TrimSpace(in)
	out := strings.SplitAfter(in, "...")
	out = strings.SplitAfter(out[0], ",")
	return strings.TrimRight(out[0], ",.")
}

func clearGenre(in string) string {
	in = strings.Trim(in, "().")
	return in
}
