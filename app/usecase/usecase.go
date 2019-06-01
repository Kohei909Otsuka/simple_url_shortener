package usecase

import (
	"errors"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/entity"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/store"
)

func ShortenUrl(origin string, store store.UrlMapper) (string, error) {
	originUrl := entity.OriginalUrl(origin)
	shortenUrl := originUrl.Shorten()
	err := store.Write(originUrl, shortenUrl)
	if err != nil {
		return string(shortenUrl), errors.New("could not write to store")
	}
	return string(shortenUrl), nil
}

func RestoreUrl(shorten string, store store.UrlMapper) (string, error) {
	originUrl, err := store.Read(entity.ShortenUrl(shorten))
	if err != nil {
		return string(originUrl), errors.New("could not read from store")
	}
	return string(originUrl), nil
}
