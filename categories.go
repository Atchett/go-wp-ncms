package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type categoryList []category

type category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var categoryData categoryList

var categoryDataDir = getDirectory(category{}, false)

// getCategoryData - gets the categories from the files saved from the WP API
func getCategoryData() error {

	// get a list of all the files in the dir
	fileList, err := ioutil.ReadDir(categoryDataDir)
	if err != nil {
		return err
	}

	for _, f := range fileList {
		jsonFile, err := os.Open(filepath.Join(categoryDataDir, f.Name()))
		if err != nil {
			return err
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var tmpcd categoryList
		// unmarshall the json into byte array
		err = json.Unmarshal(byteValue, &tmpcd)
		if err != nil {
			return err
		}
		// append the tag data
		categoryData = append(categoryData, tmpcd...)

	}
	return nil
}

// getCategoryName - gets the category name from the string ID
func getCategoryName(cID int) (string, error) {

	if len(categoryData) == 0 {
		err := getCategoryData()
		if err != nil {
			return "", err
		}
	}

	cn := "Not found"
	for i := range categoryData {
		if categoryData[i].ID == cID {
			cn = categoryData[i].Name
		}
	}
	return cn, nil
}
