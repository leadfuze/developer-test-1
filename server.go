package developertest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/yogihardi/developer-test-1/externalservice"
)

type (
	Server struct {
		Client externalservice.Client
		Echo   *echo.Echo
	}
)

func NewServer(eServer *echo.Echo, client externalservice.Client) *Server {
	return &Server{
		Client: client,
		Echo:   eServer,
	}
}

func (s *Server) AddRoutes() {
	// add routes
	apiEndpoint := s.Echo.Group("/api")
	apiEndpoint.POST("/posts/:id", s.Post)
	apiEndpoint.GET("/posts/:id", s.Get)
}

func (s *Server) Post(c echo.Context) error {
	// Get Params
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	title := c.FormValue("title")
	description := c.FormValue("description")

	// validate params
	if err != nil || title == "" || description == "" {
		return s.createBadResponseMsg(c, http.StatusBadRequest, "Please provide all parameters")
	}

	postPayload := &externalservice.Post{
		Title:       title,
		Description: description,
	}

	responsePost, err := s.Client.POST(id, postPayload)
	if err != nil {
		return s.createBadResponseMsg(c, http.StatusBadRequest, "Failed to get response from external client service")
	}

	if js, err := json.Marshal(responsePost); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed marshalling data to json: %s", err.Error()))
	} else {
		return c.JSONBlob(http.StatusOK, js)
	}
}

func (s *Server) Get(c echo.Context) error {
	// Get Params
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return s.createBadResponseMsg(c, http.StatusBadRequest, "Id should be an integer value")
	}

	responseGet, err := s.Client.GET(id)
	if err != nil {
		return s.createBadResponseMsg(c, http.StatusBadRequest, "Bad Request")
	}

	if js, err := json.Marshal(responseGet); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed marshalling data to json: %s", err.Error()))
	} else {
		return c.JSONBlob(http.StatusOK, js)
	}
}

func (s *Server) createBadResponseMsg(ec echo.Context, responseCode int, message string) error {
	response := fmt.Sprintf("{\"code\":%v, \"message\": \"%s\", \"path\":\"%s\"}", responseCode, message, ec.Request().URL.Path)
	return ec.JSONBlob(responseCode, []byte(response))
}
