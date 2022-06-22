package platform

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

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

const (
	mainUrl   = "www.kinopoisk.ru" // В КОНФИГ
	queryUrl  = "www.kinopoisk.ru/s/index.php"
	searchUrl = "https://www.kinopoisk.ru/s/"
)

/*
TODO
- РАБОТАТЬ В ДРУГОЙ ВЕТКЕ на винде
- можно уведмолять в тг / редис-сообщение (подписка), если весь круг капч прошел, но результата все нет
- проверить, может ли быть неск кл слов и добавить цикл в таком случае
- прокси
	- мб понадобится решение капчи - отдельный сервис?
- конфиг для этой площадки
- остальные мб убрать
- сделать тест, чтобы prepareUrl возвращал корректный урл, сравнивать с урлами из логов
*/

func (p *Platform) SearchByCondition(c *Condition) []byte {
	c = &Condition{ //ПРИМЕР
		Keyword:  "NAME",
		Type:     "фильм",
		Genres:   []string{"боевик"},
		YearFrom: "2021",
		YearTo:   "2023",
		Coutries: []string{"США"},
	}

	req := gorequest.New()
	resp, body, errs := req.Get(searchUrl).End()
	for i := range errs {
		if errs[i] != nil {
			logrus.Errorf("error(%d) while getting %s, %v", i, searchUrl, errs[i])
		}
	}

	for strings.Contains(resp.Request.URL.String(), "captcha") {
		//смена прокси / решение капчи
		//написать ф-ю на подачу следующей капчи
		reqWithProxy := gorequest.New().Proxy("http://CveAzQG5:5XLvvTcV@45.140.63.156:53023") //брать из цикла-замыкания по конфигу каждый раз новую
		resp, body, errs = reqWithProxy.Get(searchUrl).End()
		for i := range errs {
			if errs[i] != nil {
				logrus.Errorf("error(%d) while getting %s, %v", i, searchUrl, errs[i])
			}
		}
		time.Sleep(time.Second * 5) //TODO из конфига
	}

	// os.WriteFile("temp/search_page.htm", []byte(body), 0644)
	// body, _ := os.ReadFile("temp/search_page.htm")

	url := prepareUrl(c, string(body))
	fmt.Println(url)
	// req := gorequest.New()
	// resp, body, errs = req.Get(url).End()
	// fmt.Println(resp, "\n\n\n", body, errs)
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
