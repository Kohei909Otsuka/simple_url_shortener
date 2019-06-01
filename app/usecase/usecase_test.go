package usecase_test

import (
	"github.com/Kohei909Otsuka/simple_url_shortener/app/entity"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/usecase"
	"os"
	"testing"
)

type UrlMapperMock map[string]string

func (m UrlMapperMock) Write(o entity.OriginalUrl, s entity.ShortenUrl) error {
	m[string(s)] = string(o)
	return nil
}

func (m UrlMapperMock) Read(s entity.ShortenUrl) (entity.OriginalUrl, error) {
	return entity.OriginalUrl(m[string(s)]), nil
}

func TestMain(m *testing.M) {
	// before
	os.Setenv("BASE_URL", "https://shortener.com")

	code := m.Run()

	// after
	os.Setenv("BASE_URL", "")
	os.Exit(code)
}

func TestShortenThenRestore(t *testing.T) {
	var mock UrlMapperMock
	mock = make(map[string]string)
	shorten, _ := usecase.ShortenUrl("https://original.com", mock)
	restored, _ := usecase.RestoreUrl(shorten, mock)
	if restored != "https://original.com" {
		t.Errorf("can not shorten then restore, shorten was %s and restored was %s", shorten, restored)
	}
}
