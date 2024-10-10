package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Jamlie/go-paginate/db"
)

type PaginateRequest struct {
	PageSize int        `json:"pageSize"`
	Page     int        `json:"page"`
	OrderBy  db.OrderBy `json:"orderBy,omitempty"`
	Filters  db.Filters `json:"filters,omitempty"`
}

type Response struct {
	StatusCode int       `json:"status_code,omitempty"`
	Message    string    `json:"message,omitempty"`
	Data       []db.User `json:"data,omitempty"`
}

func sendJson(w http.ResponseWriter, res Response) {
	w.WriteHeader(res.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func main() {
	users, err := db.Create()
	if err != nil {
		log.Println(err)
		return
	}
	defer users.Close()

	port := flag.Int("port", 8080, "Set the port of the server")
	flag.Parse()

	app := echo.New()

	app.POST("/insert", func(c echo.Context) error {
		user := db.User{}
		if err := c.Bind(&user); err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "Could not decode JSON",
			})
		}

		err = users.Insert(user)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert user",
			})
		}

		return c.JSON(http.StatusOK, Response{
			StatusCode: http.StatusOK,
			Message:    "User inserted successfully",
		})
	})

	app.GET("/show", func(c echo.Context) error {
		allUsers, err := users.RetrieveAll()
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve users",
			})
		}

		buf := users.Show(allUsers)
		return c.String(http.StatusOK, buf.String())
	})

	app.POST("/paginate", func(c echo.Context) error {
		var req PaginateRequest
		if err := c.Bind(&req); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request payload",
			})
		}

		if req.PageSize == 0 {
			return c.JSON(http.StatusBadRequest, Response{
				StatusCode: http.StatusBadRequest,
				Message:    "pageSize must be specified",
			})
		}

		if req.Page == 0 {
			return c.JSON(http.StatusBadRequest, Response{
				StatusCode: http.StatusBadRequest,
				Message:    "page must be specified",
			})
		}

		retrievedUsers, err := users.Paginate(req.PageSize, req.Page, req.OrderBy, req.Filters)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve users",
			})
		}

		return c.JSON(http.StatusOK, Response{
			StatusCode: http.StatusOK,
			Message:    "Sent data successfully!",
			Data:       retrievedUsers,
		})
	})

	if err := app.Start(fmt.Sprintf(":%d", *port)); err != nil {
		log.Panic(err)
	}
}
