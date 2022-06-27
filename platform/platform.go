package platform

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"parser/config"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

const (
	msgNoProxie       = "no proxie provided, cannot deal with captcha"
	msgGotCaptcha     = "got captcha, try to use proxie"
	msgAllProxieInUse = "all proxies has already in use"
)

type Platform struct {
	MainUrl     string
	QueryUrl    string
	SearchUrl   string
	ImgUrlTempl string
}

func New(c config.Config) *Platform {
	m, _ := base64.StdEncoding.DecodeString(c.Urls.MainUrl)
	q, _ := base64.StdEncoding.DecodeString(c.Urls.QueryUrl)
	s, _ := base64.StdEncoding.DecodeString(c.Urls.SearchUrl)
	i, _ := base64.StdEncoding.DecodeString(c.Urls.ImgUrlTempl)
	return &Platform{
		MainUrl:     string(m),
		QueryUrl:    string(q),
		SearchUrl:   string(s),
		ImgUrlTempl: string(i),
	}
}

/*
TODO
- проверка на "нет результатов" на странице поиска
- цикл с капчей сделать как-то покрасивее

- тесты на то что возвращается корректные ответы, просто прогнать вызов SearchByCondition с разными условиями
- прибраться в cmd - Long: `Usage: quinoa и тд
- ЛИНТЕР! / vet

*/

func (p *Platform) SearchByCondition(c *Condition, proxie []string) ([]Result, error) {
	// заходим на страницу поиска, там перечисления для запроса,
	// типа США это "1" и тд (для стран и жанров)
	req := gorequest.New()
	resp, body, errs := req.Get(p.SearchUrl).End()
	for i := range errs {
		if errs[i] != nil {
			logrus.Errorf(
				"error(%d) while getting %s, %v", i, p.SearchUrl, errs[i])
			return nil, fmt.Errorf(
				"cannot get %s without proxies (request conditions: %#v)", p.SearchUrl, *c)
		}
	}

	// если редирект на решение капчи, меняем прокси и делаем новый запрос
	counter := 0
	prs := nextProxy(proxie)

	for strings.Contains(resp.Request.URL.String(), "captcha") {
		logrus.Warning(msgGotCaptcha)

		if proxie == nil {
			logrus.Warning(msgNoProxie)
			return nil, errors.New(msgNoProxie)
		}

		if counter == len(proxie)-1 {
			logrus.Warning(msgAllProxieInUse)
			counter = 0
		}

		proxy := prs()
		logrus.Info("proxy in use: ", proxy)

		reqWithProxy := gorequest.New().Proxy("http://" + proxy)
		resp, body, errs = reqWithProxy.Get(p.SearchUrl).End()
		for i := range errs {
			if errs[i] != nil {
				logrus.Errorf("error(%d) while getting %s, %v", i, p.SearchUrl, errs[i])
				return nil, fmt.Errorf("cannot get %s with proxies (request conditions: %#v)", p.SearchUrl, *c)
			}
		}

		time.Sleep(time.Second * 1)
		counter++
	}

	url := prepareUrl(c, string(body), p.MainUrl)
	resp, body, errs = req.Get(url).End()
	for i := range errs {
		if errs[i] != nil {
			logrus.Errorf("error(%d) while getting %s, %v", i, url, errs[i])
		}
	}

	counter = 0
	for resp.StatusCode != 200 || strings.Contains(resp.Request.URL.String(), "captcha") {
		logrus.Warning(msgGotCaptcha)

		if proxie == nil {
			logrus.Warning(msgNoProxie)
			return nil, errors.New(msgNoProxie)
		}

		if counter == len(proxie)-1 {
			logrus.Warning(msgAllProxieInUse)
			counter = 0
		}

		proxy := prs()
		logrus.Info("proxy in use: ", proxy)

		reqWithProxy := gorequest.New().Proxy("http://" + proxy)
		resp, body, errs = reqWithProxy.Get(url).End()
		for i := range errs {
			if errs[i] != nil {
				logrus.Errorf("error(%d) while getting %s, %v", i, url, errs[i])
				return nil, fmt.Errorf("cannot get %s with proxies (request conditions: %#v)", url, *c)
			}
		}

		time.Sleep(time.Second * 1)
		counter++
	}

	res := processPage(body, p.MainUrl, p.ImgUrlTempl)
	for i := range res {
		res[i].Img = getImg(res[i].ImgUrl)
	}

	return res, nil
}

// возвращает подготовленный урл для запроса
func prepareUrl(c *Condition, page, mainUrl string) string {
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

	q.Set("m_act[find]", c.Keyword)
	u.RawQuery = q.Encode()
	return strings.Replace(u.String(), "?", "/index.php?", 1) // лучшего решения не придумал, url.Encode кодирует '/' и ':'
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

func valueByLocator(val, page, locator string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		logrus.Errorf("cannot create NewDocumentFromReader, %v", err)
	}

	options := doc.Find(locator).Find("option")

	if options.Size() == 0 {
		return "", fmt.Errorf("cannot parse <%s> with locator <%s> on page", val, locator)
	}

	for i := 0; i < options.Size(); i++ {
		if strings.Contains(val, options.Eq(i).Text()) {
			if val, ok := options.Eq(i).Attr("value"); ok {
				return val, nil
			}
		}
	}
	return "", fmt.Errorf("not found values for <%s> with locator <%s> on page", val, locator)
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

// возвращает следующий уникальный прокси из списка/конфига
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
func processPage(page, mainUrl, imgUrlTempl string) []Result {
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

func getImg(imgUrl string) string {
	req := gorequest.New()
	resp, body, errs := req.Get(imgUrl).End()
	for i := range errs {
		if errs[i] != nil {
			logrus.Errorf("error(%d) while getting %s, %v", i, imgUrl, errs[i])
		}
	}

	if resp.StatusCode != 200 || !strings.Contains(resp.Header.Get("Content-Type"), "image") {
		logrus.Warningf("error while getting image: resp.code: %d, cont.type: %s",
			resp.StatusCode, resp.Header.Get("Content-Type"))

		noImg, err := os.ReadFile("resources/no_image.jpg")
		if err != nil {
			logrus.Error("cannot read 'no_image.jpg'")
		}

		return base64.StdEncoding.EncodeToString(noImg)
	}
	return base64.StdEncoding.EncodeToString([]byte(body))
}
