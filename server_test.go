package developertest

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/enkhalifapro/developertest/externalservice"
	"github.com/labstack/echo"
)

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

	// Kh code

	// create echo instance
	e := echo.New()
	// prepare request
	body := bytes.NewBufferString("title=Hello World!&description=Lorem Ipsum Dolor Sit Amen.")
	req := httptest.NewRequest("POST", "/api/posts", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// add param id with test value '1'
	c.SetParamNames("id")
	c.SetParamValues("1")

	server := &Server{
		Client: &externalservice.ClientMock{
			Posts: make(map[int]*externalservice.Post),
		}}

	// First call with id = 1
	err := server.PostHandler(c)

	// error should be nil
	if err != nil {
		t.Error(err.Error())
	}
	// StatusCode should be Created 201
	if rec.Code != 201 {
		t.Error("Invalid status code")
	}

	// Response content-type should be 'application/json; charset=UTF-8'
	if rec.Header().Get("Content-Type") != "application/json; charset=UTF-8" {
		t.Error("Invalid response Content-Type")
	}

	// Second call with id = 1
	err = server.PostHandler(c)

	// error message should be 'Post id is already called'
	if err == nil || err.Error() != "Post id is already called" {
		t.Error("Should get error with message  'Post id is already called'")
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

	// Kh code

	// create echo instance
	e := echo.New()
	// prepare request
	body := bytes.NewBufferString("title=Hello World!&description=Lorem Ipsum Dolor Sit Amen.")
	req := httptest.NewRequest("POST", "/api/posts", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// add param id with test value '1'
	c.SetParamNames("id")
	c.SetParamValues("1")

	server := &Server{
		Client: &externalservice.ClientMock{
			Posts: make(map[int]*externalservice.Post),
		}}

	// First call with id = 1
	err := server.PostHandler(c)

	// error should be nil
	if err != nil {
		t.Error(err.Error())
	}
	// StatusCode should be Created 201
	if rec.Code != 201 {
		t.Error("Invalid status code")
	}

	// Response content-type should be 'application/json; charset=UTF-8'
	if rec.Header().Get("Content-Type") != "application/json; charset=UTF-8" {
		t.Error("Invalid response Content-Type")
	}

	// Second call is getting the created post
	req = httptest.NewRequest("GET", "/api/posts", body)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	// add param id with test value '1'
	c.SetParamNames("id")
	c.SetParamValues("1")
	err = server.GetHandler(c)

	// status code should be 200
	if rec.Code != 200 {
		t.Error("Status code should be 200")
	}

	// third call is getting not found post
	req = httptest.NewRequest("GET", "/api/posts", body)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	// add param id with test value '1'
	c.SetParamNames("id")
	c.SetParamValues("2")
	err = server.GetHandler(c)

	// status code should be 400
	if rec.Code != 400 {
		t.Error("Status code should be 400")
	}
	// response schema should be
	//	{
	//		"code": 400,
	//		"message": "Bad Request",
	//		"path": "/api/posts/:id
	//	}
	expectedResponse := "{\"code\":\"400\",\"message\":\"Bad Request\",\"path\":\"/api/posts/2\"}"
	response := rec.Body.String()
	if response != expectedResponse {
		t.Error("Invalid json response schema")
	}
}
