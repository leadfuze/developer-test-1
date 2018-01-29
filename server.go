package developertest

import (
	"net/http"
	"strconv"

	"github.com/enkhalifapro/developertest/externalservice"
	"github.com/labstack/echo"
)

type Server struct {
	Client externalservice.Client
}

func (s *Server) PostHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return err
	}
	var post externalservice.Post
	err = c.Bind(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return err
	}
	savedPost, err := s.Client.POST(id, &post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return err
	}
	return c.JSON(http.StatusCreated, savedPost)
}

func (s *Server) GetHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return err
	}
	post, err := s.Client.GET(id)
	if err != nil {
		result := make(map[string]string)
		result["code"] = "400"
		result["message"] = "Bad Request"
		result["path"] = "/api/posts/" + c.Param("id")
		return c.JSON(http.StatusBadRequest, result)
	}
	return c.JSON(http.StatusOK, post)
}
