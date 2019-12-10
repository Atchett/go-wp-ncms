package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var mediaDataDir = getDirectory(media{}, false)
var mediaFilesDir = getDirectory(media{}, true)

type mediaList []media

var mediaData mediaList

type media struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	MediaType    string       `json:"media_type"`
	MediaDetails mediaDetails `json:"media_details"`
	PostID       int          `json:"post"`
	Slug         string       `json:"slug"`
}

type mediaDetails struct {
	Sizes sizes `json:"sizes"`
}

type sizes struct {
	Full full `json:"full"`
}

type full struct {
	FileName  string `json:"file"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	SourceURL string `json:"source_url"`
}

func getMediaData() error {

	// get a list of all the files in the dir
	fileList, err := ioutil.ReadDir(mediaDataDir)
	//fmt.Println(mediaDataDir)
	if err != nil {
		return err
	}

	// iterate over the files and read them
	for _, f := range fileList {
		//fmt.Println(filepath.Join(mediaDataDir, f.Name()))
		// Open the json file
		jsonFile, err := os.Open(filepath.Join(mediaDataDir, f.Name()))
		// handle error
		if err != nil {
			return err
		}

		defer jsonFile.Close()

		// read our opened json file as a byte array.
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var tmpmd mediaList
		// unmarshall the json into byte array
		err = json.Unmarshal(byteValue, &tmpmd)
		if err != nil {
			return err
		}
		// append the tag data
		mediaData = append(mediaData, tmpmd...)

	}
	return nil
}

// getMedianame - gets the media name from the passed in ID
func getMedianame(mID int, titleSlug string) (string, error) {

	if len(mediaData) == 0 {
		err := getMediaData()
		if err != nil {
			return "", err
		}
	}

	for i := range mediaData {
		if mediaData[i].ID == mID {
			// change the filename so that it's more obvious which post it belongs to
			fileExt := filepath.Ext(mediaData[i].MediaDetails.Sizes.Full.FileName)
			changedFilename := fmt.Sprintf("%s%s", titleSlug, fileExt)
			downloadedFile, err := getMediaFile(mediaData[i].MediaDetails.Sizes.Full.SourceURL, changedFilename)
			if err != nil {
				return "", err
			}
			return downloadedFile, nil
		}
	}

	return "", nil
}

// getMediaFile - access the data from the API for local store / processing
// URL - URL to the file for downloading
// f - the filename to use when storing the file
func getMediaFile(URL string, f string) (string, error) {

	fileToDownload := filepath.Join(mediaFilesDir, f)

	// check the file exists and don't get again
	if _, err := os.Stat(fileToDownload); os.IsNotExist(err) {
		// check output dir exists
		if _, err := os.Stat(mediaFilesDir); os.IsNotExist(err) {
			os.MkdirAll(mediaFilesDir, os.ModePerm)
		}
		// get the data
		response, err := http.Get(URL)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()

		// create the file
		file, err := os.Create(fileToDownload)
		if err != nil {
			return "", err
		}
		defer file.Close()

		// write to the file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			return "", err
		}
	}

	return fileToDownload, nil
}

func getFileExt(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	ct, err := getFileContentType(f)
	if err != nil {
		return "", err
	}
	fExt, err := fileExtFromMimeType(ct)
	if err != nil {
		return "", err
	}
	//fmt.Println("File Extension :", ext)
	return fExt, nil
}

// Read the first 512 bytes to detect the content type of the file
func getFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	defer out.Close()
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
