package platform

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

const (
	mainUrl     = "www.kinopoisk.ru" // В КОНФИГ
	queryUrl    = "www.kinopoisk.ru/s/index.php"
	searchUrl   = "https://www.kinopoisk.ru/s/"
	imgUrlTempl = "https://www.kinopoisk.ru/images/sm_film/%s.jpg"
)

type Result struct {
	Name   string
	Ref    string
	ImgUrl string //нужно еще сдеать GET и получить картинку
}

//ВЫНЕСТИ
type Condition struct {
	Keyword  string
	Type     string
	Genres   []string
	YearFrom string
	YearTo   string
	Coutries []string
}

type Platform struct {
}

func New() *Platform {
	return &Platform{}
}

/*
TODO
- РАБОТАТЬ В ДРУГОЙ ВЕТКЕ на винде
- обработать стр 97 возврат если там капча
	по сути маршаллить ни во что не надо, это сделает grpc
- SearchByCondition возвращает результат []byte или return parsePage(page), где тип - готовый json?
	SearchByCondition все-таки будет вызываться по grpc, нет смысла делать Work()
- можно уведмолять в тг / редис-сообщение (подписка), если весь круг капч прошел, но результата все нет
- проверить, может ли быть неск кл слов и добавить цикл в таком случае
- мб понадобится решение капчи - отдельный сервис?
- конфиг для этой площадки
- остальные мб убрать
- сделать тест, чтобы prepareUrl возвращал корректный урл, сравнивать с урлами из логов

ТЕСТЫ!
*/

func (p *Platform) SearchByCondition(c *Condition, proxie []string) []byte {
	c = &Condition{ //ПРИМЕР
		Keyword:  "NAME",
		Type:     "фильм",
		Genres:   []string{"боевик"},
		YearFrom: "2021",
		YearTo:   "2023",
		Coutries: []string{"США"},
	}

	// заходим на страницу поиска, там перечисления для запроса, типа США это "1" и тд (для стран и жанров)
	req := gorequest.New()
	resp, body, errs := req.Get(searchUrl).End()
	for i := range errs {
		if errs[i] != nil {
			logrus.Errorf("error(%d) while getting %s, %v", i, searchUrl, errs[i])
		}
	}

	// если редирект на решение капчи, меняем прокси и делаем новый запрос
	for strings.Contains(resp.Request.URL.String(), "captcha") {
		p := nextProxy(proxie)
		reqWithProxy := gorequest.New().Proxy("http://" + p())
		resp, body, errs = reqWithProxy.Get(searchUrl).End()
		for i := range errs {
			if errs[i] != nil {
				logrus.Errorf("error(%d) while getting %s, %v", i, searchUrl, errs[i])
			}
		}
		time.Sleep(time.Second * 5)
	}

	// os.WriteFile("temp/search_page.htm", []byte(body), 0644)
	// body, _ := os.ReadFile("temp/search_page.htm")

	url := prepareUrl(c, string(body))
	resp, body, errs = req.Get(url).End() //проверка ответа на капчу тут?
	results := processPage(body)

	fmt.Println(results)
	return nil
}

// возвращает подготовленный урл для запроса
func prepareUrl(c *Condition, page string) string {
	u := url.URL{
		Scheme: "https",
		Host:   mainUrl,
	}

	q := u.Query()
	q.Set("level", "7")
	q.Set("from", "forma")
	q.Set("result", "adv")
	q.Set("m_act[from]", "forma")
	q.Set("m_act[what]", "content")

	if t := typeByName(c.Type); t != "" {
		q.Set("m_act[content_find]", t)
	}

	if c.YearFrom != "" {
		q.Set("m_act[from_year]", c.YearFrom)
	}

	if c.YearTo != "" {
		q.Set("m_act[to_year]", c.YearTo)
	}

	for _, c := range c.Coutries {
		if val := counryNumberByName(c, page); val != "" {
			q.Set("m_act[country]", val)
		}
	}

	for _, g := range c.Genres {
		if val := genreNumberByName(g, page); val != "" {
			q.Set("m_act[genre][]", val)
		}
	}

	q.Set("m_act[find]", c.Keyword) //цикл? - может ли быть несколько слов как параметры или все сразу идут?
	u.RawQuery = q.Encode()
	return strings.Replace(u.String(), "?", "/index.php?", 1) //no better solution found
}

// возвращает номера стран вместо названий
func counryNumberByName(coutry, page string) string {
	res, err := valueByLocator(coutry, page, "#country")
	if err != nil {
		logrus.Errorf("cannot find country num, %v", err)
		return ""
	}
	return res
}

// возвращает номера жанров вместо названий
func genreNumberByName(genre, page string) string {
	res, err := valueByLocator(genre, page, `#m_act\[genre\]`)
	if err != nil {
		logrus.Errorf("cannot find genre num, %v", err)
		return ""
	}
	return res
}

func valueByLocator(coutry, page, locator string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		logrus.Errorf("cannot create NewDocumentFromReader, %v", err)
	}

	options := doc.Find(locator).Find("option")

	if options.Size() == 0 {
		return "", fmt.Errorf("cannot parse <%s> with locator <%s> on page", coutry, locator)
	}

	for i := 0; i < options.Size(); i++ {
		if strings.Contains(coutry, options.Eq(i).Text()) {
			if val, ok := options.Eq(i).Attr("value"); ok {
				return val, nil
			}
		}
	}
	return "", fmt.Errorf("not found values for <%s> with locator <%s> on page", coutry, locator)
}

// возвращает тип по запросу / части запроса
func typeByName(in string) string {
	if strings.Contains(in, "сериал") {
		return "serial"
	}
	if strings.Contains(in, "фильм") {
		return "film"
	}
	return ""
}

// возвращает следующий прокси из списка/конфига
func nextProxy(proxies []string) func() string {
	i := 0
	return func() string {
		i++
		if i >= len(proxies) {
			i = 0
		}
		return proxies[i]
	}
}

// ищет нужные данные на странице
func processPage(page string) []Result {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		logrus.Errorf("cannot create NewDocumentFromReader, %v", err)
	}

	elems := doc.Find("#block_left_pad > div > div.search_results.search_results_last > div")
	res := make([]Result, 0, elems.Size())

	for i := 1; i < elems.Size(); i++ {
		nameElem := elems.Eq(i).Find(".name > a")
		name := nameElem.Text()
		ref, ok := nameElem.Attr("href")
		if !ok {
			ref = "Нет"
		}
		if !strings.Contains(ref, "https") {
			ref = "https://" + mainUrl + ref
		}

		id := strings.TrimSuffix(ref, "/sr/1/")
		id = path.Base(id)
		imgUrl := fmt.Sprintf(imgUrlTempl, id)

		res = append(res, Result{
			Name:   name,
			Ref:    ref,
			ImgUrl: imgUrl,
		})
	}
	return res
}
