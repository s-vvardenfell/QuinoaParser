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

/*
type Query struct {
	Countries []string `form:"countries"`
	Genres    []string `form:"genres"`
	Month     int      `form:"month"`
}

var router *gin.Engine

func ReleasesThisMonth(c *gin.Context) {
	kp := kinopoisk.New()
	t := time.Now()
	res := kp.Releases(nil, nil, int(t.Local().Month()))

	c.IndentedJSON(
		http.StatusOK,
		string(res.ToJson()),
	)
}

func ReleasesByQuery(c *gin.Context) {
	var q Query
	if c.ShouldBind(&q) == nil {
		log.Println(q.Countries)
		log.Println(q.Genres)
		log.Println(q.Month)
	}

	if err := c.ShouldBind(&q); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	kp := kinopoisk.New()
	res := kp.Releases(q.Countries, q.Genres, q.Month)

	c.IndentedJSON(
		http.StatusOK,
		string(res.ToJson()),
	)
}

func initializeRoutes() {
	router.GET("ReleasesThisMonth", ReleasesThisMonth) //по сути не нужна - с пустым query и так спарсится корректно
	router.GET("Releases", ReleasesByQuery)
}

type MyInt struct {
	val int
}

func (v MyInt) String() string {
	return "value is " + strconv.Itoa(v.val)
}


//main

	// router := gin.Default()
	// router = gin.Default()
	// initializeRoutes()
	// router.Run()



curl "http://localhost:8080/Releases?countries=США&genres=ужасы&month=5"
curl "http://localhost:8080/Releases?countries=США&countries=Россия&genres=комедия&month=5"

*/
