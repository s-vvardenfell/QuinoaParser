package main

import (
	"fmt"
	"parser/kinopoisk"
)

func main() {
	// cmd.Execute()
	kp := kinopoisk.Kinopoisk{}
	fmt.Println(kp.ReleasesByMonth(6))
}

/*
#\38 35877 > meta:nth-child(1)
#\38 35877 > div:nth-child(4) > div:nth-child(1) > span:nth-child(3)
#\34 365786 > div:nth-child(4) > div:nth-child(1) > span:nth-child(1) > a:nth-child(1)
*/
