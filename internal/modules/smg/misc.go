package smg

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
)

//ParseHTML extracts all href urls within a html 'a' or 'base' tag
func ParseHTML(baseURL *url.URL, httpBody io.Reader) (urls []*url.URL, invalidURLs []string) {
	b, _ := io.ReadAll(httpBody)
	body := string(b)
	urls, invalidURLs = parseTags(baseURL, body)
	return urls, invalidURLs
}

func parseTags(baseURL *url.URL, body string) (urls []*url.URL, invalidURLs []string) {
	regexTags := regexp.MustCompile(`<\s*[a|base]\s+[^>]*href\s*=\s*[\"']?([^\"' >]+)[\"' >]`)
	regexHref := regexp.MustCompile(`href\s*=\s*[\"']?([^\"' >]+)[\"' >]`)
	tags := regexTags.FindAllString(body, -1)
	for _, tag := range tags {
		match := regexHref.FindAllString(tag, -1)
		if len(match) > 0 {
			value := match[0][6:]
			href := normaliseHref(strings.TrimSpace(value), "#")
			href = normaliseHref(strings.TrimSpace(value), "\"")
			if href == "" {
				invalidURLs = append(invalidURLs, value)
				continue
			}
			uri, err := formatURL(baseURL, href)
			if err != nil {
				invalidURLs = append(invalidURLs, href)
				continue
			}
			fmt.Println(uri)
			urls = append(urls, uri)
		}
	}
	return urls, invalidURLs
}

//formatURL returns a formatted and absolute url
func formatURL(baseURL *url.URL, href string) (*url.URL, error) {
	uri, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	uri = baseURL.ResolveReference(uri)
	if uri.Scheme == "" || uri.Host == "" {
		return nil, fmt.Errorf("url is invalid: %s", uri.String())
	}
	return uri, nil
}

//normaliseHref remove # and query parameters from a url
func normaliseHref(href string, identifier string) string {
	index := strings.Index(href, identifier)
	if index == -1 {
		return href
	}
	return href[:index]
}

// BuildSitemapFile create and complete sitemap file
func BuildSitemapFile(fileName string, urls map[string]int) error {
	err := deleteFileIfExists(fileName)
	if err != nil {
		return err
	}
	fh, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fh.Close()

	fh.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	fh.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	for loc := range urls {
		fh.WriteString("    " + "<url>\n")
		fh.WriteString("      " + "<loc>" + loc + "</loc>\n")
		fh.WriteString("      " + "<changefreq>weekly</changefreq>\n")
		fh.WriteString("      " + "<priority>0.5</priority>\n")
		fh.WriteString("    " + "</url>\n")
	}
	fh.WriteString("</urlset> ")

	return nil
}

//deleteFileIfExists deletes a file if exists
func deleteFileIfExists(fileName string) error {
	if _, err := os.Stat(fileName); err == nil {
		err := os.Remove(fileName)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
