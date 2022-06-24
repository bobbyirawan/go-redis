package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	rds "test-pubsub/redis"

	"github.com/labstack/echo/v4"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/", Find)
	e.Logger.Fatal(e.Start(":8000"))
}

func Find(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}

	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := rds.RedisClient.Publish(context.Background(), "baymax.generate.img_prescription", payload); err != nil {
		log.Println("disini")
	}

	return c.String(200, "success")
}
