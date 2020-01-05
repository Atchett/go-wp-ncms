package main

import (
	"fmt"
	"html"
	"path"
	"regexp"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

const (
	shortFormUK = "2006-01-02"
)

// formatToSlug - turns the name into a kebab case slug
func formatToSlug(n string) string {

	// make lowercase
	l := strings.ToLower(n)
	k := strings.ReplaceAll(l, " ", "-")

	return k

}

// formatHTMLToMD - converts the HTML content output by WP into markdown
func formatHTMLToMD(c string) (string, error) {
	converter := md.NewConverter("", true, nil)
	md, err := converter.ConvertString(c)
	if err != nil {
		return "", err
	}

	link := imagesFromMarkdown(md)
	if len(link) != 0 {
		//fmt.Fprintf(os.Stdout, "link: %s\n", link)

		stURL := link
		if !strings.Contains(link, "http") {
			stURL = fmt.Sprintf("%s%s", WPSiteURL, link)
		}
		//fmt.Fprintf(os.Stdout, "getting file: %s\n", stURL)
		mediaFileFromURL(stURL, path.Base(link))

	}

	return md, nil
}

// get the images from MD
func imagesFromMarkdown(md string) string {

	// Regex to extract title, link, and description
	re := regexp.MustCompile(`(?m)(^!\[\]\(([^)]+)\))?`)

	// Make regex
	match := re.FindStringSubmatch(md)

	if len(match) == 3 {
		return match[2]
	}

	return ""

}

// formatCleanHTML - ensures any html entities are cleaned
func formatCleanHTML(t string) string {

	// remove the HTML encoding
	c := html.UnescapeString(t)

	return c
}

// AppendDateToSlug appends UK formattted date string to string
// https://yourbasic.org/golang/format-parse-string-time-date-example/
func appendDateToSlug(d string, sl string) string {
	// get the string date into RFC3339 format (Must have a Z at the end - WP doesn't)
	t, _ := time.Parse(time.RFC3339, d+"Z")
	// conver date to UK format seperated by -
	fd := strings.Replace(t.Format(shortFormUK), " ", "-", -1)
	// join the strings with a -
	return strings.Join([]string{fd, sl}, "-")
}

// BetweenCurlyBraces returns the data between curly braces
// see - https://stackoverflow.com/a/39913100
func betweenCurlyBraces(s string) string {
	i := strings.Index(s, "{")
	if i >= 0 {
		j := strings.Index(s[i:], "}")
		if j >= 0 {
			// remove the whitespace
			return strings.TrimSpace(s[i+1 : j-i])
		}
	}
	return ""
}
