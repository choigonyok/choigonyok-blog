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
	Whole_id     int
	Visit_num    int
}

type HTML_DATA struct {
	Data_text      template.HTML
	Data_title     string
	Data_cate      string
	Data_id        int
	Data_imagepath string
	BeforeAndAfter []Post_DATA
}

type Post_DATA struct {
	ID       int
	Category string
	Where    string
}

type PRS_id struct {
	PID int
	SID int
	RID int
}

type Recent struct {
	Sequence int
	Category string
	Post_num int
}

type Visit struct {
	Data []PRS_id
	Tc   Timecheck
}

type Retrun_visit struct {
	Dd []DB_DATA
	Tc Timecheck
}

type Timecheck struct {
	Totalnum int
	Visitnum int
}

// 게시판 전체를 합친 데이터를 DB_DATA structure에 묶어서 리턴
func whole_cate(data Visit) Retrun_visit {
	db, err := sql.Open("mysql", "choigonyok:@/blog")// Change When Deploy
	if err != nil {
		log.Fatalln("DB IS NOT CONNECTED")
	}
	defer db.Close()
	data_html := Retrun_visit{}
	for i := 0; i < 6; i++ {
		temp := DB_DATA{}
		if data.Data[i].PID != 0 {
			query := "SELECT title, body, imagepath from Projects where pid =" + strconv.Itoa(data.Data[i].PID)
			r, err := db.Query((query))
			if err != nil {
				log.Fatalln("QUERY ERROR occured :", err)
			}
			for r.Next() {
				r.Scan(&temp.Title, &temp.Body, &temp.Imagepath)
				temp.Category = "Projects"
				temp.Id = data.Data[i].PID
			}
		} else if data.Data[i].SID != 0 {
			query := "SELECT title, body, imagepath from Study where sid =" + strconv.Itoa(data.Data[i].SID)
			r, err := db.Query((query))
			if err != nil {
				log.Fatalln("QUERY ERROR occured :", err)
			}
			for r.Next() {
				r.Scan(&temp.Title, &temp.Body, &temp.Imagepath)
				temp.Category = "Study"
				temp.Id = data.Data[i].SID
			}
		} else {
			query := "SELECT title, body, imagepath from Review where rid =" + strconv.Itoa(data.Data[i].RID)
			r, err := db.Query((query))
			if err != nil {
				log.Fatalln("QUERY ERROR occured :", err)
			}
			for r.Next() {
				r.Scan(&temp.Title, &temp.Body, &temp.Imagepath)
				temp.Category = "Review"
				temp.Id = data.Data[i].RID
			}
		}
		data_html.Dd = append(data_html.Dd, temp)
	}
	data_html.Tc = data.Tc
	return data_html
}

