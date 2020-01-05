package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func wpDataFromURL(URL string, numberPerPage int) error {

	var i interface{} // empty interface to pass into get default dir
	outputDir := directoryFromStruct(i, false)
	// check to see if the API output directory exists
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		// if it doesn't run the get
		for _, t := range contentTypes() {
			fmt.Fprintf(os.Stdout, "Getting %s data from WP API...\n", t)
			err := wpDataFromWpAPI(WPSiteURL, t, numberPerPage)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Wordpress API data saved to disk...")
	return nil
}

// wpDataFromWpAPI - access the data from the API for local store / processing
func wpDataFromWpAPI(URL string, contentType string, numberPerPage int) error {

	_, err := url.ParseRequestURI(URL)
	if err != nil {
		//fmt.Println(err)
		fmt.Println("Invalid URL. Must be in the format - http://...")
		return err
	}

	i := 1
	for {

		// the url to the API
		wpURL := fmt.Sprintf("%s/wp-json/wp/v2/%s?page=%d&per_page=%d", URL, contentType, i, numberPerPage)
		fmt.Println(wpURL)
		// Get the data
		resp, err := http.Get(wpURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// send the filepath
		apiOutputPath := fmt.Sprintf("data/api/%s", contentType)
		outputFileName := fmt.Sprintf("%s-%d.json", contentType, i)
		// if the response is ok, write a file
		if resp.StatusCode == 200 {
			// check the body length
			body, _ := ioutil.ReadAll(resp.Body)
			if len(body) < 10 {
				break
			}
			// write the response body to a file
			_, err := writeResponseToFile(apiOutputPath, outputFileName, body)
			if err != nil {
				return err
			}

		}
		// break the loop if the link doesn't work
		if resp.StatusCode != 200 {
			break
		}
		i++
	}
	return nil

}

func writeResponseToFile(path string, filename string, content []byte) (int64, error) {

	// check dir exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	// create a file path
	file, fileErr := os.Create(filepath.Join(path, filename))
	if fileErr != nil {
		return 0, fileErr
	}
	defer file.Close()

	// write to the file
	numBytes, writeErr := io.Copy(file, bytes.NewReader(content))
	//fmt.Println(numBytes)
	return numBytes, writeErr

}
