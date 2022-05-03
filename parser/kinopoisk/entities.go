package kinopoisk

type Film struct {
	Title         string
	Country       string
	ReleaseDate   string
	MiniPosterUrl string //большой постер - на странице фильма
	FilmUrl       string
	Genre         string //скорее всего можно узнать только на странице фильма
	Description   string //нужно спарсить на странице самого фильма
}
