package main

import (
	"errors"
)

//contentTypes - a list of types to use in the retreival of data from WP
func contentTypes() []string {
	return []string{"posts", "media", "users", "categories", "tags"}
}

// imageFileExtFromMimeType - returns the file extension based on the mime type
// t - mime type
func imageFileExtFromMimeType(t string) (string, error) {

	m := make(map[string]string)

	m["image/png"] = "png"
	m["image/jpeg"] = "jpg"
	m["image/gif"] = "gif"
	m["image/tiff"] = "tiff"
	m["image/bmp"] = "bmp"
	m["image/webp"] = "webp"
	m["image/svg+xml"] = "svg"

	if v, ok := m[t]; ok {
		return v, nil
	}

	err := errors.New("Mime type does not exist")

	return "", err

}

// directoryFromStruct - gets the directory to use based on the type of the struct passed in
// s - struct to use (e.g. post)
// isExport - wether this is for the exprt data
func directoryFromStruct(s interface{}, isExport bool) string {

	d := ""
	switch s.(type) {
	case post:
		d = "data/api/posts"
		if isExport {
			d = "data/export/posts"
		}
	case media:
		d = "data/api/media"
		if isExport {
			d = "data/export/files"
		}
	case author:
		d = "data/api/users"
		if isExport {
			d = "data/export/authors"
		}
	case category:
		d = "data/api/categories"
	case tag:
		d = "data/api/tags"
	default:
		d = "data/api"
	}
	return d
}
