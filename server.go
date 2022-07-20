package observerip

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Server struct {
	r *http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, PUT, POST, OPTIONS")
	if acrh, ok := r.Header["Access-Control-Request-Headers"]; ok {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(acrh, ", "))
	} else {
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Proxy-Authorization, Proxy-Authenticate")
	}
	if WU_PASSTHROUGH {
		var success chan error
		go passthrough(r, success)
		if err := <-success; err != nil {
			log.Printf("ERROR handling passthrough: %s\n", err)
		}
	}
	s.r.ServeHTTP(w, r)
}

func passthrough(r *http.Request, ec chan error) error {
	url, err := url.Parse(fmt.Sprintf("%s/%s", WU_URL, r.RequestURI))
	if err != nil {
		ec <- err
		return err
	}
	c := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(r.Method, url.String(), nil)
	if err != nil {
		ec <- err
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		ec <- err
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ec <- err
		return err
	}
	if string(body) != "success" || resp.StatusCode != http.StatusOK {
		log.Printf("Error passing through url=%s resp=%+v\n", url, resp)
		err = errors.New("Request not accepted!")
		ec <- err
		return err
	}
	return nil
}
