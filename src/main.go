package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type writing struct {
	Title        string
	Body         string
	Written_Time string
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
	//왜 images 폴더는 캐싱이 안될까?

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

	// DB에서 프로젝트 read
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
			st.Written_Time = wt
			ss = append(ss, st)
		}
		c.HTML(http.StatusOK, "generic.html", ss)
	})

	// DB에서 리뷰 read
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
			st.Written_Time = wt
			ss = append(ss, st)
		}
		c.HTML(http.StatusOK, "review.html", ss)
	})

	// DB에서 스터디 read
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
			st.Written_Time = wt
			ss = append(ss, st)
		}
		c.HTML(http.StatusOK, "study.html", ss)
	})

	// // Write한 글을 DB에 저장
	// eg.POST("/receive", func(c *gin.Context) {
	// 	p_categoty := c.PostForm("post-category")
	// 	p_title := c.PostForm("post-title")
	// 	p_body := c.PostForm("post-write")
	// 	// 클라이언트가 카테고리를 설정하지 않고 업로드했을 때
	// 	if p_categoty == "" || p_title == "" || p_body == "" {
	// 		http.Error(c.Writer, "THERE IS EMPTY BOX", http.StatusBadRequest)
	// 	}
	// 	// 클라이언트가 카테고리를 프로젝트로 설정하고 업로드했을 때
	// 	if p_categoty == "projects" {
	// 		r, err := db.Query("SELECT pid FROM projects order by pid desc limit 1")
	// 		if err != nil {
	// 			log.Fatalln("@@@DB connection error :", err, "@@@")
	// 		}
	// 		defer r.Close()
	// 		var p_id int
	// 		for r.Next() {
	// 			r.Scan(&p_id)
	// 		}
	// 		db.Query("INSERT INTO projects values (?,?,?,?)", p_id+1, p_title, p_body, string(time.Now().Format(time.DateTime)))
	// 		c.Redirect(http.StatusSeeOther, "/generic.html")
	// 	}
	// 	// 클라이언트가 카테고리를 스터디로 설정하고 업로드했을 때
	// 	if p_categoty == "study" {
	// 		r, err := db.Query("SELECT sid FROM study order by sid desc limit 1")
	// 		if err != nil {
	// 			log.Fatalln("@@@DB connection error :", err, "@@@")
	// 		}
	// 		defer r.Close()
	// 		var s_id int
	// 		for r.Next() {
	// 			r.Scan(&s_id)
	// 		}
	// 		db.Query("INSERT INTO study values (?,?,?,?)", s_id+1, p_title, p_body, string(time.Now().Format(time.DateTime)))
	// 		c.Redirect(http.StatusSeeOther, "/study.html")
	// 	}
	// 	// 클라이언트가 카테고리를 리뷰로 설정하고 업로드했을 때
	// 	if p_categoty == "review" {
	// 		r, err := db.Query("SELECT rid FROM review order by rid desc limit 1")
	// 		if err != nil {
	// 			log.Fatalln("@@@DB connection error :", err, "@@@")
	// 		}
	// 		defer r.Close()
	// 		var r_id int
	// 		for r.Next() {
	// 			r.Scan(&r_id)
	// 		}
	// 		db.Query("INSERT INTO review values (?,?,?,?)", r_id+1, p_title, p_body, string(time.Now().Format(time.DateTime)))
	// 		c.Redirect(http.StatusSeeOther, "/review.html")
	// 	}
	// })

	// Write한 글을 DB에 저장
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

		// 외부 API에 POST 요청 보내기
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
		fmt.Println("@@@@@@@@", string(body), "@@@@@@@@")

		// if err := c.ShouldBindJSON(&data); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		// 	return
		// }

		fmt.Println("json text :", data.Text)
		fmt.Println("json title :", ftt)
		fmt.Println("json cate :", fc)

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

			fmt.Println(string(body))
			data := struct {
				TagHTML template.HTML
			}{
				TagHTML: template.HTML(string(body)),
			}
			c.HTML(http.StatusOK, "project1.html", data)

			//escapedHTML := template.HTMLEscapeString(string(body))

			// c.Redirect(http.StatusSeeOther, "/generic.html")
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
			c.Redirect(http.StatusSeeOther, "/study.html")
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
			c.Redirect(http.StatusSeeOther, "/review.html")
		}
	})

	// text := c.PostForm(`post-write`)
	// apiURL := "https://api.github.com/markdown"
	// payload := fmt.Sprintf(`{"text": "%s"}`, text)
	// fmt.Println(payload)
	// req, err := http.NewRequest("POST", apiURL, strings.NewReader(payload))
	// if err != nil {
	// 	fmt.Println("Failed to create request:", err)
	// 	return
	// } //strings.NewReader(payload)

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Accept", "text/html")
	// req.Header.Set("Authorization", "Bearer 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")
	// req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// client := http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println("API request failed:", err)
	// 	return
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Println("API request failed with status:", resp.StatusCode)
	// 	return
	// }

	// respBody, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Json EDEEEERROR3333")
	// 	return
	// }
	// fmt.Println("HTML Response:")
	// fmt.Println(string(respBody))

	// c.HTML(http.StatusOK, "project1.html", string(respBody))

	// decodedResult := html.UnescapeString(response.Result)
	// c.HTML(http.StatusOK, "project1.html", decodedResult)
	//	result := c.GetString("result")

	// c.HTML(http.StatusOK, "project1.html", result)

	// 	p_categoty := c.PostForm("post-category")
	// 	p_title := c.PostForm("post-title")
	// 	p_body := c.PostForm("post-write")
	// 	// 클라이언트가 카테고리를 설정하지 않고 업로드했을 때
	// 	if p_categoty == "" || p_title == "" || p_body == "" {
	// 		http.Error(c.Writer, "THERE IS EMPTY BOX", http.StatusBadRequest)
	// 	}
	// 	// 클라이언트가 카테고리를 프로젝트로 설정하고 업로드했을 때
	// 	if p_categoty == "projects" {
	// 		r, err := db.Query("SELECT pid FROM projects order by pid desc limit 1")
	// 		if err != nil {
	// 			log.Fatalln("@@@DB connection error :", err, "@@@")
	// 		}
	// 		defer r.Close()
	// 		var p_id int
	// 		for r.Next() {
	// 			r.Scan(&p_id)
	// 		}
	// 		db.Query("INSERT INTO projects values (?,?,?,?)", p_id+1, p_title, resp.Body, string(time.Now().Format(time.DateTime)))
	// 		c.Redirect(http.StatusSeeOther, "/generic.html")
	// 	}
	// 	// 클라이언트가 카테고리를 스터디로 설정하고 업로드했을 때
	// 	if p_categoty == "study" {
	// 		r, err := db.Query("SELECT sid FROM study order by sid desc limit 1")
	// 		if err != nil {
	// 			log.Fatalln("@@@DB connection error :", err, "@@@")
	// 		}
	// 		defer r.Close()
	// 		var s_id int
	// 		for r.Next() {
	// 			r.Scan(&s_id)
	// 		}
	// 		db.Query("INSERT INTO study values (?,?,?,?)", s_id+1, p_title, resp.Body, string(time.Now().Format(time.DateTime)))
	// 		c.Redirect(http.StatusSeeOther, "/study.html")
	// 	}
	// 	// 클라이언트가 카테고리를 리뷰로 설정하고 업로드했을 때
	// 	if p_categoty == "review" {
	// 		r, err := db.Query("SELECT rid FROM review order by rid desc limit 1")
	// 		if err != nil {
	// 			log.Fatalln("@@@DB connection error :", err, "@@@")
	// 		}
	// 		defer r.Close()
	// 		var r_id int
	// 		for r.Next() {
	// 			r.Scan(&r_id)
	// 		}
	// 		db.Query("INSERT INTO review values (?,?,?,?)", r_id+1, p_title, resp.Body, string(time.Now().Format(time.DateTime)))
	// 		c.Redirect(http.StatusSeeOther, "/review.html")
	// 	}
	// })

	// eg.POST("/markdown", func(c *gin.Context) {
	// 	text := c.PostForm("text")
	// 	apiURL := "https://api.github.com/markdown"
	// 	payload := []byte(fmt.Sprintf(`{"text": "%s"}`, text))
	// 	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	req.Header.Set("Accept", "application/vnd.github+json")
	// 	req.Header.Set("Authorization", "Bearer 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")
	// 	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// 	client := http.Client{}
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	defer resp.Body.Close()

	// 	respBody, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	c.HTML(http.StatusOK, "project1.html", respBody)

	// 	// 응답 반환
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Success",
	// 		"result":  string(respBody),
	// 	})
	// })

	eg.Run(":8080")

}
