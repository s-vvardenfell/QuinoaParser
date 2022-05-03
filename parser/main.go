/*
Copyright © 2022 s.vvardenfell, hel0len53

*/
package main

import (
	"log"
	"net/http"
	"parser/kinopoisk"
	"time"

	"github.com/gin-gonic/gin"
)

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

func main() {
	// cmd.Execute()

	// router := gin.Default()
	router = gin.Default()
	initializeRoutes()
	router.Run()
}

/*
curl "http://localhost:8080/Releases?countries=США&genres=ужасы&month=5"
curl "http://localhost:8080/Releases?countries=США&countries=Россия&genres=комедия&month=5"
*/
