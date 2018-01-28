package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/yogihardi/developer-test-1/externalservice"
)

var (
	server           *Server
	e                *echo.Echo
	postRequestCount = 0
	getRequestCount  = 0
	getId            = 0
)

type MockClient struct{}

func (mockClient *MockClient) GET(id int) (*externalservice.Post, error) {
	getRequestCount += 1
	getId = id
	return nil, errors.New("error")
}

func (mockClient *MockClient) POST(id int, post *externalservice.Post) (*externalservice.Post, error) {
	postRequestCount += 1
	post.ID = id
	return post, nil
}

func init() {
	e = echo.New()
	server = NewServer(e, &MockClient{})
}

func TestPOSTCallsAndReturnsJSONfromExternalServicePOST(t *testing.T) {
	// Descirption
	//
	// Write a test that accepts a POST request on the server and sends it the
	// fake external service with the posted form body return the response.
	//
	// Use the externalservice.Client interface to create a mock and assert the
	// client was called <n> times.
	//
	// ---
	//
	// Server should receive a request on
	//
	//  [POST] /api/posts/:id
	//  application/json
	//
	// With the form body
	//
	//  application/x-www-form-urlencoded
	//	title=Hello World!
	//	description=Lorem Ipsum Dolor Sit Amen.
	//
	// The server should then relay this data to the external service by way of
	// the Client POST method and return the returned value out as JSON.
	//
	// ---
	//
	// Assert that the externalservice.Client#POST was called 1 times with the
	// provided `:id` and post body and that the returned Post (from
	// externalservice.Client#POST) is written out as `application/json`.

	req := httptest.NewRequest(echo.POST, "/api/posts/87", nil)
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")

	if err := req.ParseForm(); err != nil {
		panic(err.Error())
	}

	req.Form.Add("title", "Hello World!")
	req.Form.Add("description", "Lorem Ipsum Dolor Sit Amen.")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("87")

	// execute get post method
	server.Post(c)

	if rec.Code != http.StatusOK {
		t.Fail()
	}

	expectedBody := `{"id":87,"title":"Hello World!","description":"Lorem Ipsum Dolor Sit Amen."}`
	if rec.Body.String() != expectedBody {
		t.Fail()
	}

	if postRequestCount != 1 {
		t.Fail()
	}

}

func TestPOSTCallsAndReturnsErrorAsJSONFromExternalServiceGET(t *testing.T) {
	// Description
	//
	// Write a test that accepts a GET request on the server and returns the
	// error returned from the external service.
	//
	// Use the externalservice.Client interface to create a mock and assert the
	// client was called <n> times.
	//
	// ---
	//
	// Server should receive a request on
	//
	//	[GET] /api/posts/:id
	//
	// The server should then return the error from the external service out as
	// JSON.
	//
	// The error response returned from the external service would look like
	//
	//	400 application/json
	//
	//	{
	//		"code": 400,
	//		"message": "Bad Request"
	//	}
	//
	// ---
	//
	// Assert that the externalservice.Client#GET was called 1 times with the
	// provided `:id` and the returned error (above) is output as the response
	// as
	//
	//	{
	//		"code": 400,
	//		"message": "Bad Request",
	//		"path": "/api/posts/:id
	//	}
	//
	// Note: *`:id` should be the actual `:id` in the original request.*

	req := httptest.NewRequest(echo.GET, "/api/posts/87", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("87")

	server.Get(c)

	if rec.Code != http.StatusBadRequest {
		t.Fail()
	}

	expectedBody := `{"code":400, "message": "Bad Request", "path":"/api/posts/87"}`
	if rec.Body.String() != expectedBody {
		t.Fail()
	}

	if getRequestCount != 1 {
		t.Fail()
	}

	if getId != 87 {
		t.Fail()
	}
}
