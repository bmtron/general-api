package main

import (
    "database/sql"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
)

var db *sql.DB

func main() {
    router := gin.Default()
    db = setupDb()
    router.GET("/ip", getIps)
    router.GET("/ip/active", getActiveIps)
    router.PUT("/ip/deactivateall", deactivateIps)
    router.POST("/ip/new", addNewIp)

    router.Run("localhost:8080")
}

func getIps(c *gin.Context) {
    ipAdresses, _ := getAllIpAddresses(db)
    c.IndentedJSON(http.StatusOK, ipAdresses)
}

func getActiveIps(c *gin.Context) {
    ipAddresses, _ := getActiveIpAddresses(db)
    c.IndentedJSON(http.StatusOK, ipAddresses)
}

func deactivateIps(c *gin.Context) {
    rowsAffected, err := deactivateIpAddresses(db)
    fmt.Sprintf("%d total rows affected.", rowsAffected)
    if err != nil {
        return
    }
    c.IndentedJSON(http.StatusNoContent, gin.H{"result": "entities modified successfully"})
}

func addNewIp(c *gin.Context) {
    var newIp IpAddress

    if err := c.BindJSON(&newIp); err != nil {
        return
    }
    rowsAffected, err := insertNewIpAddress(db, newIp)
    fmt.Sprintf("%d total rows affected.", rowsAffected)
    if err != nil {
        return
    }
    c.IndentedJSON(http.StatusCreated, gin.H{"result": "entity added successfully"})

}
