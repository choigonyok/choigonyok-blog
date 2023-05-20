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

type writing struct {
	Title        string
	Body         string
	Written_Time string
	Category     string
}

type HTML_DATA struct {
	Data_text  template.HTML
	Data_title string
	Data_cate  string
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

	// elements page
	eg.GET("/elements.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "elements.html", nil)
	})

	// writing page
	eg.GET("/writing.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "writing.html", nil)
	})

	eg.GET("/project1.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "project1.html", nil)
	})

	// DB에서 프로젝트 read해서 게시판에 올리기
	eg.GET("/generic.html", func(c *gin.Context) {
		st := writing{}
		ss := []writing{}
		r, err := db.Query("SELECT title, body, datetime FROM projects order by pid desc")
		if err != nil {
			log.Fatalln("@@@DB connection error :", err, "@@@")
		}
		var t, b, wt string
		for r.Next() {
			r.Scan(&t, &b, &wt)
			st.Title = t
			st.Body = b
			st.Category = "projects" //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
			st.Written_Time = wt
			ss = append(ss, st)
		}
		c.HTML(http.StatusOK, "generic.html", ss)
	})

	// DB에서 리뷰 read해서 게시판에 올리기
	eg.GET("/review.html", func(c *gin.Context) {
		st := writing{}
		ss := []writing{}
		r, err := db.Query("SELECT title, body, datetime FROM review order by rid desc")
		if err != nil {
			log.Fatalln("@@@DB connection error :", err, "@@@")
		}
		var t, b, wt string
		for r.Next() {
			r.Scan(&t, &b, &wt)
			st.Title = t
			st.Body = b
			st.Category = "review" //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
			st.Written_Time = wt
			ss = append(ss, st)
		}
		c.HTML(http.StatusOK, "review.html", ss)
	})

	// DB에서 스터디 read해서 게시판에 올리기
	eg.GET("/study.html", func(c *gin.Context) {
		st := writing{}
		ss := []writing{}
		r, err := db.Query("SELECT title, body, datetime FROM study order by sid desc")
		if err != nil {
			log.Fatalln("@@@DB connection error :", err, "@@@")
		}
		var t, b, wt string
		for r.Next() {
			r.Scan(&t, &b, &wt)
			st.Title = t
			st.Body = b
			st.Category = "study" //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
			st.Written_Time = wt
			ss = append(ss, st)
		}
		c.HTML(http.StatusOK, "study.html", ss)
	})

	// index, category 기반으로 DB에서 글 읽어와서 내용 보여주기
	eg.GET("/lookpage", func(c *gin.Context) {
		index := c.Query("index")
		cate := c.Query("cate")
		id, err := strconv.Atoi(index)
		if err != nil {
			log.Fatalln(err)
		}
		id = id + 1

		var uid string
		switch cate {
		case "projects":
			uid = "pid"
		case "review":
			uid = "rid"
		case "study":
			uid = "sid"
		}

		query := "SELECT title, body FROM " + cate + " order by " + uid + " desc limit " + strconv.Itoa(id)
		r, err := db.Query(query)
		if err != nil {
			log.Fatalln("THRER IS NO POST : ", err)
		}
		var b, t string
		for r.Next() {
			r.Scan(&t, &b)
		}
		//HTML로 보낼 data를 담는 구조체의 필드를 설정
		body_data1 := HTML_DATA{
			Data_text:  template.HTML(b),
			Data_title: t,
			Data_cate:  cate,
		}

		c.HTML(http.StatusOK, "project1.html", body_data1)
	})

	// Write한 글을 DB에 저장하고 쓴 게시물로 이동
	eg.POST("/markdown", func(c *gin.Context) {
		// 요청 데이터 가져오기
		fc := c.Request.FormValue("Cate")
		ftt := c.Request.FormValue("Title")
		ft := c.Request.FormValue("Text")

		ft = strings.ReplaceAll(ft, "\r\n", `\r\n`)
		data := struct {
			Text string `json:"text"`
		}{
			Text: ft,
		}
		// 외부 API에 POST 요청 보내기 - 본문관련
		url := "https://api.github.com/markdown"
		jsonStr := []byte(fmt.Sprintf(`{"text": "%s"}`, data.Text))
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to external API"})
			return
		}
		defer resp.Body.Close()
		// 응답 데이터 읽기
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from external API"})
			return
		}
		//HTML로 보낼 data를 담는 구조체의 필드를 설정
		body_data := HTML_DATA{
			Data_text:  template.HTML(string(body)),
			Data_title: ftt,
			Data_cate:  fc,
		}
		//
		//
		//

		// 클라이언트가 카테고리를 설정하지 않고 업로드했을 때
		if fc == "" || ft == "" || ftt == "" {
			http.Error(c.Writer, "THERE IS EMPTY BOX", http.StatusBadRequest)
		}
		// 클라이언트가 카테고리를 프로젝트로 설정하고 업로드했을 때
		if fc == "projects" {
			r, err := db.Query("SELECT pid FROM projects order by pid desc limit 1")
			if err != nil {
				log.Fatalln("@@@DB connection error :", err, "@@@")
			}
			defer r.Close()
			var p_id int
			for r.Next() {
				r.Scan(&p_id)
			}
			db.Query("INSERT INTO projects values (?,?,?,?)", p_id+1, ftt, data.Text, string(time.Now().Format(time.DateTime)))

			// 깃허브 마크다운 API 결과값 + 제목 + 카테고리 data를 project1.html로 전송
			c.HTML(http.StatusOK, "project1.html", body_data)
		}
		// 클라이언트가 카테고리를 스터디로 설정하고 업로드했을 때
		if fc == "study" {
			r, err := db.Query("SELECT sid FROM study order by sid desc limit 1")
			if err != nil {
				log.Fatalln("@@@DB connection error :", err, "@@@")
			}
			defer r.Close()
			var s_id int
			for r.Next() {
				r.Scan(&s_id)
			}
			db.Query("INSERT INTO study values (?,?,?,?)", s_id+1, ftt, data.Text, string(time.Now().Format(time.DateTime)))

			// 깃허브 마크다운 API 결과값 + 제목 + 카테고리 data를 project1.html로 전송
			c.HTML(http.StatusOK, "project1.html", body_data)
		}
		// 클라이언트가 카테고리를 리뷰로 설정하고 업로드했을 때
		if fc == "review" {
			r, err := db.Query("SELECT rid FROM review order by rid desc limit 1")
			if err != nil {
				log.Fatalln("@@@DB connection error :", err, "@@@")
			}
			defer r.Close()
			var r_id int
			for r.Next() {
				r.Scan(&r_id)
			}
			db.Query("INSERT INTO review values (?,?,?,?)", r_id+1, ftt, data.Text, string(time.Now().Format(time.DateTime)))

			// 깃허브 마크다운 API 결과값 + 제목 + 카테고리 data를 project1.html로 전송
			c.HTML(http.StatusOK, "project1.html", body_data)
		}
	})

	eg.Run(":8080")

}
