package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestAPI_Serve(t *testing.T) {
	for _, test := range []struct {
		name          string
		body          string
		contentType   string
		expectPayload interface{}
	}{
		{"Plain text request", "hello world!", "text/plain", nil},
		{"JSON content", `{"hello":"world"}`, "application/json", map[string]interface{}{"hello": "world"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseRequest(req)

			req.SetBodyString(test.body)
			req.SetRequestURI("http://example.com/")
			req.Header.SetMethod("POST")
			req.Header.Add("content-type", test.contentType)
			req.Header.Add("my-test-value", "hello!")

			resp := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseResponse(resp)

			ctx := fasthttp.RequestCtx{
				Request:  *req,
				Response: *resp,
			}

			API{true}.Handle(&ctx)

			responseData := response{}
			err := json.Unmarshal(ctx.Response.Body(), &responseData)
			if err != nil {
				t.Fatalf("error unmarshaling response data: %+v", err)
			}

			for _, vTest := range []struct {
				name   string
				val    string
				expect string
			}{
				{"method string", responseData.Method, "POST"},
				{"uri", responseData.URI, "http://example.com/"},
				{"test header", responseData.Headers["My-Test-Value"][0], "hello!"},
			} {
				t.Run(vTest.name, func(t *testing.T) {
					if vTest.expect != vTest.val {
						t.Errorf("expected %q, received %q", vTest.expect, vTest.val)
					}
				})

			}

			t.Run("parsed body", func(t *testing.T) {
				received := responseData.Payload

				if !reflect.DeepEqual(test.expectPayload, received) {
					t.Errorf("expected %+v, received %+v", test.expectPayload, received)
				}
			})
		})
	}
}
