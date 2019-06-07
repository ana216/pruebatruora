package main

import (
	"net/http"

	"golang.org/x/net/html"
)

//verify if a html node is a title
func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

//verify if a html node is an icon
func isIconLogo(n *html.Node) (IsIcon bool, ref string) {
	if n.Type == html.ElementNode {
		for i := 0; i < len(n.Attr); i++ {
			if n.Attr[i].Key == "type" && n.Attr[i].Val == "image/x-icon" {
				IsIcon = true
			}
			if n.Attr[i].Key == "href" {
				ref = n.Attr[i].Val
			}
		}
	}
	return
}

//traverse the html structure to find the title
func traverseForTitle(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverseForTitle(c)
		if ok {
			return result, ok
		}
	}
	return "", false
}

//traverse the html structure to find the icon
func traverseForIcon(n *html.Node) (string, bool) {
	isIcon, ref := isIconLogo(n)
	if isIcon {
		return ref, isIcon
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverseForIcon(c)
		if ok {
			return result, ok
		}
	}
	return "", false
}

//GetHTMLTitle allows to get the title of a given html
func GetHTMLTitle(domainName string) (string, bool) {
	var response *http.Response
	var err error
	response, err = http.Get("https://" + domainName)
	if response == nil {
		response, err = http.Get("http://" + domainName)
	}
	doc, err3 := html.Parse(response.Body)
	if err != nil || err3 != nil {
		panic("Fail to parse html")
	}
	return traverseForTitle(doc)
}

//GetLogoURL allows to get the URL Logo of a given html
func GetLogoURL(domainName string) (string, bool) {
	var response *http.Response
	var err error
	response, err = http.Get("https://" + domainName)
	if response == nil {
		response, err = http.Get("http://" + domainName)
	}
	doc, err2 := html.Parse(response.Body)
	if err != nil || err2 != nil {
		panic("Fail to parse html")
	}
	return traverseForIcon(doc)
}

//update the logo of a given domain
func updateLogo(responseObject *Host) {
	logoURL, findLogo := GetLogoURL(responseObject.Name)
	if findLogo {
		responseObject.Logo = logoURL
	}
}

//update the title of a given domain
func updateTitle(responseObject *Host) {
	title, findTitle := GetHTMLTitle(responseObject.Name)
	if findTitle {
		responseObject.Title = title
	}

}
