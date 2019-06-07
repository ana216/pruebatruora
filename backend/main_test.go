package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//Testing the endpoint about get info of specific domain
func TestGetDomainServersEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/servers/truora.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetDomainsReviewedEndpoint)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestGetDomainsReviewedEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/servers/alldomains", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetDomainsReviewedEndpoint)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
