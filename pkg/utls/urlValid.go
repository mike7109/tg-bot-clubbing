package utls

import "net/url"

func IsAddCmd(text string) bool {
	return IsURL(text)
}

func IsURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
