package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

type addr struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

type conn struct {
	LocalAddr  addr `json:"local"`
	RemoteAddr addr `json:"remote"`
}

type response struct {
	Method     string              `json:"method"`
	URI        string              `json:"uri"`
	Headers    map[string][]string `json:"headers"`
	Connection conn                `json:"connection"`

	// RawPayload is almost always set (note: some requests, such
	// as GET requests, wont have a payload); Payload may be empty
	// if we're unable to parse it, or when there's no point doing
	// so
	RawPayload      string `json:"raw_payload"`
	rawPayloadBytes []byte `json:"-"`

	Payload interface{} `json:"payload"`
}

func (r *response) parsePayload() (err error) {
	payload := r.rawPayloadBytes

	if len(payload) == 0 {
		return
	}

	t, ok := r.Headers["Content-Type"]
	if !ok {
		return
	}

	switch t[0] {
	case "application/json":
		return json.Unmarshal(payload, &r.Payload)
	}

	return
}

// API exposes a super simple API which echos back
// any request made to it
type API struct {
	verbose bool
}

// Handle handles new connections and returns a payload
func (a API) Handle(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	r := &response{
		Method:          string(ctx.Method()),
		URI:             ctx.URI().String(),
		Headers:         make(map[string][]string),
		RawPayload:      string(body),
		rawPayloadBytes: body,
	}

	connection := ctx.Conn()
	if connection != nil {
		remote := connection.RemoteAddr()
		local := connection.LocalAddr()

		r.Connection = conn{
			LocalAddr: addr{
				Network: local.Network(),
				Address: local.String(),
			},
			RemoteAddr: addr{
				Network: remote.Network(),
				Address: remote.String(),
			},
		}
	}

	ctx.Request.Header.VisitAll(func(k, v []byte) {
		key := string(k)
		value := string(v)

		_, ok := r.Headers[key]
		if !ok {
			r.Headers[key] = make([]string, 0)
		}

		r.Headers[key] = append(r.Headers[key], value)
	})

	err := r.parsePayload()
	if err != nil && a.verbose {
		log.Printf("Error parsing payload: %+v", err)
	}

	body, err = json.Marshal(r)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}

	fmt.Fprintf(ctx, string(body))
}
