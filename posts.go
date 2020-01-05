package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Posts - a slice of posts from the post export
type posts []post

// Post - an individual post
type post struct {
	ID            int     `json:"id"`
	Date          string  `json:"date"`
	Title         title   `json:"title"`
	Content       content `json:"content"`
	Slug          string  `json:"slug"`
	Author        int     `json:"author"`
	FeaturedImage int     `json:"featured_media"`
	Categories    []int   `json:"categories"`
	Tags          []int   `json:"Tags"`
}

type title struct {
	Rendered string
}

type content struct {
	Rendered string
}

var postList posts

var postDataDir = directoryFromStruct(post{}, false)
var postsFilePath = directoryFromStruct(post{}, true)

// postDataFromFiles - gets the post data from the post files saved from the API
func postDataFromFiles() error {

	// get a list of all the files in the dir
	fileList, err := ioutil.ReadDir(postDataDir)
	if err != nil {
		return err
	}

	for _, f := range fileList {
		jsonFile, err := os.Open(filepath.Join(postDataDir, f.Name()))
		if err != nil {
			return err
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var tmppd posts
		// unmarshall the json into byte array
		err = json.Unmarshal(byteValue, &tmppd)
		if err != nil {
			return err
		}
		// append the tag data
		postList = append(postList, tmppd...)

	}
	return nil
}

// generatePostFiles - creates the markdown files of posts
func generatePostFiles() error {

	fmt.Println("Creating the post markdown files...")

	if len(postList) == 0 {
		err := postDataFromFiles()
		if err != nil {
			return err
		}
	}

	for _, post := range postList {

		// append the date to create the filename
		ad := appendDateToSlug(post.Date, post.Slug)

		// check dir exists
		if _, err := os.Stat(postsFilePath); os.IsNotExist(err) {
			os.MkdirAll(postsFilePath, os.ModePerm)
		}

		file, err := os.Create(filepath.Join(postsFilePath, filepath.Base(fmt.Sprintf("%v.md", ad))))
		if err != nil {
			return err
		}
		defer file.Close()

		w := bufio.NewWriter(file)

		// frontmatter
		fmt.Fprintln(w, "---")
		fmt.Fprintf(w, "title: %q\n", formatCleanHTML(post.Title.Rendered))
		fmt.Fprintln(w, "type: blog")
		fmt.Fprintln(w, "author:", authorNameFromID(post.Author))
		fmt.Fprintf(w, "date: %q\n", post.Date)
		mname, err := mediaNameFromID(post.FeaturedImage, post.Slug)
		if err != nil {
			mname = ""
		}
		fmt.Fprintln(w, "featuredImage: /assets/", mname)
		fmt.Fprintln(w, "featured: false")
		for _, cID := range post.Categories {
			catName, err := categoryNameFromID(cID)
			if err != nil {
				return err
			}
			fmt.Fprintln(w, "category:", catName)
		}
		fmt.Fprintln(w, "tags:")
		for _, tID := range post.Tags {
			tname, err := tagNameFromID(tID)
			if err != nil {
				return err
			}
			fmt.Fprintln(w, " -", tname)
		}
		fmt.Fprintln(w, "---")
		// end frontmatter - begin content
		content, err := formatHTMLToMD(post.Content.Rendered)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, content)
		// flush to the file
		w.Flush()
	}
	fmt.Println("Post markdown files created.")

	return nil
}
