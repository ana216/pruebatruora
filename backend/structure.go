package main

import (
	"time"
)

//Host Struct to map every domain.
type Host struct {
	Servers          []Server `json:"endpoints"`
	ServersChanged   bool     `json:"servers_changed"`
	SslGrade         string   `json:"ssl_grade"`
	PreviousSslGrade string   `json:"previous_ssl_grade"`
	Logo             string   `json:"logo"`
	Title            string   `json:"title"`
	Down             bool     `json:"is_down"`
	DateReview       time.Time

	Name   string `json:"host"`
	Status string `json:"status"`
}

//Server Struct to map every server.
type Server struct {
	Ip       string `json:"ipAddress"`
	SslGrade string `json:"grade"`
	Country  string `json:"country"`
	Company  string `json:"owner"`
}
