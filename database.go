package main

import (
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"net/url"
	"strconv"
)

type Result[T any] struct {
	Value T
}
type IpAddress struct {
	IpAddressId int    `json:"ipaddressid"`
	IpAddress   string `json:"ipaddress"`
	IsActive    bool   `json:"isactive"`
	DateCreated string `json:"datecreated"`
}
type IpRunLog struct {
	IpRunLogId int    `json:"iprunlogid"`
	RunDate    string `json:"rundate"`
	IpUpdated  bool   `json:"ipupdated"`
}

func setupDb() *sql.DB {
	var db *sql.DB
	//dbUser := os.Getenv("DB_USER")
	dbUser := "bmtron"
	//dbPass := os.Getenv("DB_PASS")
	dbPass := "Burtreynolds#30"
	//dbIp := os.Getenv("DB_IP")
	dbIp := "192.168.50.16"
	//dbPort := os.Getenv("DB_PORT")
	dbPort := "1433"
	dbPortAsInt, convertErr := strconv.Atoi(dbPort)
	if convertErr != nil {
		fmt.Errorf("error setting up db, cant parse env var DB_PORT")
		return nil
	}
	query := url.Values{}
	query.Add("database", "General")

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(dbUser, dbPass),
		Host:     fmt.Sprintf("%s:%d", dbIp, dbPortAsInt),
		RawQuery: query.Encode(),
	}

	var err error
	db, err = sql.Open("sqlserver", u.String())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(err)
	}

	return db
}

func getAllIpAddresses(db *sql.DB) ([]IpAddress, error) {
	var ips []IpAddress

	rows, err := db.Query(`select IpAddressId, IpAddress, IsActive, DateCreated from IpAddress;`)
	if err != nil {
		return nil, fmt.Errorf("ipaddresses %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ip IpAddress
		if err := rows.Scan(&ip.IpAddressId, &ip.IpAddress, &ip.IsActive, &ip.DateCreated); err != nil {
			return nil, fmt.Errorf("ipaddresses %v", err)
		}
		ips = append(ips, ip)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ipaddresses %v", err)
	}

	return ips, nil
}

func getActiveIpAddresses(db *sql.DB) (Result[[]IpAddress], error) {
	var ips []IpAddress

	rows, err := db.Query(`select IpAddressId, IpAddress, IsActive, DateCreated from IpAddress where IsActive = 1;`)
	if err != nil {
		return Result[[]IpAddress]{}, fmt.Errorf("ipadresses: active %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ip IpAddress
		if err := rows.Scan(&ip.IpAddressId, &ip.IpAddress, &ip.IsActive, &ip.DateCreated); err != nil {
			return Result[[]IpAddress]{}, fmt.Errorf("ipaddresses: active %v", err)
		}
		ips = append(ips, ip)
	}
	if err := rows.Err(); err != nil {
		return Result[[]IpAddress]{}, fmt.Errorf("ipaddresses %v", err)
	}
	var res = Result[[]IpAddress]{}
	res.Value = ips
	return res, nil
}

func deactivateIpAddresses(db *sql.DB) (int64, error) {

	res, err := db.Exec(`UPDATE IpAddress SET IsActive = 0;`)

	if err != nil {
		return 0, fmt.Errorf("ipaddress: deactivate %v", err)
	}

	rowsAffected, rowErr := res.RowsAffected()
	if rowErr != nil {
		return 0, fmt.Errorf("ipaddress: deactivate %v", rowErr)
	}
	return rowsAffected, nil
}

func insertNewIpAddress(db *sql.DB, address IpAddress) (int64, error) {
	res, err := db.Exec(`INSERT INTO IpAddress (IpAddress, IsActive, DateCreated) VALUES (@ipaddr, @active, getdate());`,
		sql.Named("ipaddr", &address.IpAddress), sql.Named("active", &address.IsActive))
	if err != nil {
		return 0, fmt.Errorf("ippaddress: add %v", err)
	}

	rowsAffected, rowErr := res.RowsAffected()
	if rowErr != nil {
		return 0, fmt.Errorf("ipaddress: deactivate %v", rowErr)
	}

	return rowsAffected, nil
}
func getAllRunLogs(db *sql.DB) (Result[[]IpRunLog], error) {
	var runLogs []IpRunLog
	var res Result[[]IpRunLog]
	rows, err := db.Query(`SELECT IpRunLogId, RunDate, IpUpdated FROM IpRunLog;`)
	if err != nil {
		return Result[[]IpRunLog]{}, fmt.Errorf("iprunlog: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var runLog IpRunLog
		if err := rows.Scan(&runLog.IpRunLogId, &runLog.RunDate, &runLog.IpUpdated); err != nil {
			return Result[[]IpRunLog]{}, fmt.Errorf("iprunlog: %v", err)
		}

		runLogs = append(runLogs, runLog)
	}

	if err := rows.Err(); err != nil {
		return Result[[]IpRunLog]{}, fmt.Errorf("iprunlog: %v", err)
	}
	res.Value = runLogs
	return res, nil
}

func addNewRunLog(db *sql.DB, runLog IpRunLog) (int64, error) {
	res, err := db.Exec(`INSERT INTO IpRunLog (RunDate, IpUpdated) VALUES (@rundate, @updated);`,
		sql.Named("rundate", runLog.RunDate), sql.Named("updated", runLog.IpUpdated))
	if err != nil {
		return 0, fmt.Errorf("add runlog: %v", err)
	}
	rowsAffected, rowErr := res.RowsAffected()
	if rowErr != nil {
		return 0, fmt.Errorf("add runlog: %v", err)
	}

	return rowsAffected, nil
}
