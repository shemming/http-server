package proxy

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Proxy struct {
}

func NewProxy() *Proxy {
	return &Proxy{}
}

type Test struct {
	Test  string `json:"Test"`
	Test2 string `json:"Test2"`
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	_ = mux.Vars(r)
	ctx := r.Context()

	req := &Test{}
	err := readJSONBody(r, req)
	if err != nil {
		writeHTTPError(ctx, w, err)
		return
	}

	writeJSONResponse(ctx, w, req)
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
		writeHTTPError(ctx, w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Printf("could not write response: %s", err.Error())
	}
}
