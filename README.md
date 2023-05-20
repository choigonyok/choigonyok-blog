<br>

# Choigonyok Blog

<br>
<br>

# Contents
>
        1. Outline
        2. Why?
        3. What I want to learn
        4. Technologies Used
        5. Features
        6. Details
        7. ReFactorying
        8. What I actually learned
	9. Used source


<br><br>

## Outline
>>개발 / 공부 블로그=

프로젝트의 개발과정이나 배운 것들, 또 개인적으로 학습한 내용들을 정리할 수 있는 블로그


<br><br>

## Why?
>>앞으로 유용하게 사용할 수 있어서
>      
<br>
이제 막 개발에 열정을 가지고 달리기 시작했는데, 기록하는 것이 중요하다고 생각했다
<br><br>

>>첫 프로젝트라서
>
<br>
너무 흔한 프로젝트이기 때문에 오히려 프로젝트를 개발하다 부딪히는 문제점들을 해결하기 쉬울 것이다

<br><br>

## What I want to learn
* GO 실전 적용
* 프레임워크(Gin)의 역할
* DB 설계 및 커넥션
* 배포
* HTML, CSS 기초 지식

<br><br>

## Technologies Used
* Language : <mark>GO</mark>
* BE framework : <mark>Gin</mark>
* DB : <mark>MySQL</mark>
* Deploy : <mark>AWS lightsail</mark>
* Version managing : <mark>Git</mark> / <mark>Github</mark>

<br><br>

## Features

1. 글의 CRUD가 가능해야 함
2. 글이 카테고리 별로 분류되어야 함 ( 프로젝트, 공부, 후기 등 )
3. 글을 md 형식으로 작성할 수 있어야 함
4. 클라우드에 배포해서 언제든지 접속이 가능해야함

<br><br>

## Details

<br><br>

## ReFactorying

<br><br>

## What I actually learned

<br><br>

## Used sources
* HTML templates : https://html5up.net/editorial





<br>
<br>



# 진행상황 (추후 블로그에 기록용)

>
        ## 언젠가 구현할 것
        1. 맨 위로 스크롤 올려주는 버튼 구현
        2. 게시글 휴지통 기능(삭제/복원) 구현





> ## 5/19 (약 5시간)

개발 환경 설정 (디렉토리 준비, HTML 템플릿 적용, 깃허브 레포 생성 및 연동)

개발 시작

서버와의 접속 연결

HTML - 사이드바 이름 변경

기본 html template template(?) 핸들러 구현
홈페이지로 돌아오도록 리다이렉트 일부 구현

Gin engine 생성 후 template 전부 파싱, static resource 설정

* static resource가 있는 폴더로 경로를 설정했는데, 브라우저가 폴더 중 css, js는 읽는데 images파일들을 못읽는 문제 발생




## 다음에 구현할 것

~~1. Image static resources 오류 해결~~ (5/20) -> 상대경로와 절대경로에 대한 개념 익혔음
~~2. 핸들러 추가 구현~~ (5/20) -> 리다이렉션 몇가지 구현 (글 작성 이후 해당 글 카테고리 페이지로 이동 등)

## 느낀 점

1. FE 디자인은 생각보다 더 품이 많이 든다. 시간이 너무 많이 소요됨.
2. 메인 코드는 30줄 밖에 안썼지만 아무 클론할 예제 없이 혼자 머리 굴리며 0부터 시작하니까 너무 재밌음.
3. 깃허브에 커밋하고 푸시할 때마다 개발자 된 기분(?)




> ## 5/20 (약 5시간)

사이드바 구현 완료

db스키마 정의 완료

글 작성하면 카테고리별로 DB에 저장하는 기능 구현완료

DB에 있는 글들을 읽어서 각 카테고리 별로 게시판에 정렬하는 기능 구현완료

## 다음에 구현할 것

* 게시판의 글들 클릭하면 자동으로 html 생성해서 내용 볼 수 있도록 구현하기
* 그 페이지에서 삭제하면 게시판과 DB에서도 사라지도록 삭제기능 구현하기
* 그 페이지에서 수정하면 게시판과 DB에서도 수정되도록 수정기능 구현하기

## 느낀 점

 * 코드를 짜고, 짠 코드를 당일에 리뷰까지 하는 건 너무 버거운 것 같다 아직은
 * 손목이 엄청 아프다. 괜히 개발자들이 비싼 키보드 사는게 아닌 것 같다.
 * 문제를 해결하는 그 과정은 너무 짜릿한 것 같다.
 * 클라이언트 쪽 지식이 있고 없고 차이가 프로그램의 완성도에 꽤 큰 영향을 미치는 것 같다.




# weg

* 이제 드디어
* 되나?

> 진짜 잠깐이었지만
>포기하고싶었다
>


aweg
weg

