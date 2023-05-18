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


5/19

개발 환경 설정 (디렉토리 준비, HTML 템플릿 적용, 깃허브 레포 생성 및 연동)

개발 시작

서버와의 접속 연결

HTML - 사이드바 이름 변경

기본 html template template(?) 핸들러 구현
홈페이지로 돌아오도록 리다이렉트 일부 구현

Gin engine 생성 후 template 전부 파싱, static resource 설정

* static resource가 있는 폴더로 경로를 설정했는데, 브라우저가 폴더 중 css, js는 읽는데 images파일들을 못읽는 문제 발생


## 다음에 구현할 것

1. Image static resources 오류 해결
2. 핸들러 추가 구현

## 느낀 점

1. FE 디자인은 생각보다 더 품이 많이 든다. 시간이 너무 많이 소요됨.
2. 메인 코드는 30줄 밖에 안썼지만 아무 클론할 예제 없이 혼자 머리 굴리며 0부터 시작하니까 너무 재밌음.
3. 깃허브에 커밋하고 푸시할 때마다 개발자 된 기분(?)