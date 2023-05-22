package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	Data_text      template.HTML
	Data_title     string
	Data_cate      string
	Data_id        int
	Data_imagepath string
}

// type PRS_id struct {
// 	ID      int
// 	Cate_id int
// }

type PRS_id struct {
	PID int //	'projects' table's primary key
	SID int //	'study' table's primary key
	RID int //	'review' table's primary key
}

type Recent struct {
	Sequence int
	Category string
	Post_num int
}

func main() {
	db, err := sql.Open("mysql", "root:andromeda0085@/blog")
	if err != nil {
		log.Fatalln("DB IS NOT CONNECTED")
	}
	defer db.Close()

	eg := gin.Default()
	eg.LoadHTMLGlob("./templates/**/*.html")
	eg.Static("/assets", "./assets")

	// default homepage + 아래에 최근 게시물 6개 표시
	eg.GET("/", func(c *gin.Context) {
		query := "SELECT COALESCE(p_id, '0'), COALESCE(s_id, '0'), COALESCE(r_id, '0') from whole order by id desc limit 6"
		r, err := db.Query(query)
		if err != nil {
			log.Fatalln("DB Connetion ERROR occured :", err)
		}
		data := []PRS_id{}
		for r.Next() { // 최근 post 보여줄 게 6개라서
			temp := PRS_id{}
			r.Scan(&temp.PID, &temp.SID, &temp.RID)
			data = append(data, temp)
		}

		fmt.Println(data)
		data_html := []DB_DATA{}
		for i := 0; i < 6; i++ {
			temp := DB_DATA{}
			if data[i].PID != 0 {
				query = "SELECT title, body, imagepath from projects where pid =" + strconv.Itoa(data[i].PID)
				r, err := db.Query((query))
				if err != nil {
					log.Fatalln("QUERY ERROR occured :", err)
				}
				for r.Next() {
					r.Scan(&temp.Title, &temp.Body, &temp.Imagepath)
					temp.Category = "Projects"
					temp.Id = data[i].PID
				}
			} else if data[i].SID != 0 {
				query = "SELECT title, body, imagepath from study where sid =" + strconv.Itoa(data[i].SID)
				r, err := db.Query((query))
				if err != nil {
					log.Fatalln("QUERY ERROR occured :", err)
				}
				for r.Next() {
					r.Scan(&temp.Title, &temp.Body, &temp.Imagepath)
					temp.Category = "Study"
					temp.Id = data[i].SID
				}
			} else {
				query = "SELECT title, body, imagepath from review where rid =" + strconv.Itoa(data[i].RID)
				r, err := db.Query((query))
				if err != nil {
					log.Fatalln("QUERY ERROR occured :", err)
				}
				for r.Next() {
					r.Scan(&temp.Title, &temp.Body, &temp.Imagepath)
					temp.Category = "Review"
					temp.Id = data[i].RID
				}
			}
			data_html = append(data_html, temp)
		}

		c.HTML(http.StatusOK, "index.html", data_html)
		// c.HTML(http.StatusOK, "index.html", data1)
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

	eg.GET("/elements", func(c *gin.Context) {
		c.HTML(http.StatusOK, "elements.html", nil)
	})

	// DB에서 데이터 가져와서 게시판 보여주기
	eg.GET("/postlist", func(c *gin.Context) {
		listname := c.Query("listname")
		var id string
		switch listname {
		case "Projects":
			id = "pid"
		case "Review":
			id = "rid"
		case "Study":
			id = "sid"
		}
		show_list_query := "SELECT * from " + listname + " order by " + id + " desc"
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
		fmt.Println(data)
		c.HTML(http.StatusOK, "postlist.html", data)
	})

	//write하면 DB에 저장
	eg.POST("/write", func(c *gin.Context) {
		image, err := c.FormFile("image")
		if err != nil {
			log.Fatalln("Image Uploading ERROR occured :", err)
		}
		c.SaveUploadedFile(image, "./assets/images/"+image.Filename) // 여기 다시 공부
		already := c.Request.FormValue("Id")
		cate := c.Request.FormValue("Cate")
		title := c.Request.FormValue("Title")
		body := c.Request.FormValue("Text")
		var id, id_whole string
		switch cate {
		case "Projects":
			id = "pid"
			id_whole = "p_id"
		case "Review":
			id = "rid"
			id_whole = "r_id"
		case "Study":
			id = "sid"
			id_whole = "s_id"
		}
		if cate == "" || title == "" || body == "" {
			http.Error(c.Writer, "THERE IS EMPTY BOX", http.StatusBadRequest)
		}
		if already != "" { // 원래 있던 게시글의 수정 내용을 저장
			query := "UPDATE " + cate + ` set title = "` + title + `" , body = "` + body + `" where ` + id + ` = ` + already
			_, err := db.Query(query)
			if err != nil {
				log.Fatalln("DB Update failed ERROR :", err)
			}
			c.Redirect(http.StatusSeeOther, "/postlist?listname="+cate)
		} else { // 새롭게 게시물을 작성할 때 DB에 저장
			// 카테고리에 맞는 다음 id는 몇 번인지 확인
			query := "SELECT " + id + " FROM " + cate + " order by " + id + " desc limit 1"
			r, err := db.Query(query)
			if err != nil {
				log.Fatalln("DB Uploading ERROR :", err)
			}
			defer r.Close()
			var idnum int
			for r.Next() {
				r.Scan(&idnum)
			}
			// 카테고리에 따라 맞는 DB table에 저장
			query = "INSERT INTO " + cate + ` values ( ` + strconv.Itoa(idnum+1) + `,"` + title + `","` + body + `","` + string(time.Now().Format(time.DateTime)) + `","/assets/images/` + image.Filename + `")`
			db.Query(query)

			// whole table의 다음 저장할 id num 찾기
			query = "SELECT id FROM whole ORDER BY id desc limit 1"
			r, err = db.Query(query)
			if err != nil {
				log.Fatalln("DB Uploading ERROR :", err)
			}
			defer r.Close()
			idnum_whole := 0
			for r.Next() {
				r.Scan(&idnum_whole)
			}
			// 맞는 카테고리에만 카테고리별 id num을 저장하고 나머지는 null로 초기화
			switch id_whole {
			case "p_id":
				query = "INSERT INTO whole values (" + strconv.Itoa(idnum_whole+1) + "," + strconv.Itoa(idnum+1) + ",NULL,NULL)"
				db.Query(query)
			case "s_id":
				query = "INSERT INTO whole values (" + strconv.Itoa(idnum_whole+1) + ",NULL,NULL," + strconv.Itoa(idnum+1) + ")"
				db.Query(query)
			case "r_id":
				query = "INSERT INTO whole values (" + strconv.Itoa(idnum_whole+1) + ",NULL," + strconv.Itoa(idnum+1) + ",NULL)"
				db.Query(query)
			}

			// 게시판으로 이동
			c.Redirect(http.StatusSeeOther, "/postlist?listname="+cate)
		}
	})

	//게시판에서 게시글을 누르면 게시글 내용을 markdown 적용된 상태로 확인 가능
	eg.GET("/markdown", func(c *gin.Context) {
		index := c.Query("index")
		cate := c.Query("listname")
		var id string
		switch cate {
		case "Projects":
			id = "pid"
		case "Review":
			id = "rid"
		case "Study":
			id = "sid"
		}
		query := "SELECT * from " + cate + " where " + id + " = " + index
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

		html_data := HTML_DATA{
			Data_text:  template.HTML(string(md_body)),
			Data_title: data.Title,
			Data_cate:  data.Category,
			Data_id:    data.Id,
		}
		c.HTML(http.StatusOK, "postlook.html", html_data)
	})

	// 게시글 삭제
	eg.GET("/delete", func(c *gin.Context) {
		index := c.Query("index")
		cate := c.Query("cate")
		var id string
		switch cate {
		case "Projects":
			id = "pid"
		case "Review":
			id = "rid"
		case "Study":
			id = "sid"
		}
		query := "SELECT imagepath from " + cate + " where " + id + " = " + index
		r, err := db.Query(query)
		if err != nil {
			log.Fatalln("DELETE Image ERROR occured :", err)
		}
		var imagepath string
		for r.Next() {
			r.Scan(&imagepath)
		}
		os.Remove("." + imagepath)

		query = "DELETE from " + cate + " where " + id + " = " + index // 해당 레코드 DB에서 지우기
		_, err = db.Query(query)
		if err != nil {
			log.Fatalln("Query has ERROR :", err)
		}

		c.Redirect(http.StatusSeeOther, "/postlist?listname="+cate)
	})

	//게시글 수정
	eg.GET("/modify", func(c *gin.Context) {
		index := c.Query("index")
		cate := c.Query("listname")
		var id string
		switch cate {
		case "Projects":
			id = "pid"
		case "Review":
			id = "rid"
		case "Study":
			id = "sid"
		}
		query := "SELECT * from " + cate + " where " + id + " = " + index
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
		c.HTML(http.StatusOK, "modify.html", data)
	})

	// // 검색 기능
	// eg.POST("/search", func(c *gin.Context) {
	// 	word := c.Request.FormValue("Find")
	// 	query := ""
	// 	db.Query()
	// })

	eg.Run(":8080")

}
