package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)
//main file which contains the main methods used to satisfy the system requirements
//main method
func main() {
	//create router
	r := chi.NewRouter()
	r.Get("/servers/{domainName}", GetDomainServersEndpoint)
	r.Get("/servers/alldomains", GetDomainsReviewedEndpoint)
	http.ListenAndServe(":2020", r)
}

//GetDomainServersEndpoint allows to get information about a specific domain and its servers
func GetDomainServersEndpoint(w http.ResponseWriter, req *http.Request) {
	var domain Host
	//instance the Host variable which is going to store all the required information
	domainName := string(chi.URLParam(req, "domainName"))
	//Consume the SSLLab service to get information about the domain and its servers
	consumeSSLLabOfDomain(domainName, &domain)
	//Use the who-is lib to get information about a specifi server
	consumeWhoIsOfDomainServers(&domain)
	//If an error didn't ocurre, we can search its title and logo
	if domain.Status != "ERROR" {
		updateTitle(&domain)
		updateLogo(&domain)

	} else {
		domain.Down = true
	}

	//update domain in the database
	err3 := updateDomain(&domain)
	//shape the response
	jsonAnswer, err := createJSON(&domain)

	if err != nil || err3 != nil {
		w.Write([]byte("error"))
	} else {
		w.Write(jsonAnswer)
	}

}

//GetDomainsReviewedEndpoint allows to get information about all domains that have been checked
func GetDomainsReviewedEndpoint(w http.ResponseWriter, req *http.Request) {

	//retrieve all domains
	domains, err := selectAllDomains()

	//struct that helps to create the desired json
	var tmp struct {
		Domains []string `json:"items"`
	}

	//traverse each domain to insert it in the database
	for _, domain := range domains {

		tmp.Domains = append(tmp.Domains, domain.Name+" info")

	}

	jsonAnswer, err2 := json.Marshal(tmp)

	if err != nil || err2 != nil {
		w.Write([]byte("error"))
	} else {
		w.Write(jsonAnswer)
	}

}
