package main

import (
	"fmt"
	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var (
	db             DB
	templateSignup = pongo2.Must(pongo2.FromFile("template/signup.html"))
	templatelogin  = pongo2.Must(pongo2.FromFile("template/login.html"))
	templateIndex  = pongo2.Must(pongo2.FromFile("template/index.html"))
)

type DB struct {
	*sqlx.DB
}

// Handler

func handlerIndex(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		log.Printf("session error:%v\n", err)
	}

	body, err := templateIndex.Execute(
		pongo2.Context{
			"userID": sess.Values["userID"],
		},
	)
	if err != nil {
		log.Printf("pongo2 error:%v\n", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, body)
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

	if authResult {
		sess, err := session.Get("session", c)
		if err != nil {
			log.Printf("session get error:%v\n", err)
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		sess.Values["userID"] = userID
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			log.Printf("session save error:%v\n", err)
		}
		return c.Redirect(http.StatusFound, "/")
	} else {
		body, err := templatelogin.Execute(
			pongo2.Context{
				"flash":  "ログイン失敗",
				"userID": userID,
			},
		)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return c.HTML(http.StatusUnauthorized, body)
	}

}

func handlerLogout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		log.Printf("session error:%v\n", err)
	}
	sess.Options = &sessions.Options{
		MaxAge: -1,
		Path:   "/",
	}
	sess.Save(c.Request(), c.Response())
	return c.String(http.StatusOK, "Hello, World!")
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
		"site",
	))
	if err != nil {
		log.Fatalf("DB Connection Error: %v", err)
		return
	}
	db = DB{sqlxdb}

	store, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte("secret-key"))
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// Echo instance
	e := echo.New()
	e.Use(session.Middleware(store))
	// Routes
	e.GET("/", handlerIndex)
	e.GET("/signup", handlerGetSingUp)
	e.POST("/signup", handlerPostSignUp)
	e.GET("/login", handlerGetLogin)
	e.POST("/login", handlerPostLogin)
	e.DELETE("/logout", handlerLogout)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
