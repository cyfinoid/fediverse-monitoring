package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/fediverse/Db"
	"github.com/gin-gonic/gin"
)

func main() {
	err := Db.Connect()
	if err != nil {
		fmt.Println("Couldn't connect to Postgres DB")
		return
	}

	r := gin.Default()

	r.LoadHTMLGlob("static/templates/*.html")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Home.html", nil)
	})

	r.GET("/api/softwareUsed", func(c *gin.Context) {
		query := `SELECT name, COUNT(*) AS name_count
		FROM Nodes
		GROUP BY name
		ORDER BY name_count DESC
		LIMIT 10;`
		result, err := Db.Readpsql(query)
		if err != nil {
			panic(err)
		}
		j, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(j))
		c.Header("Content-Type", "application/json")
		c.Writer.Write(j)
	})

	r.GET("/api/softwareVersionUsed", func(c *gin.Context) {
		query := `SELECT name, version, COUNT(*) AS name_count
		FROM Nodes
		GROUP BY name, version
		ORDER BY name_count DESC
		LIMIT 10;`
		result, err := Db.Readpsql(query)
		if err != nil {
			panic(err)
		}
		j, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(j))
		c.Header("Content-Type", "application/json")
		c.Writer.Write(j)
	})

	r.GET("/api/tld", func(c *gin.Context) {
		query := `SELECT tld, COUNT(*) AS tld_count
		FROM Nodes
		GROUP BY tld
		ORDER BY tld_count DESC
		LIMIT 10;`
		result, err := Db.Readpsql(query)
		if err != nil {
			panic(err)
		}
		j, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(j))
		c.Header("Content-Type", "application/json")
		c.Writer.Write(j)
	})

	r.GET("/api/dailyServerCount", func(c *gin.Context) {
		query := `SELECT date_time, node_count from dailygrowth;`
		result, err := Db.Readpsql(query)
		if err != nil {
			panic(err)
		}
		j, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(j))
		c.Header("Content-Type", "application/json")
		c.Writer.Write(j)
	})

	r.GET("/api/dailyUserCount", func(c *gin.Context) {
		query := `SELECT date_time, total_users from dailygrowth;`
		result, err := Db.Readpsql(query)
		if err != nil {
			panic(err)
		}
		j, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(j))
		c.Header("Content-Type", "application/json")
		c.Writer.Write(j)
	})

	r.GET("/api/dailyPostCount", func(c *gin.Context) {
		query := `SELECT date_time, total_posts from dailygrowth;`
		result, err := Db.Readpsql(query)
		if err != nil {
			panic(err)
		}
		j, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(j))
		c.Header("Content-Type", "application/json")
		c.Writer.Write(j)
	})

	r.GET("/api/dailyCommentCount", func(c *gin.Context) {
		query := `SELECT date_time, total_comments from dailygrowth;`
		result, err := Db.Readpsql(query)
		if err != nil {
			panic(err)
		}
		j, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(j))
		c.Header("Content-Type", "application/json")
		c.Writer.Write(j)
	})

	r.Run()

}
