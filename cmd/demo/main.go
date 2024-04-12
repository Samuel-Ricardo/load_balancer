package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.Int("Port", 8081, "Port to start the demo service on")

type DemoServer struct{}

func (s *DemoServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte(fmt.Sprintf("All Good! from server: %d", *port)))
}

func main() {
	flag.Parse()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), &DemoServer{}); err != nil {
		log.Fatal(err)
	}
}
