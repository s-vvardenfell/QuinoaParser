package kinopoisk

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Films []Film

func (f *Films) ToJson() []byte {
	jsData, err := json.Marshal(f)
	if err != nil {
		logrus.Fatal("cannot marshall parsed films to json")
	}
	return jsData
}

type Film struct {
	Title         string
	Country       string
	ReleaseDate   string
	MiniPosterUrl string //большой постер - на странице фильма
	FilmUrl       string
	Genre         string //скорее всего можно узнать только на странице фильма
	Description   string //нужно спарсить на странице самого фильма
}
