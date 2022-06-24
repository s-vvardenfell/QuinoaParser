package platform

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
