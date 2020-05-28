package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
func addSrv(c echo.Context) error {

	addsrv := Server{}
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadGateway, fmt.Sprintf("%s", "error happend"))
	}
	err = json.Unmarshal(b, &addsrv)
	return c.String(http.StatusOK, fmt.Sprintf("%s", addsrv.Srv_name))
}
func getSrv(c echo.Context) error {

	db, err := gorm.Open("mysql", "gouser:password@(localhost)/godb?charset=utf8&parseTime=True&loc=Local")
	if err != nil {

		return c.String(http.StatusBadRequest, fmt.Sprintf("%s", err.Error()))

	}
	defer db.Close()

	var srv1 Server
	srv1.Srv_name = c.QueryParam("srv_name")
	srv1.Ram = (c.QueryParam("ram"))
	srv1.Cpu = c.QueryParam("cpu")
	srv1.Owner = c.QueryParam("owner")
	srv1.Login_user = c.QueryParam("login_user")
	srv1.Description = c.QueryParam("desc")

	datatype := c.Param("data")

	if !db.HasTable(srv1) {
		db.CreateTable(srv1)
	}
	db.Create(&srv1)

	if datatype == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("%s", srv1.Owner))
	}

	if datatype == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"owner":       srv1.Owner,
			"server_name": srv1.Srv_name,
		})
	}

	return c.String(http.StatusBadRequest, string("error"))
}
func main() {
	e := echo.New()
	e.GET("/", ping)
	e.GET("/servers/:data", getSrv)
	e.GET("/newserver/", addSrv)

	e.Logger.Fatal(e.Start(":8000"))
}

type Server struct {
	gorm.Model
	Srv_name    string
	Ram         string
	Cpu         string
	Owner       string
	Login_user  string
	Login_pass  string
	Ips         string
	Description string
}
