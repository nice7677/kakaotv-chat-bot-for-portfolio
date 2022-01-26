package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tebeka/selenium"
	"kakaotv-chat-bot/domain"
	"kakaotv-chat-bot/handler"
	"log"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func main() {
	//db := database.Connect()
	//database.Close(db)
	go func() {
		for {
			StatusCheck = broadcastStatusChecker()
			if StatusCheck == true && FirstTimeCheck == 0 {
				FirstTimeCheck = 1
				getLiveLinkId()
				go func() {
					Start()
				}()
				time.Sleep(5 * time.Second)
				getChatContent()
			}
			time.Sleep(10 * time.Second)
		}
	}()
	e := echo.New()
	e.File("/", "public/index.html")
	e.File("/user", "public/hello.html")
	e.DELETE("/instruct", func(c echo.Context) error {
		function := c.FormValue("function")
		allHandler := handler.AllHandler{}
		instVM := &domain.InstructionVM{
			Function: function,
		}
		allHandler.DeleteInstruct(instVM)
		return c.JSON(http.StatusOK, "삭제완료")
	})
	e.DELETE("/noword", func(c echo.Context) error {
		word := c.FormValue("word")
		allHandler := handler.AllHandler{}
		noVM := &domain.NoWordVM{
			Word: word,
		}
		allHandler.DeleteNoword(noVM)
		return c.JSON(http.StatusOK, "삭제완료")
	})
	e.POST("/user-info", func(c echo.Context) error {
		userId := c.FormValue("userid")
		lolid := c.FormValue("lolid")
		UserId = userId
		lolId = lolid
		return c.JSON(http.StatusOK, "등록완료")
	})
	e.POST("/instruction", func(c echo.Context) error {
		function := c.FormValue("function")
		dab := c.FormValue("dab")

		instructVM := &domain.InstructionVM{
			Idx:      0,
			Word:     dab,
			Function: function,
		}

		allHandler := handler.AllHandler{}
		allHandler.SaveInstruction(instructVM)

		return c.JSON(http.StatusOK, "등록완료")
	})
	e.POST("/noword", func(c echo.Context) error {
		noword := c.FormValue("noword")
		pdUserID := c.FormValue("pd-user-id")

		noWordVM := &domain.NoWordVM{
			Idx:      0,
			Word:     noword,
			PDUserID: pdUserID,
		}

		allHandler := handler.AllHandler{}
		allHandler.SaveNoword(noWordVM)

		return c.JSON(http.StatusOK, "등록완료")
	})
	e.GET("/noword", func(c echo.Context) error {
		allHandler := handler.AllHandler{}

		return c.JSON(http.StatusOK, allHandler.GetNoword())
	})
	e.GET("/instruction", func(c echo.Context) error {
		allHandler := handler.AllHandler{}

		return c.JSON(http.StatusOK, allHandler.GetInstruction())
	})
	e.GET("/pd-list", func(c echo.Context) error {
		allHandler := handler.AllHandler{}

		return c.JSON(http.StatusOK, allHandler.GetPD())
	})
	e.GET("/pd-user-id", func(c echo.Context) error {
		//allHandler := handler.AllHandler{}

		return c.JSON(http.StatusOK, PDUSERID)
	})
	e.POST("/login", func(c echo.Context) error {
		id := c.FormValue("id")
		pw := c.FormValue("pw")

		// Throws unauthorized error
		if id != "jinwoo" || pw != "123!" {
			return echo.ErrUnauthorized
		}

		// Set custom claims
		claims := &jwtCustomClaims{
			"진우챗봇",
			true,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	})
	r := e.Group("/restricted")
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("", restricted)
	e.Logger.Fatal(e.Start(":80"))
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

/// ---------------------------
// 밑으로는  셀레니움 부분 및 채팅 보내기 및 관리자 권한 부분

//관리자 권한 asmg

// 이성진 그룹아이디 3614183
// 김치민 그룹아이디 3615649

var (
	StartTime int64 = 0

	LiveLinkId = ""

	UserId = ""

	Groupid   = ""
	RoomId    = ""
	SessionId = "3831298"

	Kadu    = ""
	Karmt   = ""
	Karmtea = ""
	Kawlt   = ""
	Kawltea = ""
	Klimtc  = ""
	TIARA   = ""

	Ip   = ""
	Port = ""

	// false 방송안함 true 방송중
	StatusCheck = false

	lolId = "hide on bush"

	Help = "/티어 , /업타임"
	//Help = "/티어 , /업타임, /금지어"

	wdSessionId    = ""
	FirstTimeCheck = 0

	TargetsId = ""

	PDUSERID = ""
)
var warningCount = make(map[string]int)
var helpCheck = make(map[string]string)   // help
var tierCheck = make(map[string]string)   // 티어
var uptimeCheck = make(map[string]string) // 업타임

/**
라이브 링크 아이디 구하기
*/
func getLiveLinkId() {

	c := colly.NewCollector()

	c.OnHTML(".link_contents", func(e *colly.HTMLElement) {
		if len(strings.Split(e.Attr("href"), "/channel/"+UserId+"/livelink/")) == 2 {
			liveLinkIds := strings.Split(e.Attr("href"), "/channel/"+UserId+"/livelink/")
			LiveLinkId = liveLinkIds[1]
			getGroupId(LiveLinkId)
		}
	})

	c.Visit("https://tv.kakao.com/channel/" + UserId)

}

func broadcastStatusChecker() bool {
	var statusChecker bool

	c := colly.NewCollector()

	c.OnHTML(".txt_none", func(e *colly.HTMLElement) {
		log.Println(e.Text)
		if e.Text == "콘텐츠가 없어요.." {
			//log.Println("방송안함")
			statusChecker = false
		}
	})

	c.OnHTML(".tit_vod", func(e *colly.HTMLElement) {
		//log.Println(e.Text)
		if e.Text != "" {
			//log.Println("방송중")
			statusChecker = true
		}
	})

	c.Visit("https://tv.kakao.com/channel/" + UserId)

	return statusChecker
}

/**
그룹아이디 구하기
*/
func getGroupId(liveLinkId string) {

	c := colly.NewCollector()

	c.OnHTML("body", func(e *colly.HTMLElement) {
		html, _ := e.DOM.Html()
		groupIdFirst := strings.Split(html, "groupid: '")
		groupIdSecond := strings.Split(groupIdFirst[1], "'")
		Groupid = groupIdSecond[0]
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", "https://tv.kakao.com/channel/2866299/livelink/7018784?metaObjectType=Channel")
	})

	c.Visit("https://live-tv.kakao.com/kakaotv/live/chat/user/" + liveLinkId)

}

func getStartTime() int64 {
	client := resty.New()
	resp, err := client.R().
		//SetHeader("Content-Type", "application/json").
		SetHeaders(map[string]string{
			//"Content-Type":   "application/x-www-form-urlencoded",
			//"Origin":         "https://live-tv.kakao.com",
			"Referer": "https://live-tv.kakao.com/kakaotv/live/chat/user/" + LiveLinkId, // 유저마다 번호가다름 유저번호
			//"Sec-Fetch-Mode": "cors",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36",
		}).
		SetFormData(map[string]string{
			"groupid": Groupid,
		}).
		Get("https://t1.daumcdn.net/potplayer/chat/prohibit_words.json")
	if err != nil {
		log.Println(err)
	}

	byt := []byte(resp.Body())
	var jsonMap map[string]interface{}

	if err := json.Unmarshal(byt, &jsonMap); err != nil {
		panic(err)
	}

	startTime := int64(jsonMap["ts"].(float64) * 0.001)
	//log.Println(startTime)

	return startTime

}

func getChatContent() {
	client := resty.New()
	resp, err := client.R().
		//SetHeader("Content-Type", "application/json").
		SetHeaders(map[string]string{
			"Content-Type":   "application/x-www-form-urlencoded",
			"Origin":         "https://live-tv.kakao.com",
			"Referer":        "https://live-tv.kakao.com/kakaotv/live/chat/user/" + LiveLinkId, // 유저마다 번호가다름 유저번호
			"Sec-Fetch-Mode": "cors",
			"User-Agent":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36",
		}).
		SetFormData(map[string]string{
			"groupid": Groupid,
		}).
		Post("https://play.kakao.com/chat/service/api/room")
	if err != nil {
		log.Println(err)
	}

	byt := []byte(resp.Body())
	var jsonMap map[string]interface{}
	strInfo := string(resp.Body())

	if err := json.Unmarshal(byt, &jsonMap); err != nil {
		panic(err)
	}
	enter := jsonMap["enter"]
	testMap := jsonMap["roomInfo"]

	roomInfoMapStr := fmt.Sprintf("%v", testMap)
	ipFirstSplit := strings.Split(roomInfoMapStr, "pos:")
	ipSecondSplit := strings.Split(ipFirstSplit[1], ":")
	roomIdss := strings.Split(roomInfoMapStr, "roomid:")

	Ip = ipSecondSplit[0]
	Port = ipSecondSplit[1]
	RoomId = strings.Split(roomIdss[1], " serverip")[0]
	//StartTime = getStartTime()

	startTimeFisrt := strings.Split(strInfo, `"time":`)[1]
	startTimeSecond := strings.Split(startTimeFisrt, ",")[0]
	startTime3rd, _ := strconv.ParseFloat(startTimeSecond, 64)
	startTime4th := int64(startTime3rd * 0.001)
	StartTime = startTime4th
	log.Println(StartTime)

	s := fmt.Sprintf("ENTER %s\n", enter)

	conn, err := net.Dial("tcp", Ip+":"+Port)
	if nil != err {
		log.Println(err)
	}

	what, error := conn.Write([]byte(s))
	log.Println(what)
	if error != nil {
		log.Println(error.Error())
	}

	data := make([]byte, 40960)

	//go func() {

	go func() {
		for {
			StatusCheck = broadcastStatusChecker()
			if StatusCheck == false {
				selenium.DeleteSession(fmt.Sprintf("http://localhost:%d/wd/hub", 5556), wdSessionId)
				//selenium.DeleteSession(fmt.Sprintf("http://www.kakaotv.xyz:%d/wd/hub", 5556), wdSessionId)
				FirstTimeCheck = 0
				break
			}
			time.Sleep(10 * time.Second)
		}
	}()

	for {
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			restart()
			break
			//return
		}

		msg := string(data[:n])
		//fmt.Println(msg)
		log.Println(msg)
		getMsgCheck(msg)

		if StatusCheck == false {
			break
		}

	}
	//}()
}

func keepDoingSomething() (bool, error) {
	for {
		go sendMsg("카카오 챗봇입니다. 테스트 중 입니다. 매너채팅 [ 도움말-> /명령어 ]")
		//go sendMsg("금지어 2회 경고 누적시 자동 채금 입니다. 테스트 중 입니다. 매너채팅 [ 도움말-> /금지어 ]")
		//sendMsg("테스트 중 입니다.")
		time.Sleep(300 * time.Second)
	}
}

func restart() {
	getChatContent()
}

func getMsgCheck(msg string) {

	checkLen := strings.Split(msg, "NORMAL ")

	//checkFindLen := strings.Split(msg, " JOIN")

	/*if len(checkFindLen) == 2 {
		getIdOneNumber := strings.Split(checkFindLen[0], ":")[1]
		log.Println(getIdOneNumber)
		TargetsId = getIdOneNumber
		perm5(getIdOneNumber)
	}*/

	if len(checkLen) == 2 {

		getIdOne := strings.Split(msg, " ALL")
		//getIdTwo := strings.Split(getIdOne[0], "MSG ")[1]

		getNumberId := strings.Split(strings.Split(getIdOne[0], " MSG ")[0], ":")[1]
		TargetsId = getNumberId
		//log.Println(getNumberId)
		//realMSG := strings.Split(msg, "NORMAL ")[1]
		typeCheck := fmt.Sprint(reflect.TypeOf(checkLen[1]))
		var jsonMap map[string]string

		if err := json.Unmarshal([]byte(checkLen[1]), &jsonMap); err != nil {
			log.Println("pd or ad의 개소리")
		}

		if typeCheck == "string" && jsonMap["font"] == "" {

			allHandler := handler.AllHandler{}
			instructList := allHandler.GetInstruction()
			nowordList := allHandler.GetNoword()

			mapMSG := jsonMap["msg"]
			//splitMSG := strings.Split(mapMSG, m_fail)

			fmt.Println(mapMSG)

			for _, element := range *instructList {

				if mapMSG == "/"+element.Function {
					go sendMsg(element.Word)
				}

			}

			// 채금
			/*for _, element := range *nowordList {

				r, _ := regexp.Compile(element.Word)
				match := r.MatchString(mapMSG)
				//fmt.Println(match)
				log.Println(warningCount[TargetsId])
				if match == true && warningCount[TargetsId] == 0 {
					warningCount[TargetsId] = 1
					go sendMsg(getIdTwo + " 금지어 사용 경고 1회 경고 2회 누적시 채금")
				} else if match == true && warningCount[TargetsId] == 1 {
					warningCount[TargetsId] = 2
					perm5(getIdTwo)
				} else if match == true && warningCount[TargetsId] == 2 {
					warningCount[TargetsId] = 1
					go sendMsg(getIdTwo + " 금지어 사용 경고 1회 경고 2회 누적시 채금")
				}

			}*/

			//if mapMSG == "채금 테스트" {
			//	perm5(getIdTwo)
			//}

			if mapMSG == "/명령어" {

				var test string
				for _, element := range *instructList {
					test += " , /" + element.Function
				}
				go sendMsg(Help + test)

			}

			if mapMSG == "/티어" {

				//check := helpCheck[getIdTwo]
				//if check == "" {
				//	tierCheck[getIdTwo] = "nono"
				info := getUserInfoUseFowKR()
				var statusMSGTEXT = lolId + "," + info.Grade + "," + info.Point
				go sendMsg(statusMSGTEXT)
				//} else if check == "nono" {

				//}

			}

			//if mapMSG == "/전판" {
			//	afterOneGame := getUserInfoUseFowKRStatus()
			//	go sendMsg(afterOneGame)
			//}

			if mapMSG == "/업타임" {

				//check := helpCheck[getIdTwo]
				//if check == "" {
				//	uptimeCheck[getIdTwo] = "nono"
				upTime := getUpTimeDate()
				go sendMsg(upTime)
				//} else if check == "nono" {
				//
				//}

			}

			if mapMSG == "/금지어" {

				var test = ""
				for _, element := range *nowordList {
					test += element.Word + " , "
				}

				go sendMsg(test)

			}

			//if mapMSG == sibal {
			//	go sendMsg("바른말 고운말을 사용합시다.")
			//} else if mapMSG == ssibal {
			//	go sendMsg("바른말 고운말을 사용합시다.")
			//} else if mapMSG == qt {
			//	go sendMsg("바른말 고운말을 사용합시다.")
			//}

		}

	}

}

// 메세지 보내기
func sendMsg(sendMSGTEXT string) {

	client := resty.New()

	var cookies []*http.Cookie

	cookies = append(cookies, &http.Cookie{
		Name:     "_kadu",
		Value:    Kadu,
		Path:     "/",
		Domain:   ".kakao.com",
		HttpOnly: true,
	})
	cookies = append(cookies, &http.Cookie{
		Name:     "_karmt",
		Value:    Karmt,
		Path:     "/",
		Domain:   ".kakao.com",
		HttpOnly: true,
		Secure:   true,
	})
	cookies = append(cookies, &http.Cookie{
		Name:     "_karmtea",
		Value:    Karmtea,
		Path:     "/",
		Domain:   ".kakao.com",
		HttpOnly: true,
		Secure:   true,
	})
	cookies = append(cookies, &http.Cookie{
		Name:     "_kawlt",
		Value:    Kawlt,
		Path:     "/",
		Domain:   ".kakao.com",
		HttpOnly: true,
		Secure:   true,
	})
	cookies = append(cookies, &http.Cookie{
		Name:     "_kawltea",
		Value:    Kawltea,
		Path:     "/",
		Domain:   ".kakao.com",
		HttpOnly: true,
		Secure:   true,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "_klimtc",
		Value:  Klimtc,
		Path:   "/",
		Domain: ".kakao.com",
	})

	client.SetCookies(cookies)

	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":   "application/x-www-form-urlencoded",
			"Origin":         "https://live-tv.kakao.com",
			"Referer":        "https://live-tv.kakao.com/kakaotv/live/chat/user/" + LiveLinkId,
			"Sec-Fetch-Mode": "cors",
			//"User-Agent":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36",
		}).
		SetFormData(map[string]string{
			"sessionid": SessionId,
			"roomid":    RoomId,
			"msg":       `{"msg":"` + sendMSGTEXT + `"}`,
		}).
		Post("https://play.kakao.com/chat/service/api/msg")

	if err != nil {
		log.Println(err.Error())
	}
	log.Println(resp)

}

