package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/likexian/whois-go"
)

//Consumes and read the response of the SSLab Endpoint
func consumeSSLLabOfDomain(domainName string, responseObject *Host) {
	consumeSSLLabs(domainName, responseObject)
	//Retry petitions until all the information about servers and domain is READY, or stop when
	//an error ocurres
	for responseObject.Status != "READY" && responseObject.Status != "ERROR" {
		consumeSSLLabs(domainName, responseObject)
	}

}

//Auxiliar method of consumeSSLLabOfDomain, which make a get request and update the data of the host object
func consumeSSLLabs(domainName string, responseObject *Host) {

	//Get petition to the SSLab endpoint
	response, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + domainName)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		responseData, _ := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		//Retrieve JSON information
		json.Unmarshal(responseData, &responseObject)
	}
}

//Use the lib who-is, to get information about the country and owner of all servers of the domain
func consumeWhoIsOfDomainServers(responseObject *Host) {
	var i int
	//Traverse all servers to update its data about country and owner
	for i = 0; i < len(responseObject.Servers); i++ {
		consumeWhois(&responseObject.Servers[i])
	}
}

//Auxiliar method of consumeWhoIsOfDomainServers, which find the country and owner of a server, using its ip address
func consumeWhois(server *Server) {

	//Verify the server has a valid ip address
	if whois.IsIpv4(server.Ip) {

		result, err := whois.Whois(server.Ip)
		//Parse the result
		owner, country, err2 := parseWhoisInfo(result)
		if err != nil || err2 != nil {

			log.Fatal(err)
			log.Fatal(err2)
		}

		//update server's data
		server.Company = owner
		server.Country = country
	}

}

//auxiliar method of consumeWhois which search for the information about country and owner in the
//the who-is lib response
func parseWhoisInfo(text string) (owner, country string, err error) {

	var (
		TextReplacer = regexp.MustCompile(`\n\[(.+?)\][\ ]+(.+?)`)
	)
	whoisText := strings.Replace(text, "\r", "", -1)
	whoisText = TextReplacer.ReplaceAllString(whoisText, "\n$1: $2")

	whoisLines := strings.Split(whoisText, "\n")
	var findOwner, findCountry bool
	for i := 0; i < len(whoisLines); i++ {
		line := strings.TrimSpace(whoisLines[i])
		if len(line) < 5 || !strings.Contains(line, ":") {
			continue
		}

		fChar := line[:1]
		if fChar == ">" || fChar == "%" || fChar == "*" {
			continue
		}

		if line[len(line)-1:] == ":" {
			i++
			for ; i < len(whoisLines); i++ {
				thisLine := strings.TrimSpace(whoisLines[i])
				if strings.Contains(thisLine, ":") {
					break
				}
				line += thisLine + ","
			}
			line = strings.Trim(line, ",")
			i--
		}

		lines := strings.SplitN(line, ":", 2)
		name := strings.TrimSpace(lines[0])
		value := strings.TrimSpace(lines[1])

		if value == "" {
			continue
		}

		//Search for owner
		if name == "OrgName" || name == "owner" {
			owner = value
			findOwner = true
		}

		//Search for country
		if name == "Country" || name == "country" {
			country = value
			findCountry = true
		}

		//stop the loop when the country and owner is already found
		if findCountry && findOwner {
			break
		}

	}

	return
}

//Adapt the information obtained with the instance of 'Host' structure to the desired response
func createJSON(responseObject *Host) ([]byte, error) {

	var tmp struct {
		Servers          []Server `json:"servers"`
		ServersChanged   bool     `json:"servers_changed"`
		SslGrade         string   `json:"ssl_grade"`
		PreviousSslGrade string   `json:"previous_ssl_grade"`
		Logo             string   `json:"logo"`
		Title            string   `json:"title"`
		Down             bool     `json:"is_down"`
	}
	tmp.Servers = responseObject.Servers
	tmp.ServersChanged = responseObject.ServersChanged
	tmp.SslGrade = responseObject.SslGrade
	tmp.PreviousSslGrade = responseObject.PreviousSslGrade
	tmp.Logo = responseObject.Logo
	tmp.Title = responseObject.Title
	tmp.Down = responseObject.Down

	jsonAnswer, err := json.Marshal(tmp)
	return jsonAnswer, err

}