func main() {
	totalnum := 0
	visitnum := 0
	db, err := sql.Open("mysql", "choigonyok:@/blog") // Change When Deploy
	if err != nil {
		log.Fatalln("DB IS NOT CONNECTED")
	}
	defer db.Close()

	eg := gin.Default()
	eg.LoadHTMLGlob("./../templates/**/*.html")
	eg.Static("/assets", "./../assets")

	// default homepage + 아래에 최근 게시물 6개 표시
	eg.GET("/", func(c *gin.Context) {
		_, err := c.Cookie("visit")
		if err == http.ErrNoCookie {
			c.SetCookie("visit", "OK", 3600 , "", "", false, false)
			visitnum += 1
		}
		go func (){
			for(true){
			if time.Now().Format("15:04:05") == "00:00:00" {
				totalnum += visitnum
				visitnum = 0
				time.Sleep(time.Second*2)
			}
		}
		}()

		query := "SELECT COALESCE(p_id, '0'), COALESCE(s_id, '0'), COALESCE(r_id, '0') from whole order by id desc"
		r, err := db.Query(query)
		if err != nil {
			log.Fatalln("DB Connetion ERROR occured :", err)
		}
		data := []PRS_id{}
		for r.Next() {
			temp := PRS_id{}
			r.Scan(&temp.PID, &temp.SID, &temp.RID)
			data = append(data, temp)
		}
		tc := Timecheck{
			Totalnum: totalnum,
			Visitnum: visitnum,
		}
		data_html := Visit{
			data,
			tc,
		}

		c.HTML(http.StatusOK, "index.html", whole_cate(data_html))
	})

	// redirect to homepage
	eg.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/")
	})

	// writing page
	eg.GET("/writing.html", func(c *gin.Context) {
		value, err := c.Cookie("admistrator")
		if err != http.ErrNoCookie {
			if value == "OK" {
				c.HTML(http.StatusOK, "writing.html", nil)
			}
		}
		c.Redirect(http.StatusSeeOther, "/loginpage")

	})

	eg.GET("/project1.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "project1.html", nil)
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
		c.HTML(http.StatusOK, "postlist.html", data)
	})

	//write하면 DB에 저장
	eg.POST("/write", func(c *gin.Context) {
		image, err := c.FormFile("image")
		var imagename string
		if err != nil {
			imagename = "A5CC2571-3CEE-490B-A0FD-B7F1B0F4A642_4_5005_c.jpeg"
		} else {
			c.SaveUploadedFile(image, "./../assets/images/"+image.Filename) // 여기 다시 공부
			imagename = image.Filename
		}
		

		already := c.Request.FormValue("Id")
		cate := c.Request.FormValue("Cate")
		title := c.Request.FormValue("Title")
		body := c.Request.FormValue("Text")
		body = strings.ReplaceAll(body, `"`, `$`)
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
			query := "UPDATE " + cate + ` set title = "` + title + `" , body = "` + body + `", imagepath = "./../assets/images/` + imagename + `" where ` + id + ` = ` + already
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
			query = "INSERT INTO " + cate + ` values ( ` + strconv.Itoa(idnum+1) + `,"` + title + `","` + body + `","` + string(time.Now().Format("2006-01-02 15:04:05")) + `","./../assets/images/` + imagename + `")`
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
		}
		for row.Next() {
			row.Scan(&data.Id, &data.Title, &data.Body, &data.Written_time, &data.Imagepath)
			data.Category = cate
		}
		data.Body = strings.ReplaceAll(data.Body, "$", `"`)
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
		}

		body := string(md_body)
		html_data := HTML_DATA{
			Data_text:      template.HTML(body),
			Data_title:     data.Title,
			Data_cate:      data.Category,
			Data_id:        data.Id,
			Data_imagepath: data.Imagepath,
		}

		// 이전 포스트, 다음 포스트로 넘어가기 위해 before/after post의 id, cate 찾기
		exist := ""
		// 이전 포스트 존재 여부 확인
		query = "SELECT " + id + " FROM " + cate + " where " + id + " < " + index + " order by " + id + " desc limit 1"
		r_before, err := db.Query(query)
		if err != nil {
			log.Fatalln("BEFORE POST UPLOADING ERROR occured :", err)
		}
		if r_before.Next() {
			var temp_before Post_DATA
			r_before.Scan(&temp_before.ID)
			temp_before.Category = cate
			html_data.BeforeAndAfter = append(html_data.BeforeAndAfter, temp_before)
			exist += "b"
		}

		// 다음 포스트 존재 여부 확인
		query2 := "SELECT " + id + " FROM " + cate + " where " + id + " > " + index + " order by " + id + " asc limit 1"
		r_after, err := db.Query(query2)
		if err != nil {
			log.Fatalln("AFTER POST UPLOADING ERROR occured :", err)
		}
		if r_after.Next() { // 다음 포스트가 있으면
			var temp_after Post_DATA
			r_after.Scan(&temp_after.ID)
			temp_after.Category = cate
			html_data.BeforeAndAfter = append(html_data.BeforeAndAfter, temp_after)
			exist += "a"
		}
		// case 에 따라서 버튼에 추가할 이름 정의
		if exist == "a" {
			html_data.BeforeAndAfter[0].Where = "다음 포스트"
		} else if exist == "b" {
			html_data.BeforeAndAfter[0].Where = "이전 포스트"
		} else {
			html_data.BeforeAndAfter[0].Where = "이전 포스트"
			html_data.BeforeAndAfter[1].Where = "다음 포스트"
		}

		// 전체 게시판에서의 예외처리 헤야함

		c.HTML(http.StatusOK, "postlook.html", html_data)
	})

	// 게시글 삭제
	eg.GET("/delete", func(c *gin.Context) {
		value, err := c.Cookie("admistrator")
		if err == http.ErrNoCookie {
			c.Redirect(http.StatusSeeOther, "/loginpage")
		} else if value != "OK" {
			c.Redirect(http.StatusSeeOther, "/loginpage")
		} else {
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
			// default image 이면 지우지 않도록
			if imagepath != "./../assets/images/ACAC4F47-2E85-412C-98F0-E3750922E1D9_4_5005_c.jpeg"{
				os.Remove("." + imagepath)
			}
			col := ""
			switch id {
			case "pid" : col = "p_id"
			case "sid" : col = "s_id"
			case "rid" : col = "r_id"
			}
			query = "DELETE from whole where " + col + " = " + index // 해당 레코드 외래키 테이블에서 지우기 (참조무결성 원칙)
			_, err = db.Query(query)
			if err != nil {
				log.Fatalln("DELETE in Whole Table ERROR :", err)
			}
			query = "DELETE from " + cate + " where " + id + " = " + index // 해당 레코드 기본키 테이블에서 지우기
			_, err = db.Query(query)
			if err != nil {
				log.Fatalln("DELETE in Cate Table ERROR :", err)
			}

			c.Redirect(http.StatusSeeOther, "/postlist?listname="+cate)
		}
	})

	//게시글 수정
	eg.GET("/modify", func(c *gin.Context) {
		value, err := c.Cookie("admistrator")
		if err == http.ErrNoCookie {
			c.Redirect(http.StatusSeeOther, "/loginpage")
		} else if value != "OK" {
			c.Redirect(http.StatusSeeOther, "/loginpage")
		} else {
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
			}
			for row.Next() {
				var imagepath string
				row.Scan(&data.Id, &data.Title, &data.Body, &data.Written_time, &imagepath)
				data.Category = cate
				if strings.Contains(imagepath, "./../assets"){
					_, imagepath, _ = strings.Cut(imagepath, "./../assets/images/")
				}
				data.Imagepath = imagepath
			}
			data.Body = strings.ReplaceAll(data.Body, `$`, `"`)
			c.HTML(http.StatusOK, "modify.html", data)
		}
	})

	// 검색 기능
	eg.POST("/search", func(c *gin.Context) {
		str := c.Request.FormValue("Find")
		str = strings.ReplaceAll(str, `'`, "")
		data := []DB_DATA{}
		query := "select id, pid, title, body, datetime, imagepath from whole join Projects on whole.p_id = Projects.pid where body like '%" + str + "%' order by id desc"
		r, err := db.Query(query)
		if err != nil {
			log.Fatalln("DB Connecting ERROR occured :", err)
		}
		for r.Next() {
			var temp DB_DATA
			r.Scan(&temp.Whole_id, &temp.Id, &temp.Title, &temp.Body, &temp.Written_time, &temp.Imagepath)
			temp.Category = "Projects"
			data = append(data, temp)
		}

		query = "select id, sid, title, body, datetime, imagepath from whole join Study on whole.s_id = Study.sid where body like '%" + str + "%' order by id desc"
		r, err = db.Query(query)
		if err != nil {
			log.Fatalln("DB Connecting ERROR occured :", err)
		}
		for r.Next() {
			var temp DB_DATA
			r.Scan(&temp.Whole_id, &temp.Id, &temp.Title, &temp.Body, &temp.Written_time, &temp.Imagepath)
			temp.Category = "Study"
			data = append(data, temp)
		}

		query = "select id, rid, title, body, datetime, imagepath from whole join Review on whole.r_id = Review.rid where body like '%" + str + "%' order by id desc"
		r, err = db.Query(query)
		if err != nil {
			log.Fatalln("DB Connecting ERROR occured :", err)
		}
		for r.Next() {
			var temp DB_DATA
			r.Scan(&temp.Whole_id, &temp.Id, &temp.Title, &temp.Body, &temp.Written_time, &temp.Imagepath)
			temp.Category = "Review"
			data = append(data, temp)
		}

		// 검색결과를 최신순으로 sorting
		for i := 0; i < len(data)-1; i++ {
			for j := 0; j < len(data)-1-i; j++ {
				if data[i].Whole_id < data[i+j+1].Whole_id {
					data[i], data[i+j+1] = data[i+j+1], data[i]
				}
			}
		}

		if len(data) == 0 {
			c.HTML(http.StatusOK, "nopost.html", str)
		} else {
			c.HTML(http.StatusOK, "postlist.html", data)
		}

	})

	eg.POST("/login", func(c *gin.Context) {
		id := c.PostForm("ID")
		pw := c.PostForm("PASSWORD")
		from := c.PostForm("from")
		if id == "achoistic98" {
			if pw == "levor0805" {
				c.SetCookie("admistrator", "OK", 7200, "/", "", false, false)
				if from == "writing" {
					c.HTML(http.StatusOK, "writing.html", nil)
				} else {
					c.Redirect(http.StatusSeeOther, "/")
				}
			} else {
				c.Redirect(http.StatusSeeOther, "/")
			}
		} else {
			c.Redirect(http.StatusSeeOther, "/")
		}
	})

	// 로그인페이지로 넘어가는 핸들러
	eg.GET("/loginpage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	eg.Run(":8080")

}
