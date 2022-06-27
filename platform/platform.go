package platform

import (
	"encoding/base64"
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
- протестить получение картинки с новой логикой!
- если ф-я на получение новых прокси работает корректно, почему в реальных запросах всегда 1 прокся?????????

- сделать break и возврат ошибки, если прошелся по всем прокси и не получил результата
- цикл с капчей сделать как-то покрасивее
- что-то возвращать если ошибка в цикле с прокси/просто гет запроса? пока просто пишу варнинг
	например, в getImg есть только варниги, нет фаталов или возврата - что за бред?
- проверка на "нет результатов" на странице поиска

- возможно, стоит возвращать из все ф-й ошибки

- еще сервисы: - это в общий TODO
	- решение капчи
	- уведомление в тг
		вообще, можно уведомлятор в тг сделать как подписку на очередь в редисе

- тесты на то что возвращается корректные ответы, просто прогнать вызов SearchByCondition с разными условиями
- для других сервисов тоже можно сделать тесты
- прибраться в cmd - Long: `Usage: quinoa и тд
- ЛИНТЕР! / vet

- grpc

*/

//мб должен возвр err
func (p *Platform) SearchByCondition(c *Condition, proxie []string) []Result {
	// c = &Condition{ //ПРИМЕР
	// 	Keyword:  "NAME",
	// 	Type:     "фильм",
	// 	Genres:   []string{"боевик"},
	// 	YearFrom: "2021",
	// 	YearTo:   "2023",
	// 	Coutries: []string{"США"},
	// }

	// заходим на страницу поиска, там перечисления для запроса, типа США это "1" и тд (для стран и жанров)
	req := gorequest.New()
	resp, body, errs := req.Get(p.SearchUrl).End()
	for i := range errs {
		if errs[i] != nil {
			logrus.Errorf("error(%d) while getting %s, %v", i, p.SearchUrl, errs[i])
		}
	}

	// если редирект на решение капчи, меняем прокси и делаем новый запрос
	counter := 0
	for strings.Contains(resp.Request.URL.String(), "captcha") {
		logrus.Warning("got captcha, try to use proxie")

		if counter == len(proxie)-1 {
			logrus.Warning("all proxies has already in use")
			counter = 0
		}

		if proxie == nil {
			logrus.Warning("no proxies provided, cannot deal with captcha") //мб возвращать err, чтобы передать в grpc
			return nil
		}
		prs := nextProxy(proxie)
		prxy := prs()
		fmt.Println(prxy)
		reqWithProxy := gorequest.New().Proxy("http://" + prxy)
		resp, body, errs = reqWithProxy.Get(p.SearchUrl).End()
		for i := range errs {
			if errs[i] != nil {
				logrus.Errorf("error(%d) while getting %s, %v", i, p.SearchUrl, errs[i])
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
		logrus.Warning("got captcha, try to use proxie")

		if proxie == nil {
			logrus.Warning("no proxies provided, cannot deal with captcha") //мб возвращать err, чтобы передать в grpc
			return nil
		}

		if counter == len(proxie)-1 {
			logrus.Warning("all proxies has already in use")
			counter = 0
		}

		prs := nextProxy(proxie)
		reqWithProxy := gorequest.New().Proxy("http://" + prs())
		resp, body, errs = reqWithProxy.Get(url).End()
		for i := range errs {
			if errs[i] != nil {
				logrus.Errorf("error(%d) while getting %s, %v", i, p.SearchUrl, errs[i])
			}
		}

		time.Sleep(time.Second * 1)
		counter++
	}

	// for i := range results {
	// 	fmt.Println(results[i].Name, results[i].Ref, results[i].ImgUrl)
	// 	os.WriteFile(fmt.Sprintf("temp/%d.jpg", i), []byte(results[i].Img), 0644)
	// }
	res := processPage(body, p.MainUrl, p.ImgUrlTempl)
	for i := range res {
		res[i].Img = getImg(res[i].ImgUrl)
	}

	return res
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
