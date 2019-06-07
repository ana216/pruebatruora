package main

import (
	"context"
	"database/sql"
	"log"
)

const (
	//name of the database where is stored the data
	databaseName      string  = "pruebaTruora"
	//hour limit to check changes about servers and ssl grades
	maxTimeHour       float64 = 1
	//address to connect to the database
	connectionAddress string  = "postgresql://root@LAPTOP-8LU1DH0A:26257?sslmode=disable"
	//name of the driver
	driverName        string  = "postgres"
)

//map which stores each SSLstate with its respectively grade
var sSLGrade = map[string]int{
	"A":  4,
	"A+": 3,
	"B":  2,
	"C":  1,
}

//creates the connection to the database
func connectToDataBase() (db *sql.DB) {

	db, err := sql.Open(driverName, connectionAddress)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	return
}

//Update or insert the domain entered as a parameter
func updateDomain(currentDomain *Host) error {
	db := connectToDataBase()
	domainExists, errp := domainExists(context.Background(), db, currentDomain.Name)
	if errp != nil {
		log.Fatal(errp)
		return errp
	}

	if domainExists {

		//since the domain exists, we need to store the data before the domain is updated, to verify
		//the last reviewed hour and if the servers have changed
		prevDomain, errp := findDomainSQL(context.Background(), db, currentDomain.Name)
		if errp != nil {
			log.Fatal(errp)
			return errp
		}

		//update the domain's servers info
		err := updateDomainServers(context.Background(), db, currentDomain)
		if err != nil {
			log.Fatal(err)
			return err
		}

		//update the domain info
		err2 := updateDomainSQL(context.Background(), db, currentDomain)
		if err2 != nil {
			log.Fatal(err)
			return err
		}

		//compare previous data with the current data
		comparePreviousData(&prevDomain, currentDomain)
	} else {

		//the domain is not yet stored in the database, se we need to insert it
		err := insertDomainSQL(context.Background(), db, currentDomain)
		if err != nil {
			log.Fatal(err)
			return err
		}

		//insert the domain's servers
		err2 := updateDomainServers(context.Background(), db, currentDomain)
		if err2 != nil {
			log.Fatal(err2)
			return err2
		}
	}
	defer db.Close()

	return nil
}

//checks if a domain exists in the database
func domainExists(ctx context.Context, db *sql.DB, nameDomain string) (bool, error) {
	selectSQL := "SELECT name FROM " + databaseName + ".domains WHERE name = $1;"
	rows, err := db.QueryContext(ctx, selectSQL, nameDomain)
	if !rows.Next() {

		return false, err
	}
	return true, err

}

//checks if a server exists in the database
func serverExists(ctx context.Context, db *sql.DB, address string) (bool, error) {
	selectSQL := "SELECT * FROM " + databaseName + ".servers WHERE address = $1;"
	rows, err := db.QueryContext(ctx, selectSQL, address)
	if !rows.Next() {

		return false, err
	}
	return true, err

}

