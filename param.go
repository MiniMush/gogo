package gogo

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/golib/httprouter"
)

type AppParams struct {
	mux     sync.RWMutex
	request *http.Request
	params  httprouter.Params
	rawBody []byte
	rawErr  error
	readed  bool
}

func NewAppParams(r *http.Request, params httprouter.Params) *AppParams {
	return &AppParams{
		request: r,
		params:  params,
	}
}

// HasQuery returns whether named param is exist for URL query string.
func (p *AppParams) HasQuery(name string) bool {
	_, ok := p.request.URL.Query()[name]

	return ok
}

// HasForm returns whether named param is exist for POST/PUT request body.
func (p *AppParams) HasForm(name string) bool {
	p.request.ParseMultipartForm(DefaultMaxMultiformBytes)

	_, ok := p.request.PostForm[name]

	return ok
}

func (p *AppParams) RawBody() ([]byte, error) {
	if !p.readed {
		p.mux.Lock()
		if !p.readed {
			p.rawBody, p.rawErr = ioutil.ReadAll(p.request.Body)

			// close the request.Body
			p.request.Body.Close()

			p.readed = true
		}
		p.mux.Unlock()
	}

	return p.rawBody, p.rawErr
}

// Get returns the first value for the named component of the request.
// NOTE: httprouter.Params takes precedence over URL query string values.
func (p *AppParams) Get(name string) string {
	value := p.params.ByName(name)

	// trye URL query string if value of route is empty
	if value == "" {
		value = p.request.URL.Query().Get(name)
	}

	return value
}

// Post returns the named comonent of the request by calling http.Request.FormValue()
func (p *AppParams) Post(name string) string {
	return p.request.FormValue(name)
}

// File retrieves multipart uploaded file of HTTP POST request
func (p *AppParams) File(name string) (multipart.File, *multipart.FileHeader, error) {
	return p.request.FormFile(name)
}

// Json unmarshals request body with json codec
func (p *AppParams) Json(v interface{}) error {
	data, err := p.RawBody()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// Xml unmarshals request body with xml codec
func (p *AppParams) Xml(v interface{}) error {
	data, err := p.RawBody()
	if err != nil {
		return err
	}

	return xml.Unmarshal(data, v)
}

// Gob decode request body with gob codec
func (p *AppParams) Gob(v interface{}) error {
	data, err := p.RawBody()
	if err != nil {
		return err
	}

	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}

// Bson unmarshals request body with bson codec
func (p *AppParams) Bson(v interface{}) error {
	data, err := p.RawBody()
	if err != nil {
		return err
	}

	return bson.Unmarshal(data, v)
}
