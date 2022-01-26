# 카카오TV 플랫폼 채팅 봇

### 프로젝트 소개

이 프로젝트는 카카오 TV라는 라이브 스트리밍 플랫폼에서 작동할 수 있는 전용 채팅봇입니다.

트위치의 [싹둑](https://ssakdook.twip.kr/) 같은 채팅봇입니다.

Golang으로 개발되었으며 사용된 프레임워크 및 라이브러리는 다음과 같습니다.

- Echo web framework
- Websocket
- Selenium standalone
- Postgresql

폴더 구조는 다음과 같습니다.

```
├── config
│   └── database
├── domain
├── handler
├── repository
└── public
```

트위치와 다르게 카카오티비는 개발자를 위한 Developer API를 제공하지 않습니다.

모든 데이터는 [KakaoTv](https://tv.kakao.com/) 웹 클라이언트에서 나오는 값들에 맞춰 개발되었습니다.

작동 방식은 다음과 같습니다.

1. PD(이하 스트리머) 등록
2. 스트리머의 번호를 등록 후 Get Request를 통한 방송을 시작할 때까지 여부 확인
3. 방송이 시작된 경우
   1. Selenium Standalone을 통해 docker 내부에 Selenium 실행
   2. Websocket을 통해 해당 방송 채팅방에 접속 후 메시지 읽기 실행(읽는 건 로그인이 필요 없음)
4. 채팅 전용 벗이 웹에 Selenium에 설정된 아이디, 비밀번호 값으로 로그인을 시도
5. 봇이 로그인이 된 후부터 이 플랫폼에 작성된 해당 스트리머의 명령어 및 금지어들을 읽어 기능 실행
6. 해당 스트리머의 방송이 꺼질 경우 Selenium 종료. 후 다시 위의 2번으로 가 방송 여부 확인

`셀레니움 및 채팅 컨트롤 부분은 main.go:178부터 진행됩니다.`

웹에 사용된 자세한 미들웨어는 [Echo](https://echo.labstack.com/) 여기서 확인할 수 있습니다.

제공되는 API는 다음과 같습니다. (main.go에서 확인할 수 있습니다.)

```go
e.File("/", "public/index.html")
e.File("/user", "public/hello.html")
e.DELETE("/instruct", func (c echo.Context) error {})
e.DELETE("/noword", func (c echo.Context) error {})
e.POST("/user-info", func (c echo.Context) error {})
e.POST("/instruction", func (c echo.Context) error {})
e.POST("/noword", func (c echo.Context) error {})
e.GET("/noword", func (c echo.Context) error {})
e.GET("/instruction", func (c echo.Context) error {})
e.GET("/pd-list", func (c echo.Context) error {})
e.GET("/pd-user-id", func (c echo.Context) error {})
e.POST("/login", func (c echo.Context) error {})
r := e.Group("/restricted")
```

웹에서 제공되는 화면은 다음과 같습니다.

![images1](/document/images/1.png)
(메인)

![images2](/document/images/2.png)
(기능)