package platform

// результат парсинга
type Result struct {
	Name   string
	Ref    string
	ImgUrl string
	Img    string
}

// условия для поиска
type Condition struct {
	Keyword  string
	Type     string
	Genres   []string
	YearFrom string
	YearTo   string
	Coutries []string
}
