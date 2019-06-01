package usecase

import (
	"errors"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/entity"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/store"
)

func ShortenUrl(origin string, store store.UrlMapper) (string, error) {
	originUrl := entity.OriginalUrl(origin)
	shortenUrl := originUrl.Shorten()
	originUrlStr := string(originUrl)
	shortenUrlStr := string(shortenUrl)

	err := store.Write(originUrlStr, shortenUrlStr)
	if err != nil {
		return shortenUrlStr, errors.New("could not write to store")
	}
	return shortenUrlStr, nil
}

func RestoreUrl(shorten string, store store.UrlMapper) (string, error) {
	originUrlStr, err := store.Read(shorten)
	if err != nil {
		return originUrlStr, errors.New("could not read from store")
	}
	return originUrlStr, nil
}
