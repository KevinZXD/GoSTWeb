package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func Routes() map[string]echo.HandlerFunc {
	routes := map[string]echo.HandlerFunc{
		"/get_recom": Recommend,
		"/aitao":     AITao,
	}

	return routes
}

func Recommend(c echo.Context) error {
	return c.String(http.StatusOK, "recom")
}

func AITao(c echo.Context) error {
	return c.String(http.StatusOK, "aitao")
}
