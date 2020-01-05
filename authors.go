package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type authorList []author

type author struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Avatar      avatar `json:"avatar_urls"`
}

type avatar struct {
	URL string `json:"96"`
}

var authorData authorList
var authorDataDir = directoryFromStruct(author{}, false)
var authorsFilePath = directoryFromStruct(author{}, true)

// authorDataFromFiles - gets the author data from the files saved from the WP API
func authorDataFromFiles() {

	// get a list of all the files in the dir
	fileList, err := ioutil.ReadDir(authorDataDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range fileList {
		jsonFile, err := os.Open(filepath.Join(authorDataDir, f.Name()))
		if err != nil {
			log.Fatal(err)
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var tmpad authorList
		// unmarshall the json into byte array
		err = json.Unmarshal(byteValue, &tmpad)
		if err != nil {
			log.Fatal(err)
		}
		// append the tag data
		authorData = append(authorData, tmpad...)

	}
}

// authorNameFromID - returns the author name from the ID
func authorNameFromID(id int) string {

	if len(authorData) == 0 {
		authorDataFromFiles()
	}

	a := "Not found"
	for i := range authorData {
		if authorData[i].ID == id {
			a = authorData[i].Name
		}
	}
	return a
}

// generateAuthorFiles - creates the markdown author files
func generateAuthorFiles() error {

	fmt.Println("Creating the author markdown files...")

	for _, author := range authorData {

		// check dir exists
		if _, err := os.Stat(authorsFilePath); os.IsNotExist(err) {
			os.MkdirAll(authorsFilePath, os.ModePerm)
		}

		// create a file path
		file, fileErr := os.Create(filepath.Join(authorsFilePath, filepath.Base(fmt.Sprintf("%v.md", formatToSlug(author.Name)))))
		if fileErr != nil {
			return fileErr
		}
		defer file.Close()

		w := bufio.NewWriter(file)

		// frontmatter
		fmt.Fprintln(w, "---")
		fmt.Fprintf(w, "name: %q\n", author.Name)
		fmt.Fprintln(w, "type: author")
		fmt.Fprintln(w, "short_desc:", author.Description)
		avatar, err := avatarFromURL(author.Avatar.URL, author.Slug)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, "thumbnail:", fmt.Sprintf("/assets/%s", filepath.Base(avatar)))
		fmt.Fprintln(w, "---")
		// Bio - same as short_desc
		fmt.Fprintln(w, author.Description)
		// end content
		w.Flush()
	}
	fmt.Println("Author markdown files created.")
	return nil
}

// avatarFromURL - gets the avatar set in the WP API
func avatarFromURL(URL string, filename string) (string, error) {

	// download the file
	downloadedFile, err := mediaFileFromURL(URL, filename)
	if err != nil {
		return "", err
	}

	// check the filepath
	ext := filepath.Ext(downloadedFile)
	if len(ext) == 0 {
		// find the mime type
		// check the filename (if no extension)
		// get mime type
		ext, err = fileExtFromFile(downloadedFile)
		if err != nil {
			return "", err
		}
		// create a destination file with extension
		destFile := fmt.Sprintf("%s.%s", downloadedFile, ext)
		dest, err := os.Create(destFile)
		if err != nil {
			return "", err
		}
		defer dest.Close()

		// open the source file
		source, err := os.Open(downloadedFile)
		if err != nil {
			return "", err
		}
		defer source.Close()

		// copy file
		_, err = io.Copy(dest, source)
		if err != nil {
			return "", err
		}
		// delete the file with no extension
		err = os.Remove(downloadedFile)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(os.Stdout, "File downloaded : %s\n", destFile)

		return destFile, nil
	}
	fmt.Fprintf(os.Stdout, "File downloaded : %s\n", downloadedFile)

	return downloadedFile, nil
}
