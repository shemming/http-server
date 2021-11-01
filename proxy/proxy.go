package proxy

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

type Proxy struct {
	Router         *mux.Router
	urlToShortCode map[string]string
	shortCodeToUrl map[string]string
}

func NewProxy(router *mux.Router) *Proxy {
	p := &Proxy{
		urlToShortCode: make(map[string]string),
		shortCodeToUrl: make(map[string]string),
		Router:         router,
	}
	p.Router.Handle("/", http.HandlerFunc(p.SetShortCode)).Methods(http.MethodPost)
	return p
}

type ShortCodeRequest struct {
	Url string `json:"url"`
}

type ShortCodeResponse struct {
	OriginalUrl  string `json:"url"`
	ShortUrlCode string `json:"short_url_code"`
}

type UrlRequest struct{}
type UrlResponse struct {
}

func (p *Proxy) SetShortCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &ShortCodeRequest{}
	err := readJSONBody(r, req)
	if err != nil {
		writeJsonHttpError(ctx, w, err)
		return
	}

	if _, ok := p.urlToShortCode[req.Url]; !ok {
		shortCode := "/" + generateCode(6)
		p.urlToShortCode[req.Url] = shortCode
		p.shortCodeToUrl[shortCode] = req.Url
		p.Router.Handle(shortCode, http.HandlerFunc(p.GetShortCode)).Methods(http.MethodGet)
	}

	resp := &ShortCodeResponse{
		OriginalUrl:  req.Url,
		ShortUrlCode: p.urlToShortCode[req.Url],
	}

	writeJSONResponse(ctx, w, resp)
}

func (p *Proxy) GetShortCode(w http.ResponseWriter, r *http.Request) {
	response := ""
	if oldUrl, ok := p.shortCodeToUrl[r.URL.Path]; ok {
		w.WriteHeader(http.StatusMovedPermanently)
		response = "Location: " + oldUrl
	} else {
		response = "Not relocated yet!"
	}
	w.Write([]byte(response))
}

func generateCode(lenCode int) string {
	// can also use a UUID for much greater scaling
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, lenCode)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func readJSONBody(r *http.Request, obj interface{}) error {
	tmp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(tmp, &obj); err != nil {
		return err
	}
	return nil
}

func writeJSONResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) {

	body, err := json.Marshal(resp)
	if err != nil {
		writeJsonHttpError(ctx, w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Printf("could not write response: %s", err.Error())
	}
}
