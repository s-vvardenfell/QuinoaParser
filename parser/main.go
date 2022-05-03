package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// cmd.Execute()
	f, _ := os.Open("resources/Fiddler_19-57-29.htm")

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	sel := doc.Find(".prem_list")
	items := sel.Eq(1).Find(".premier_item")

	for i := 0; i < items.Size(); i++ {
		// fmt.Println(items.Eq(i).Find(`[itemprop="startDate"]`).Attr("content")) //дата выхода
		// fmt.Println(items.Eq(i).Find(`[itemprop="image"]`).Attr("content"))     //постер
		fmt.Println(items.Eq(i).Find(".text").Find(".name").Text())                              //название
		fmt.Println(items.Eq(i).Find(".text").Find(".name").Find("a:nth-child(1)").Attr("href")) //урл
		// fmt.Println(items.Eq(i).Find(`.text > div > span:nth-child(3)`).Text()) //страна и режиссер
		// fmt.Println(items.Eq(i).Find(`.text > div > span:nth-child(4)`).Text()) //жанр

	}
}

/*
#\38 35877 > meta:nth-child(1)
#\38 35877 > div:nth-child(4) > div:nth-child(1) > span:nth-child(3)
#\34 365786 > div:nth-child(4) > div:nth-child(1) > span:nth-child(1) > a:nth-child(1)
*/
