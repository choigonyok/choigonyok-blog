<br>

## **Project #1 - 2 : 개발환경 구성하기**
---
<br>

>* ## 깃허브 레포지토리

        이번 프로젝트를 통해서 깃허브를 처음 써보았다.
<br>

계정은 이전에 만들어두었으나, '깃허브엔 레포라는 게 있다' 정도로만 알고있었고, 개념조차도 모르고 있었다.
<br>

프로젝트의 버전관리라는 명목의 사실상 백업을 위해 깃허브 레포지토리를 새롭게 하나 만들고
프로젝트 폴더와 연동했다.
<br><br>


>* ## HTML 템플릿 다운로드

        네가 알던 내가 아냐 (feat.template)
<br>

템플릿이라고 하면 뭔가 일이 아주 간단해질 것 같은 느낌이 든다. 
<br><br>

### **나만 그런가...?**

<br>

네이버에서 블로그 만들기를 누르면 사이드바 메뉴 이름은 뭘로 바꿀지만 정하면 되고, 몇 가지 디자인이나 글을 쓰고 지우는 것만 하면 된다.
<br>

그러나 HTML템플릿이 있다고 이런 것들이 가능해지는 건 아니다.
<br>

HTML템플릿은 그저 네모를 좀 굵은 네모로 바꿔주고, 색이 있는 네모로 만들어주는 정도? 
<br>

나머지는 다 "**ON MY OWN**" 이다.

그래도 보기 좋은 떡이 먹기도 좋고, 이왕이면 다홍치마라고
<br>

이런 너지만,,, 이거라도 활용해보기로 했다.

HTML 무료 템플릿 사이트에서 적당히 보기좋은 3페이지 짜리 HTML파일을 다운로드하고 프로젝트 폴더에 넣어주었다.
<br><br>


>* ## 인프라 세팅

        로컬서버 접속 / MySQL DB 연동
<br>

* 로컬서버 접속
<br>

로컬 서버에는 eg.RUN(":8080")으로 연결해준다.
GO의 template 패키지의 ListenAndServe()와 같은 역할을 한다.
<br>

```
func main(){
        eg := gin.Default()
        eg.Run(":8080")
}
```
<br>

eg는 engine pointer type으로 Gin의 기능들을 사용할 수 있게 해주는 컴포넌트이다.
<br>

eg은 
```
gin.Default()
```
혹은

```
gin.New()
```

로 만들 수 있는데, Default 는 literally 하게, 여러 미들웨어가 디폴트로 들어가있는 것이다.
<br>

미들웨어란 핸들러와 클라이언트 사이에서 여러가지 기능 (인증, 보안 등) 을 제공하는 일종의 코드 묶음 같은 것인데, 추후에 더 공부해볼 예정이다.
<br>

<br>

* MySQL DB 연동
<br>

```
db, err := sql.Open("mysql", "root:/*PASSWORD*/@/blog")
	if err != nil {
		log.Fatalln("DB IS NOT CONNECTED")
	}
	defer db.Close()
```
<br>

미리 MySQL DB에 새로운 Blog DB를 create 해뒀고, 스키마는 아직 정의되지 않은 상태이다.
<br>

코드 상 MySQL PW는 보안을 위해 주석처리 했다.
<br>
<br>

>* ## 개발 시작

        template 파싱, 사이드바 이름 변경, 기본 핸들러 구현, 리다이렉트 일부 구현, static resource 설정
<br>        

템플릿 파싱해주고
```
eg.LoadHTMLGlob("./templates/**/*.html")
```
<br>

HTML파일에서 사이드바 이름 변경해주었다.
<br>

```
// redirect to homepage
	eg.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/")
	})
```
<br>

등의 기본적인 핸들러와 메인페이지로의 리다이렉트까지 구현했는데, 
<br><br>


## <mark>**문제 발생**</mark>
<br>

### **Static resource 설정에서 문제가 생겼다.**
<br>

```
eg.Static("/assets", "yunsuk/choigonyok-blog/assets")
```

정적 파일 리소스를 브라우저로 전송하기 위해 경로를 설정했는데, '이런 경로를 찾을 수 없다'는 오류가 발생했다.
<br>

한참을 고민하고 구글링하다가 답을 찾았다.
<br>

## **경로에는 절대경로 / 상대경로가 있다**
<br>

코드를 작성하고 있는 main.go 의 경로를 살펴보자
<br>
절대경로는 
> /Users/yunsuk/choigonyok-blog/src/main.go
<br>

처럼 경로상의 Root("/")를 기반으로 찾는 절대적인 경로이고,
<br>

상대경로는
> ./src/main.go
<br>

처럼 내 위치 (여기서는 프로젝트 폴더 경로)를 기반으로 상대적인 경로를 찾는 방식이다. 여기서 **.** 은 현재 파일의 위치를 의미한다.

이 이후에 코드를
```
eg.Static("/assets", "./assets")
```

수정해줬고, 메인페이지에 사용되기 위한 내 프로필 사진이 정적 리소스 파일로 브라우저에 잘 전될되게 되었다.
<br>

## 느낀점

### 1. FE 디자인은 생각보다 더 품이 많이 드는 작업이다. 
<br>
-> CSS와 HTML, JavaScrpit 의 공부 필요성을 느꼈다.
<br><br>

### 2. 깃허브에 커밋하고 푸시할 때마다 개발자 된 기분(?)이다
<br>
-> 즐거워서 좋다. 그러나 중요한 것은 열정보다는 끈기! 열정도 좋지만 기분에 방해받지 않는 꾸준함이 필요하다.    
<br><br>

## 다음에 구현할 것

* DB schema 정의하기
* 게시물 작성 시 DB에 저장(CREATE) 시키기
* DB를 통해, 저장된 게시글 리스트를 게시판 별로 볼 수 있는 기능 구현하기

