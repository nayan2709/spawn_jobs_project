package mockServer

import (
	"errors"
	"fmt"
	"github.com/dunzoit/projects/nirwana_spwanJob/entities"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func NirwanaMockServer(payload entities.Payload) error {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))

	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		return errors.New("error while making api call")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("error while making api call")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("error while reading response body")
	}
	fmt.Println(string(body))
	return nil
}