//update or insert in the database, the fieldas of the servers of a domain entered as a parameter
func updateDomainServers(ctx context.Context, db *sql.DB, domain *Host) error {

	for _, server := range domain.Servers {

		//update or insert each server in the database
		err := updateServer(ctx, db, server, domain)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	//delete in the database the servers that were eliminated 
	err := deleteInexistentServers(ctx, db, domain)
	return err

}

//update the domain entered as a parameter in the database
func updateDomainSQL(ctx context.Context, db *sql.DB, domain *Host) error {
	const updateDomainSQLStatement = "UPDATE " + databaseName + ".domains SET servers_changed = $1, ssl_grade = $2, previous_ssl_grade = $3, logo = $4, title = $5, is_down = $6, review_time = $7  WHERE name = $8;"
	if _, err := db.ExecContext(ctx, updateDomainSQLStatement, domain.ServersChanged, domain.SslGrade, domain.PreviousSslGrade, domain.Logo, domain.Title, domain.Down, domain.DateReview, domain.Name); err != nil {
		return err
	}
	return nil

}

//insert the domain entered as a parameter in the database
func insertDomainSQL(ctx context.Context, db *sql.DB, domain *Host) error {
	const createDomainSQLStatement = "INSERT INTO " + databaseName + ".domains VALUES ($1,$2,$3,$4,$5,$6,$7,$8);"
	if _, err := db.ExecContext(ctx, createDomainSQLStatement, domain.Name, domain.ServersChanged, domain.SslGrade, domain.PreviousSslGrade, domain.Logo, domain.Title, domain.Down, domain.DateReview); err != nil {
		return err
	}
	return nil
}

//insert the server entered as a parameter in the database
func insertServerSQL(ctx context.Context, db *sql.DB, server Server, domainName string) error {
	const insertServerSQLStatement = "INSERT INTO " + databaseName + ".servers VALUES( $1, $2, $3, $4, $5);"
	if _, err := db.ExecContext(ctx, insertServerSQLStatement, server.Ip, server.SslGrade, server.Country, server.Company, domainName); err != nil {
		return err
	}
	return nil
}

//update the server entered as a parameter in the database
func updateServerSQL(ctx context.Context, db *sql.DB, server Server, domainName string) error {
	const updateServerSQLStatement = "UPDATE " + databaseName + ".servers SET ssl_grade = $1, country = $2, owner = $3, domain_name = $4 WHERE address = $5;"
	if _, err := db.ExecContext(ctx, updateServerSQLStatement, server.SslGrade, server.Country, server.Company, domainName, server.Ip); err != nil {
		return err
	}
	return nil
}

//update or insert the server entered as a parameter in the database
func updateServer(ctx context.Context, db *sql.DB, server Server, domain *Host) error {

	//check if the server exists
	serverExists, err := serverExists(context.Background(), db, server.Ip)
	if err != nil {
		log.Fatal(err)
		return err
	}
	//update the SSLGrade attribute of the domain
	verifyLowestSslGrade(domain, &server)
	if serverExists {

		err := updateServerSQL(context.Background(), db, server, domain.Name)
		if err != nil {
			log.Fatal(err)
			return err
		}

	} else {
		//new server is added, thus the servers changed attribute has to be true
		domain.ServersChanged = true
		err := insertServerSQL(context.Background(), db, server, domain.Name)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

//find in the database info about the domain entered as a parameter
func findDomainSQL(ctx context.Context, db *sql.DB, domainName string) (Host, error) {
	var domainFound Host
	const selectDomainSQLStatement = "SELECT * FROM " + databaseName + ".domains WHERE name = $1;"
	if err := db.QueryRowContext(ctx, selectDomainSQLStatement, domainName).Scan(&domainFound.Name, &domainFound.ServersChanged, &domainFound.SslGrade, &domainFound.PreviousSslGrade, &domainFound.Logo, &domainFound.Title, &domainFound.Down, &domainFound.DateReview); err != nil {
		return domainFound, err
	}
	return domainFound, nil
}

//update the sslgrade of a domain entered as a paremeter, depending of the server also entered as a parameter 
func verifyLowestSslGrade(domain *Host, server *Server) {
	if domain.SslGrade != "" {
		sSlgradeDomain := sSLGrade[domain.SslGrade]
		sSLGradeServer := sSLGrade[server.SslGrade]

		if sSLGradeServer < sSlgradeDomain {
			domain.SslGrade = server.SslGrade
		}
	} else {
		domain.SslGrade = server.SslGrade
	}

}

//deletes in the database servers that are not in the current domain
func deleteInexistentServers(ctx context.Context, db *sql.DB, domain *Host) error {

	const selectSQL = "SELECT * FROM " + databaseName + ".servers WHERE domain_name = $1;"
	serverRows, _ := db.QueryContext(ctx, selectSQL, domain.Name)
	defer serverRows.Close()

	for serverRows.Next() {
		var server Server
		err := serverRows.Scan(&server.Ip, &server.SslGrade, &server.Country, &server.Company, &domain.Name)
		if err != nil {
			log.Fatal(err)
			return err
		}

		if !containsServer(domain.Servers, server.Ip) {
			err := deleteServer(ctx, db, &server)
			//A server is deleted, thus the servers changed attribute has to be true
			domain.ServersChanged = true
			if err != nil {
				log.Fatal(err)
				return err
			}
		}

	}
	return nil

}

//checks if an array of servers contains a specific a server which has a specific ipaddress
func containsServer(servers []Server, ipServer string) bool {
	found := false
	for _, server := range servers {
		if server.Ip == ipServer {
			found = true
			break
		}
	}
	return found
}

//deletes a server in the database
func deleteServer(ctx context.Context, db *sql.DB, server *Server) error {
	const deleteServerSQL = "DELETE FROM " + databaseName + ".servers WHERE address = $1;"
	_, err := db.ExecContext(ctx, deleteServerSQL, server.Ip)
	if err != nil {
		log.Fatal(err)
	}
	return err

}

//compares the current data with the previous data of a domain
func comparePreviousData(previousDomain *Host, currentDomain *Host) {

	//checks how many time has been transcurred
	diff := previousDomain.DateReview.Sub(currentDomain.DateReview).Hours()

	if diff < maxTimeHour {
		currentDomain.PreviousSslGrade = previousDomain.SslGrade
	} else {
		//Not enough time has been trasncurred, thus the attribute servers changed has to be false
		//and the previous ssl grade is the same of the current ssl grade
		currentDomain.ServersChanged = false
		currentDomain.PreviousSslGrade = currentDomain.SslGrade
	}

}

//allows to retrieve all the domains stored in the database
func selectAllDomains() (domains []Host, err2 error) {
	db := connectToDataBase()
	const selectSQL = "SELECT * FROM " + databaseName + ".domains;"
	serverDomains, _ := db.QueryContext(context.Background(), selectSQL)
	defer serverDomains.Close()
	for serverDomains.Next() {
		var domainFound Host
		err := serverDomains.Scan(&domainFound.Name, &domainFound.ServersChanged, &domainFound.SslGrade, &domainFound.PreviousSslGrade, &domainFound.Logo, &domainFound.Title, &domainFound.Down, &domainFound.DateReview)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		domains = append(domains, domainFound)
	}

	return

}
