package kinopoisk

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"parser/utility"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

const timeout = 20
const mainIrl = "https://www.kinopoisk.ru/"

type Kinopoisk struct {
	client *http.Client
	token  string
}

func New() *Kinopoisk {
	jar, err := cookiejar.New(nil)
	if err != nil {
		logrus.Fatalf("got error while creating cookie jar, %v", err)
	}

	var kp Kinopoisk
	kp.client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * time.Duration(timeout),
	}
	return &kp
}

// returns a list of releases filtred by the given month, coutry, genre
func (kp *Kinopoisk) Releases(countries, genres []string, month int) Films {
	films := make(Films, 0)
	p := kp.releasesFirstPageAndToken(month)
	films = append(films, kp.parseReleasesPage(p)...)

	//приходит только первая страница, нужен дебаг в фиддлере
	p = kp.releasesNextPage(month, 1)
	films = append(films, kp.parseReleasesPage(p)...)

	if countries == nil && genres == nil {
		return films
	}

	res := make(Films, 0)

	for _, film := range films {

		var matchByCountry, matchByGenre = true, true
		if countries != nil {
			matchByCountry = func() bool {
				for _, c := range countries {
					if strings.Contains(film.Country, c) {
						return true
					}
				}
				return false
			}()
		}

		if genres != nil {
			matchByGenre = func() bool {
				for _, g := range genres {
					if strings.Contains(film.Genre, g) {
						return true
					}
				}
				return false
			}()
		}

		if matchByCountry && matchByGenre {
			res = append(res, film)
		}
	}

	return res
}

func (kp *Kinopoisk) releasesFirstPageAndToken(month int) []byte {
	url := fmt.Sprintf("https://www.kinopoisk.ru/premiere/ru/2022/month/%d/", month)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Fatalf("cannot create NewRequest to %s, %v", url, err)
	}

	resp, err := kp.client.Do(req)
	if err != nil {
		logrus.Fatalf("cannot get %s, %v", url, err)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "_csrf_csrf_token" {
			kp.token = cookie.Value
		}
	}
	return utility.BytesFromReader(resp.Body)
}

func (kp *Kinopoisk) releasesNextPage(month, page int) []byte {
	url := fmt.Sprintf("https://www.kinopoisk.ru/premiere/ru/2022/month/%d/", month)
	body := fmt.Sprintf("token=%s&page=%d&ajax=true", kp.token, page)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	if err != nil {
		logrus.Fatalf("cannot create NewRequest to %s, %v", url, err)
	}

	resp, err := kp.client.Do(req)
	if err != nil {
		logrus.Fatalf("cannot get %s, %v", url, err)
	}
	return utility.BytesFromReader(resp.Body)
}

func (kp *Kinopoisk) parseReleasesPage(page []byte) Films {
	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// path := filepath.Join(filepath.Dir(wd), "parser/resources/Fiddler_19-57-29.htm")
	// f, _ := os.Open(path)

	// doc, err := goquery.NewDocumentFromReader(f)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}

	sel := doc.Find(".prem_list")
	items := sel.Eq(1).Find(".premier_item")

	var res Films

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
