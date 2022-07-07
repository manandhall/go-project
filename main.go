package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/rollbar/rollbar-go"
	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error(err)
	}
	port := os.Getenv("PORT")
	app := newApp()

	app.Run(iris.Addr(":"+port), iris.WithoutServerError(iris.ErrServerClosed))

}

func newApp() *iris.Application {

	app := iris.New()
	app.Use(recover2.New())

	app.Use(requestid.New())
	app.Use(func(c iris.Context) {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				file, line := c.HandlerFileLine()
				rollbar.Critical(errors.New(fmt.Sprint(r)), iris.Map{
					"request_id":  c.GetID(),
					"request_ip":  c.RemoteAddr(),
					"request_uri": c.FullRequestURI(),
					"handler": iris.Map{
						"name": c.HandlerName(),
						"file": fmt.Sprintf("%s%d", file, line),
					},
				})
				c.StopWithStatus(iris.StatusInternalServerError)
			}
		}()

		c.Next()
	})

	app.Logger().Install(log.StandardLogger())

	app.Get("ming", func(context *context.Context) {
		context.Text("out")
		return
	})

	return app
}
