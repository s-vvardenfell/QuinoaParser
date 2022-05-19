/*
Copyright © 2022 s.vvardenfell, hel0len53

*/
package main

import (
	"fmt"
	"parser/kinoafisha"
)

func main() {
	// cmd.Execute()
	// f, _ := os.Open("logs/Fiddler_10-45-15.htm")
	// res := processSearchResults(utility.BytesFromReader(f))
	// fmt.Println(res)

	ka := kinoafisha.New()
	res := ka.ParseSeriesCalendar()
	fmt.Println(res)
}
