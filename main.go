package main

import (
	"fmt"
	"net/http"

	"github.com/boj/redistore"
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var (
	db             DB
	templateSignup = pongo2.Must(pongo2.FromFile("template/signup.html"))
	templatelogin  = pongo2.Must(pongo2.FromFile("template/login.html"))
)

type DB struct {
	*sqlx.DB
}

// Handler
func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func handlerGetSingUp(c echo.Context) error {
	body, err := templateSignup.Execute(
		pongo2.Context{},
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, body)
}

func handlerPostSignUp(c echo.Context) error {
	userID := c.FormValue("userID")
	userName := c.FormValue("userName")

	password, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("password hash error:%v\n", err)
	}

	IDNumber, err := db.InsertUser(userID, userName, string(password))
	if err != nil {
		log.Printf("insert error:%v\n", err)
	}
	fmt.Printf("insert number %d\n", IDNumber)
	db.FetchUsers()
	return c.String(http.StatusOK, "Hello, World!")
}

func handlerGetLogin(c echo.Context) error {
	body, err := templatelogin.Execute(
		pongo2.Context{},
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, body)
}

func handlerPostLogin(c echo.Context) error {
	userID := c.FormValue("userID")

	password := c.FormValue("password")
	authResult := authentication(userID, password)
	fmt.Println(authResult)

	if authResult {
		return c.String(http.StatusOK, "こんにちは")
	} else {
		return c.String(http.StatusUnauthorized, "ログイン失敗")
	}

}

func authentication(userID, password string) bool {
	user, err := db.FetchUserByID(userID)
	if err != nil {
		log.Printf("userInfo fetch error:%v\n", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
		return true
	} else {
		log.Printf("authenticattion error:%v\n", err)
		return false
	}
}

func main() {
	//DB接続
	sqlxdb, err := sqlx.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		"local_user",
		"local_password",
		"localhost",
		"3306",
		"user",
	))
	if err != nil {
		log.Fatalf("DB Connection Error: %v", err)
		return
	}
	db = DB{sqlxdb}
	//ここで新規ユーザ登録
	db.FetchUsers()

	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", helloHandler)
	e.GET("/signup", handlerGetSingUp)
	e.POST("/signup", handlerPostSignUp)
	e.GET("/login", handlerGetLogin)
	e.POST("/login", handlerPostLogin)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
