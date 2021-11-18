package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	//"golang.org/x/net/http2"
	//"golang.org/x/net/http2/h2c"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

type Data struct {
	Fruit string `json:"fruit"`
}

func main() {
	startTestHTTPServer(3032, http.StatusOK, "{hello}")
	time.Sleep(100 * time.Hour)
}

func startTestHTTPServer(port, statusCode int, content string) (ts *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Data{Fruit: "12"}); err != nil {
			panic(fmt.Errorf("writeJSON failed: %w", err))
		}
	})
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic(err)
	}
	ts = &httptest.Server{
		Listener: listener,
		Config: &http.Server{
			Handler:      h2c.NewHandler(handler, &http2.Server{}),
			//Handler:      handler,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
			IdleTimeout:  time.Second,
		},
		EnableHTTP2: false,
	}
	ts.Start()
	return ts
}

//package main
//
//import (
//	"fmt"
//	"io"
//	"net/http"
//)
//
//func hello(w http.ResponseWriter, req *http.Request) {
//	fmt.Println(req.Method)
//	fmt.Println(req.URL)
//	fmt.Println(req.Proto)
//	bodyBytes, err := io.ReadAll(req.Body)
//	if err != nil {
//		fmt.Println(err)
//	}
//	bodyString := string(bodyBytes)
//	fmt.Println(bodyString)
//
//	fmt.Fprintf(w, "hello\n")
//}
//
//func main() {
//
//	ts := http.Server{
//		Addr:    ":3032",
//		//Handler: h2c.NewHandler(http.HandlerFunc(hello), &http2.Server{}),
//		Handler: http.HandlerFunc(hello),
//	}
//
//	ts.ListenAndServe()
//}
