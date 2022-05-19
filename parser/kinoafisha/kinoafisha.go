package kinoafisha

import (
	"bytes"
	"net/http"
	"net/http/cookiejar"
	"parser/utility"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type SeriesList map[string]string

type SeriesInfo struct {
	Name  string
	Descr string
	Genre []string
}

type KinoAfisha struct {
	client *http.Client
}

func New() *KinoAfisha {
	jar, err := cookiejar.New(nil)
	if err != nil {
		logrus.Fatalf("got error while creating cookie jar, %v", err)
	}

	var bot KinoAfisha
	bot.client = &http.Client{
		Transport: http.DefaultTransport,
		Jar:       jar,
		Timeout:   20 * time.Second,
	}
	return &bot
}

func (k *KinoAfisha) ParseSeriesCalendar() []SeriesInfo {
	resp := utility.DoRequest(
		k.client,
		http.MethodGet,
		"https://www.kinoafisha.info/series/calendar/",
		nil,
	)
	defer resp.Body.Close()
	b := utility.BytesFromReader(resp.Body)

	sl := processSearchResults(b)

	res := make([]SeriesInfo, 0, len(sl))

	for name, ref := range sl {
		ref, _, _ = strings.Cut(ref, "/seasons")
		resp := utility.DoRequest(
			k.client, //мб нужно создавать нового клиента/исп-ть новые прокси
			http.MethodGet,
			ref,
			nil,
		)

		defer resp.Body.Close()
		b := utility.BytesFromReader(resp.Body)

		si := processSeriesPage(b)
		si.Name = name
		res = append(res, si)
	}
	return res
}

func processSearchResults(page []byte) SeriesList {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		logrus.Fatalf("cannot create goquery document %v", err)
	}

	series := doc.Find(".seriesList_item")
	result := make(map[string]string)

	for i := 0; i < 10; /*series.Size()*/ i++ {
		title := series.Eq(i).Find(".seriesList_name").Text()
		title, _, _ = strings.Cut(title, "/")

		ref, ok := series.Eq(i).Find(".seriesList_ref").Attr("href")
		if !ok {
			logrus.Warningf("cannot find ref for %s", title)
		}
		result[title] = ref
	}
	return result
}

func processSeriesPage(page []byte) SeriesInfo {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		logrus.Fatalf("cannot create goquery document %v", err)
	}

	var si SeriesInfo
	genreItems := doc.Find(".filmInfo_genreItem")
	for i := 0; i < genreItems.Size(); i++ {
		genre, ok := genreItems.Eq(i).Attr("href")
		if !ok {
			logrus.Warning("cannot find genre")
		}
		genre = path.Base(genre)
		si.Genre = append(si.Genre, genre)
	}

	descr := doc.Find(".filmInfo_descText").Text()
	descr, _, _ = strings.Cut(descr, "Еще")
	si.Descr = strings.TrimSpace(descr)
	return si
}
