package main

import (
    "database/sql"
    "fmt"
    _ "github.com/microsoft/go-mssqldb"
    "log"
    "net/url"
    "os"
    "strconv"
)

type IpAddress struct {
    IpAddressId int    `json:"ipaddressid"`
    IpAddress   string `json:"ipaddress"`
    IsActive    bool   `json:"isactive"`
    DateCreated string `json:"datecreated"`
}

func setupDb() *sql.DB {
    var db *sql.DB
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbIp := os.Getenv("DB_IP")
    dbPort := os.Getenv("DB_PORT")
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

func getActiveIpAddresses(db *sql.DB) ([]IpAddress, error) {
    var ips []IpAddress

    rows, err := db.Query(`select IpAddressId, IpAddress, IsActive, DateCreated from IpAddress where IsActive = 1;`)
    if err != nil {
        return nil, fmt.Errorf("ipadresses: active %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var ip IpAddress
        if err := rows.Scan(&ip.IpAddressId, &ip.IpAddress, &ip.IsActive, &ip.DateCreated); err != nil {
            return nil, fmt.Errorf("ipaddresses: active %v", err)
        }
        ips = append(ips, ip)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("ipaddresses %v", err)
    }

    return ips, nil
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
