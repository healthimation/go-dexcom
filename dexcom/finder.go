package dexcom

import (
	"net/url"
)

// Finder to use with base client
func findDexcom(serviceName string, useTLS bool) (url.URL, error) {
	ret, err := url.Parse("https://api.dexcom.com/")
	if err != nil || ret == nil {
		return url.URL{}, err
	}
	return *ret, err
}

func findDexcomSandbox(serviceName string, useTLS bool) (url.URL, error) {
	ret, err := url.Parse("https://sandbox-api.dexcom.com/")
	if err != nil || ret == nil {
		return url.URL{}, err
	}
	return *ret, err
}
