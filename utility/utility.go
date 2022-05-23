package utility

import (
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func DoRequest(client *http.Client, method, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logrus.Fatalf("cannot create NewRequest to %s, %v", url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatalf("cannot get %s, %v", url, err)
	}
	return resp
}

func BytesFromReader(r io.Reader) []byte {
	byteValue, err := io.ReadAll(r)
	if err != nil {
		logrus.Fatalf("error reading in BytesFromReader(), %v", err)
	}
	return byteValue
}
