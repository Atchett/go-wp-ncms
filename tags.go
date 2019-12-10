package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type tagList []tag

type tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var tagData tagList

var tagDataDir = getDirectory(tag{}, false)

// getTagData - gets the tags from the json file
func getTagData() error {

	// get a list of all the files in the dir
	fileList, err := ioutil.ReadDir(tagDataDir)
	if err != nil {
		return err
	}

	for _, f := range fileList {
		jsonFile, err := os.Open(filepath.Join(tagDataDir, f.Name()))
		if err != nil {
			return err
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var tmptd tagList
		// unmarshall the json into byte array
		err = json.Unmarshal(byteValue, &tmptd)
		if err != nil {
			return err
		}
		tagData = append(tagData, tmptd...)

	}
	return nil
}

// getTagName - gets the tag name from the string ID
func getTagName(tID int) (string, error) {

	if len(tagData) == 0 {
		err := getTagData()
		if err != nil {
			return "", err
		}
	}

	t := "Not found"
	for i := range tagData {
		if tagData[i].ID == tID {
			t = tagData[i].Name
		}
	}
	return t, nil
}
