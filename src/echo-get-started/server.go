package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	// インスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = t

	// ルートを設定
	e.GET("/", hello) // ローカル環境の場合、http://localhost:1323/ にGETアクセスされるとhelloハンドラーを実行する

	e.GET("/page2", func(c echo.Context) error {
		data := struct {
		}{}
		return c.Render(http.StatusOK, "page2", data)
	})
	e.POST("/page3", func(c echo.Context) error {
		mail := c.Request().PostFormValue("mail")
		//password := c.Request().PostFormValue("password")
		// [ユーザ名]:[パスワード]@tcp([ホスト名]:[ポート番号])/[データベース名]?charset=[文字コード]
		dbconf := define.getConnect()

		db, err := sql.Open("mysql", dbconf)

		// 接続が終了したらクローズする
		defer db.Close()

		if err != nil {
			fmt.Println(err.Error())
		}

		err = db.Ping()

		if err != nil {
			mail = "データベース接続失敗"
		} else {
			mail = "データベース接続成功"
		}

		data := struct {
			Mail string
		}{
			Mail: mail,
		}
		return c.Render(http.StatusOK, "page3", data)
	})
	// サーバーをポート番号1323で起動
	e.Logger.Fatal(e.Start(":1323"))
}

// ハンドラーを定義
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
