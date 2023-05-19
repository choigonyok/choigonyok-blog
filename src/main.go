package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:andromeda0085@/blog")
	if err != nil {
		log.Fatalln("@@@DB IS NOT CONNECTED@@@")
	}
	defer db.Close()

	eg := gin.Default()
	eg.LoadHTMLGlob("./templates/**/*.html")
	eg.Static("/assets", "./assets")
	//왜 images 폴더는 캐싱이 안될까?

	// default homepage
	eg.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// redirect to homepage
	eg.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/")
	})

	// generic page
	eg.GET("/generic.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "generic.html", nil)
	})

	// elements page
	eg.GET("/elements.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "elements.html", nil)
	})

	// writing page
	eg.GET("/writing.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "writing.html", nil)
	})

	// 클라이언트가 글을 업로드할 때
	eg.POST("/receive", func(c *gin.Context) {
		p_categoty := c.PostForm("post-category")
		p_title := c.PostForm("post-title")
		p_body := c.PostForm("post-write")
		// 클라이언트가 카테고리를 설정하지 않고 업로드했을 때
		if p_categoty == "" || p_title == "" || p_body == "" {
			http.Error(c.Writer, "THERE IS EMPTY BOX", http.StatusBadRequest)
		}
		// 클라이언트가 카테고리를 프로젝트로 설정하고 업로드했을 때
		if p_categoty == "projects" {
			r, err := db.Query("SELECT pid FROM projects order by pid desc limit 1")
			if err != nil {
				log.Fatalln("@@@DB connection error :", err, "@@@")
			}
			defer r.Close()
			var p_id int
			for r.Next() {
				r.Scan(&p_id)
			}
			db.Query("INSERT INTO projects values (?,?,?,?)", p_id+1, p_title, p_body, string(time.Now().Format(time.DateTime)))
		}
		// 클라이언트가 카테고리를 스터디로 설정하고 업로드했을 때
		if p_categoty == "study" {
			r, err := db.Query("SELECT sid FROM study order by sid desc limit 1")
			if err != nil {
				log.Fatalln("@@@DB connection error :", err, "@@@")
			}
			defer r.Close()
			var s_id int
			for r.Next() {
				r.Scan(&s_id)
			}
			db.Query("INSERT INTO study values (?,?,?,?)", s_id+1, p_title, p_body, string(time.Now().Format(time.DateTime)))
		}
		// 클라이언트가 카테고리를 리뷰로 설정하고 업로드했을 때
		if p_categoty == "review" {
			r, err := db.Query("SELECT rid FROM review order by rid desc limit 1")
			if err != nil {
				log.Fatalln("@@@DB connection error :", err, "@@@")
			}
			defer r.Close()
			var r_id int
			for r.Next() {
				r.Scan(&r_id)
			}
			db.Query("INSERT INTO review values (?,?,?,?)", r_id+1, p_title, p_body, string(time.Now().Format(time.DateTime)))
		}
		c.Redirect(http.StatusSeeOther, "/generic.html")
	})

	eg.Run(":8080")
}
