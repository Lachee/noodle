package noodle

import (
	"io/ioutil"
	"net/http"
)

//DownloadFile fetches a URL and return the bytes
func DownloadFile(url string) ([]byte, error) {

	//Download the URL
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	//Read the contents and return it
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//DownloadString fetches a string from the URL
func DownloadString(url string) (string, error) {
	response, err := DownloadFile(url)
	if err != nil {
		return "", err
	}
	return string(response), nil
}
