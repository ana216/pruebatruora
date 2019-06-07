package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

//main method
func main() {
	r := chi.NewRouter()
	r.Get("/servers/{domainName}", GetDomainServersEndpoint)
	r.Get("/servers/alldomains", GetDomainsReviewedEndpoint)
	http.ListenAndServe(":2020", r)
}

//GetDomainServersEndpoint allows to get information about a specific domain and its servers
func GetDomainServersEndpoint(w http.ResponseWriter, req *http.Request) {
	var domain Host
	domainName := string(chi.URLParam(req, "domainName"))
	consumeSSLLabOfDomain(domainName, &domain)
	consumeWhoIsOfDomainServers(&domain)
	if domain.Status != "ERROR" {
		updateTitle(&domain)
		updateLogo(&domain)

	} else {
		domain.Down = true
	}

	err3 := updateDomain(&domain)
	jsonAnswer, err := createJSON(&domain)

	if err != nil || err3 != nil {
		w.Write([]byte("error"))
	} else {
		w.Write(jsonAnswer)
	}

}

//GetDomainsReviewedEndpoint allows to get information about all the domains that have been checked
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
