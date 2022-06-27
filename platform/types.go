package platform

type Result struct {
	Name   string
	Ref    string
	ImgUrl string
	Img    string //string because gorequest returns the response body as a string - avoid multiple conversionsto []byte
}

type Condition struct {
	Keyword  string
	Type     string
	Genres   []string
	YearFrom string
	YearTo   string
	Coutries []string
}
