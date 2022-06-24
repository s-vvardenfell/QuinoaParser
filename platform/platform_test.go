package platform

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_prepareUrl(t *testing.T) {
	var page, urls []byte
	t.Log("\tReading urls list")
	{
		wd, err := os.Getwd()
		require.NoError(t, err)
		path := filepath.Join(filepath.Dir(wd), "resources/tests_data/urls.zip")
		urls = readZip(path)
		require.NoError(t, err)
	}

	t.Log("\tReading results page")
	{
		wd, err := os.Getwd()
		require.NoError(t, err)
		path := filepath.Join(filepath.Dir(wd), "resources/tests_data/search_page.zip")
		page = readZip(path)
		require.NoError(t, err)
	}

	res := prepareUrl(&Condition{
		Keyword:  "NAME",
		Type:     "фильм",
		Genres:   []string{"боевик"},
		YearFrom: "2021",
		YearTo:   "2023",
		Coutries: []string{"США"},
	}, string(page), "www.kinopoisk.ru")

	want, err := url.Parse(string(urls))
	require.NoError(t, err)
	got, err := url.Parse(res)
	require.NoError(t, err)

	require.Equal(t, want.Query().Get("m_act[find]"), got.Query().Get("m_act[find]"))
	require.Equal(t, want.Query().Get("m_act[to_year]"), got.Query().Get("m_act[to_year]"))
	require.Equal(t, want.Query().Get("m_act[country]"), got.Query().Get("m_act[country]"))
	require.Equal(t, want.Query().Get("m_act[genre][]"), got.Query().Get("m_act[genre][]"))
	require.Equal(t, want.Query().Get("m_act[from_year]"), got.Query().Get("m_act[from_year]"))
	require.Equal(t, want.Query().Get("m_act[content_find]"), got.Query().Get("m_act[content_find]"))
}

func Test_valueByLocator(t *testing.T) {
	var page []byte
	t.Log("\tReading results page")
	{
		wd, err := os.Getwd()
		require.NoError(t, err)
		path := filepath.Join(filepath.Dir(wd), "resources/tests_data/search_page.zip")
		page = readZip(path)
		require.NoError(t, err)
	}

	t.Log("\tSearching country num")
	{
		res, err := valueByLocator("Россия", string(page), "#country")
		require.NoError(t, err)
		require.Equal(t, res, "2")
	}

	t.Log("\tSearching genre num")
	{
		res, err := valueByLocator("боевик", string(page), `#m_act\[genre\]`)
		require.NoError(t, err)
		t.Log(res)
		require.NotEmpty(t, res, "3")
	}
}

func Test_processPage(t *testing.T) {
	var page []byte
	t.Log("\tReading results page")
	{
		wd, err := os.Getwd()
		require.NoError(t, err)
		path := filepath.Join(filepath.Dir(wd), "resources/tests_data/results_page.zip")
		page = readZip(path)
		require.NoError(t, err)
	}

	t.Log("\tProcessing results page")
	{
		res := processPage(string(page), "http://example.com", "http://example.com")
		for i := range res {
			require.NotEmpty(t, res[i].Name, res[i].Ref, res[i].ImgUrl)
			require.Empty(t, res[i].Img)
		}
	}
}

func Test_nextProxy(t *testing.T) {
	tests := []struct {
		name    string
		proxies []string
	}{
		{
			name: "Check for giving diff values",
			proxies: []string{
				"some_proxy_1",
				"some_proxy_2",
				"some_proxy_3",
				"some_proxy_4",
				"some_proxy_5",
			},
		},
		{
			name: "Check for giving equal values",
			proxies: []string{
				"some_proxy_1",
				"some_proxy_2",
				"some_proxy_3",
				"some_proxy_1",
				"some_proxy_2",
			},
		},
	}

	prevVal := "previous value"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := nextProxy(tt.proxies)
			for range tt.proxies {
				got := p()
				require.NotEqual(t, got, prevVal)
				prevVal = got
			}
		})
	}
}

func readZip(path string) []byte {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	str := new(bytes.Buffer)

	for _, f := range r.File {
		page, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(str, page)
		if err != nil {
			log.Fatal(err)
		}
		page.Close()
	}
	return str.Bytes()
}