type FowKRModel struct {
	Rank  string `json:"rank"`
	Grade string `json:"grade"`
	Point string `json:"point"`
	//PromotionMode string
	Record string `json:"record"`
}

// 전적검색
func getUserInfoUseFowKR() *FowKRModel {

	fowKrModel := &FowKRModel{}

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".table_summary", func(e *colly.HTMLElement) {
		text := e.DOM.Children().Find("div").Text()
		splitInfo := strings.Split(text, ":")
		re1, _ := regexp.Compile("\n")
		re2, _ := regexp.Compile("\t")
		rank1 := re1.ReplaceAllString(strings.Split(splitInfo[1], "리그")[0], "")
		rank2 := re2.ReplaceAllString(rank1, "")
		fowKrModel.Rank = "랭킹 :" + rank2
		grade1 := re1.ReplaceAllString(strings.Split(splitInfo[3], "리그")[0], "")
		grade2 := re2.ReplaceAllString(grade1, "")
		fowKrModel.Grade = grade2
		point1 := re1.ReplaceAllString(strings.Split(splitInfo[4], "승급전")[0], "")
		point2 := re2.ReplaceAllString(point1, "")
		fowKrModel.Point = point2 + "점"
		//fowKrModel.PromotionMode = "승급전"
		record1 := re1.ReplaceAllString(strings.Split(strings.Split(splitInfo[5], "-")[1], "리그")[0], "")
		record2 := re2.ReplaceAllString(record1, "")
		fowKrModel.Record = "전적 :" + record2
	})

	var url string

	url = "http://fow.kr/find/" + lolId

	c.Visit(url)

	return fowKrModel
}

