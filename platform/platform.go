package platform

import (
	"fmt"
	"net/url"
	"strings"
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
	mainUr    = "" // В КОНФИГ
	searchUrl = "www.kinopoisk.ru/s/index.php"
)

/*
TODO
- проверить, может ли быть неск кл слов и добавить цикл в таком случае
- прокси
- конфиг для этой площадки
- остальные мб убрать
- сделать тест, чтобы prepareUrl возвращал корректный урл, сравнивать с урлами из логов
*/

func (p *Platform) SearchByCondition(c *Condition) []byte {
	c = &Condition{ //ПРИМЕР
		Keyword:  "NAME2",
		Type:     "сериал",
		Genres:   []string{"биография", "боевик", "сериал"},
		YearFrom: "2021",
		YearTo:   "2022",
		Coutries: []string{"Россия"},
	}

	url := prepareUrl(c)
	fmt.Println(url)
	// req := gorequest.New()
	// resp, body, errs := req.Get(url).End()
	// fmt.Println(resp, "\n\n\n", body, errs)
	return nil
}

// возвращает подготовленный урл для запроса
func prepareUrl(c *Condition) string {
	url := url.URL{
		Scheme: "https",
		Host:   searchUrl,
	}

	q := url.Query()
	q.Add("level", "7")
	q.Add("from", "forma")
	q.Add("result", "adv")
	q.Add("m_act[from]", "forma")
	q.Add("m_act[what]", "content")

	if t := typeByName(c.Type); t != "" {
		q.Add("m_act[content_find]", t)
	}

	if c.YearFrom != "" {
		q.Add("m_act[from_year]", c.YearFrom)
	}

	if c.YearTo != "" {
		q.Add("m_act[to_year]", c.YearTo)
	}

	for _, c := range c.Coutries {
		if val := counryNumberByName(c); val != "" {
			q.Add("m_act[country]", val)
		}
	}

	for _, g := range c.Genres {
		q.Add("m_act[genre][]", genreNumberByName(g))
	}

	q.Add("m_act[find]", c.Keyword) //цикл? - может ли быть несколько слов как параметры или все сразу идут?
	url.RawQuery = q.Encode()

	return url.String()
}

// возвращает номера стран вместо названий
func counryNumberByName(coutry string) string {
	switch coutry {
	case "Россия": //как обработать, если придет "россия" (мал букв)
		return "2"
	default:
		return ""
	}
}

// возвращает номера жанров вместо названий
func genreNumberByName(coutry string) string {
	return "3"
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
