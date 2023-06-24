package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Attraction struct {
	Id         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Detail     string `db:"detail" json:"detail"`
	Coverimage string `db:"coverimage" json:"coverimage"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/gomysql")
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	router := gin.Default()
	router.GET("/attractions", getAttractions)

	router.Use(cors.Default())
	router.Run("localhost:8080")
}

func getAttractions(c *gin.Context) {
	var attractions []Attraction
	rows, err := db.Query("select id, name, detail, coverimage from attractions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var a Attraction
		err := rows.Scan(&a.Id, &a.Name, &a.Detail, &a.Coverimage)
		if err != nil {
			log.Fatal(err)
		}
		attractions = append(attractions, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, attractions)
}
