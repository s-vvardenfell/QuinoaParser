package kinoafisha

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

//возвращает содержимое файла
func fileContent(filename string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(filepath.Dir(wd), filename)
	b, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Test_processSearchResults(t *testing.T) {
	b, err := fileContent("logs/search_result_page.txt")
	require.NoError(t, err)

	t.Log("\tTesting processing schedule page")
	{
		res := processSearchResults(b)
		require.NotEmpty(t, res)
	}
}

func Test_processSeriesPage(t *testing.T) {
	b, err := fileContent("logs/series_page.txt")
	require.NoError(t, err)

	t.Log("\tTesting processing series page")
	{
		res := processSeriesPage(b)
		require.NotEmpty(t, res.Descr, res.Genre)
		require.Empty(t, res.Name)
	}
}
