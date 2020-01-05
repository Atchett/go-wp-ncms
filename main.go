package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// WPSiteURL - Global Wordpress site URL
var WPSiteURL string

func main() {

	fmt.Println("Starting the application...")

	// get the passed in flags
	siteURL := flag.String("url", "", "URL of the Wordpress site. Don't include trailing /")
	numberPerPage := flag.Int("num", 10, "Number of values to get per page from the Wordpress API.")
	refresh := flag.Bool("refresh", false, "Clear any existing stored data, getting the latest data from the API.")
	flag.Parse()

	if *siteURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	WPSiteURL = *siteURL

	// delete the output directory - forcing a refresh
	if *refresh {
		err := os.RemoveAll("data")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	err := wpDataFromURL(*siteURL, *numberPerPage)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// generate the output files
	err = generatePostFiles()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = generateAuthorFiles()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("Terminating the application...")
}
