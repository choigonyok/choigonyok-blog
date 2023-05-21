package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type DB_DATA struct {
	Title        string
	Body         string
	Id           int
	Written_time string
	Category     string
	Imagepath    string
}

type HTML_DATA struct {
	Data_text  template.HTML
	Data_title string
	Data_cate  string
	Data_id    int
}

func main() {
	db, err := sql.Open("mysql", "root:andromeda0085@/blog")
	if err != nil {
		log.Fatalln("@@@DB IS NOT CONNECTED@@@")
	}
	defer db.Close()

	eg := gin.Default()
	eg.LoadHTMLGlob("./templates/**/*.html")
	eg.Static("/assets", "./assets")

	// default homepage
	eg.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// redirect to homepage
	eg.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/")
	})

	// writing page
	eg.GET("/writing.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "writing.html", nil)
	})

	eg.GET("/project1.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "project1.html", nil)
	})

	// DB에서 데이터 가져와서 게시판 보여주기
	eg.GET("/postlist", func(c *gin.Context) {
		listname := c.Query("listname")
		show_list_query := "SELECT * from " + listname + " order by id desc"
		row, err := db.Query(show_list_query)
		data := []DB_DATA{}
		if err != nil {
			log.Fatalln("DB Connecting ERROR occured :", err)
		}
		for row.Next() {
			var temp DB_DATA
			row.Scan(&temp.Id, &temp.Title, &temp.Body, &temp.Written_time, &temp.Imagepath)
			temp.Category = listname
			data = append(data, temp)
		}

		fmt.Println(data[1].Category)
		c.HTML(http.StatusOK, "postlist.html", data)
	})

	//write하면 DB에 저장
	eg.POST("/write", func(c *gin.Context) {
		cate := c.Request.FormValue("Cate")
		title := c.Request.FormValue("Title")
		body := c.Request.FormValue("Text")
		if cate == "" || title == "" || body == "" {
			http.Error(c.Writer, "THERE IS EMPTY BOX", http.StatusBadRequest)
		}
		// 카테고리에 맞는 다음 id는 몇 번인지 확인
		query := "SELECT id FROM " + cate + " order by id desc limit 1"
		r, err := db.Query(query)
		if err != nil {
			log.Fatalln("DB Uploading ERROR :", err)
		}
		defer r.Close()
		var id int
		for r.Next() {
			r.Scan(&id)
		}
		// 카테고리에 따라 맞는 DB table에 저장
		query = "INSERT INTO " + cate + ` values ( ` + strconv.Itoa(id+1) + `,"` + title + `","` + body + `","` + string(time.Now().Format(time.DateTime)) + `",NULL)`
		fmt.Println(query)
		db.Query(query)
		// 게시판으로 이동
		c.Redirect(http.StatusSeeOther, "/postlist?listname="+cate)
	})

	//게시판에서 게시글을 누르면 게시글 내용을 markdown 적용된 상태로 확인 가능
	eg.GET("/markdown", func(c *gin.Context) {
		index := c.Query("index")
		cate := c.Query("listname")
		query := "SELECT * from " + cate + " where id = " + index
		data := DB_DATA{}
		row, err := db.Query(query)
		if err != nil {
			log.Fatalln("DB Connecting ERROR occured :", err)
			return
		}
		for row.Next() {
			row.Scan(&data.Id, &data.Title, &data.Body, &data.Written_time, &data.Imagepath)
			data.Category = cate
		}
		data.Body = strings.ReplaceAll(data.Body, "\r\n", `\r\n`)
		URL := "https://api.github.com/markdown"
		md_data := struct {
			Text string `json:"text"`
		}{
			Text: data.Body,
		}
		jsonStr := []byte(fmt.Sprintf(`{"text": "%s"}`, md_data.Text))
		res, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Fatalln("markdown API ERROR occured :", err)
		}
		defer res.Body.Close() // 왜 굳이 body?

		md_body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalln("Reading Markdown response ERROR occured :", err)
			return
		}
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println(string(md_body))
		fmt.Println(template.HTML(string(md_body)))
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

		html_data := HTML_DATA{
			Data_text:  template.HTML(string(md_body)),
			Data_title: data.Title,
			Data_cate:  data.Category,
			Data_id:    data.Id,
		}
		c.HTML(http.StatusOK, "project1.html", html_data)
	})

	eg.Run(":8080")

}
