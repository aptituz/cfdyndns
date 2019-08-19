package main

import (
	"io/ioutil"
	"net/http"
)

func getExternaIP() (string, error) {
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}

	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}