// 전적검색
func getUserInfoUseFowKRStatus() string {

	var returnValue string
	var firstInt int = 0

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".table_recent tbody tr", func(e *colly.HTMLElement) {
		if firstInt == 0 {
			firstInt = 1
			returnValue = e.DOM.Children().Text()
			re1, _ := regexp.Compile("\n\t\n\t\n\t\t")
			s1 := re1.ReplaceAllString(returnValue, " / ")
			re2, _ := regexp.Compile("\n\t\n\n")
			returnValue = re2.ReplaceAllString(s1, " / ")
			returnValue = strings.Split(returnValue, "PLAY")[0]
			log.Println(returnValue)
		}
	})

	var url string

	url = "http://fow.kr/find/" + lolId

	c.Visit(url)

	return returnValue
}

// 셀레니움 시작
func Start() {

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 5556))
	//wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://www.kakaotv.xyz:%d/wd/hub", 5556))
	if err != nil {
		log.Println(err)
		//panic(err)
	}
	wdSessionId = wd.SessionID()
	defer wd.Quit()

	// Navigate to the simple playground interface.
	if err := wd.Get("https://accounts.kakao.com/login?continue=https%3A%2F%2Ftv.kakao.com%2Fchannel%2F" + UserId + "%2Flivelink%2F" + LiveLinkId + "%3FmetaObjectType%3DChannel"); err != nil {
		log.Println(err.Error())
	}

	idInputElem, err := wd.FindElement(selenium.ByXPATH, "//input[@id='id_email_2']")
	if err != nil {
		log.Println(err.Error())
	}

	idInputElem.SendKeys("nice7677@naver.com")

	pwdInputElem, err := wd.FindElement(selenium.ByID, "id_password_3")
	if err != nil {
		log.Println(err.Error())
	}

	pwdInputElem.SendKeys("wlsdn123")

	loginBtnElem, err := wd.FindElement(selenium.ByXPATH, "//button[@class='btn_g btn_confirm submit']")
	if err != nil {
		log.Println(err.Error())
	}

	loginBtnElem.Click()

	time.Sleep(5 * time.Second)

	cookies, _ := wd.GetCookies()

	//go func() {
	//	getChatContent()
	//}()

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "_kadu" {
			Kadu = cookies[i].Value
		} else if cookies[i].Name == "_karmt" {
			Karmt = cookies[i].Value
		} else if cookies[i].Name == "_karmtea" {
			Karmtea = cookies[i].Value
		} else if cookies[i].Name == "_kawlt" {
			Kawlt = cookies[i].Value
		} else if cookies[i].Name == "_kawltea" {
			Kawltea = cookies[i].Value
		} else if cookies[i].Name == "_klimtc" {
			Klimtc = cookies[i].Value
		} else if cookies[i].Name == "TIARA" {
			TIARA = cookies[i].Value
		}
	}

	go keepDoingSomething()

	log.Println(cookies)

	for {
		time.Sleep(100 * time.Second)
		//log.Println("test")
		//if StatusCheck == false {
		//}
	}

}

