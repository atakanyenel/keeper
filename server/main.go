package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

var (
	repeat int
	db     *sql.DB
)

const distPath = "./dist/"

func getHostHandler(c *gin.Context) {
	query := fmt.Sprintf("SELECT ipaddress from hosts where hostname='%s'", c.Param("hostname"))
	rows, err := db.Query(query)
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error:%q", err))
	}
	defer rows.Close()
	for rows.Next() {
		var ipaddress string
		if err := rows.Scan(&ipaddress); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error:%q", err))
			return
		}
		c.String(http.StatusOK, ipaddress)
	}

}
func allHostsHandler(c *gin.Context) {
	rows, err := db.Query("SELECT * from hosts")
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error:%q", err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var hostname string
		var ipaddress string
		if err := rows.Scan(&hostname, &ipaddress); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading:%q", err))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("export keeper_%s=%s\n", strings.TrimSpace(hostname), ipaddress))
	}
}

func insertHostHandler(c *gin.Context) {
	hostname := c.PostForm("hostname")
	ipaddress := c.ClientIP()
	insertQuery := fmt.Sprintf("INSERT into hosts (hostname,ipaddress) values('%s','%s') ON CONFLICT(hostname) DO UPDATE SET hostname=excluded.hostname,ipaddress=excluded.ipaddress;", hostname, ipaddress)
	if _, err := db.Exec(insertQuery); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error:%q", err))
		return
	}

	c.String(http.StatusOK, "%s --> %s", hostname, ipaddress)
}

func rootHandler(c *gin.Context) {

	c.String(http.StatusOK, "Keeper is the solution for Dynamic IP Infrastructures")

}
func setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", rootHandler)

	hosts := router.Group("/hosts")
	{
		hosts.GET("/:hostname", getHostHandler)
		hosts.GET("/", allHostsHandler)
		hosts.POST("/", insertHostHandler)
	}
	return router
}
func main() {
	var err error
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening DB:%q", err)
	}
	router := setupRouter()

	router.Run(":" + port)
}
