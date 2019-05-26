package entity_test

import (
	"github.com/Kohei909Otsuka/simple_url_shortener/app/entity"
	"net/url"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// before
	os.Setenv("BASE_URL", "https://shortener.com")

	code := m.Run()

	// after
	os.Setenv("BASE_URL", "")
	os.Exit(code)
}

func TestShorten(t *testing.T) {
	origin := entity.OriginalUrl("http://some-long-url.com")

	o1, err1 := url.Parse(string(origin.Shorten()))
	o2, err2 := url.Parse(string(origin.Shorten()))
	if err1 != nil || err2 != nil {
		t.Errorf("Shorten fail, should be parsed as url")
	}

	token1 := o1.Path[1:]
	token2 := o2.Path[1:]

	if o1.String() == o2.String() {
		t.Errorf("Shorten fail, should gen different urls one is %s, another is %s", o1, o2)
	}

	if len([]rune(token1)) != 6 || len([]rune(token2)) != 6 {
		t.Errorf("Shorten fail, token path should be 6 characters")
	}
}