// 업타임
func getUpTimeDate() string {
	times := time.Now().Unix() - StartTime
	fullMin := times / 60
	//fullHour := fullMin / 60
	hour := strconv.FormatInt(fullMin/60, 10)
	min := strconv.FormatInt(fullMin%60, 10)
	second := strconv.FormatInt(times%60, 10)
	upTime := hour + "시간" + min + "분" + second + "초"
	return upTime
}

// 채금
func perm5(id string) {
	client := resty.New()
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Host":         "play.kakao.com",
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36",
			"Cookie":       "_kadu=" + Kadu + "; _kawlt=" + Kawlt + "; _kawltea=" + Kawltea + "; _karmt=" + Karmt + "; _karmtea=" + Karmtea + ";",
		}).
		SetFormData(map[string]string{
			"roomid":    RoomId,
			"targetsid": TargetsId,
			"perm":      "PERM_NOCHAT",
			"eternal":   "0",
		}).
		Post("https://play.kakao.com/chat/service/api/perm")
	if err != nil {
		log.Println(err)
	}

	byt := []byte(resp.Body())
	var jsonMap map[string]interface{}

	if err := json.Unmarshal(byt, &jsonMap); err != nil {
		panic(err)
	}
	//log.Println(" 채금 어떻게?", jsonMap)
	go sendMsg(id + " 금지어 사용으로 인한 채금!")
}
