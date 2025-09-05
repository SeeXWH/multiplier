// main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type Generator struct {
	rtp float64
	p1  float64
	p2  float64
}

func NewGenerator(rtp float64) *Generator {
	return &Generator{
		rtp: rtp,
		p1:  1.0 - rtp,
		p2:  1.0 - rtp/10000.0,
	}
}

func (g *Generator) Sample() float64 {
	u := rand.Float64()
	if u < g.p1 {
		return 1.0
	}
	if u >= g.p2 {
		return 10000.0
	}
	x := g.rtp / (1.0 - u)
	return x
}

type response struct {
	Result float64 `json:"result"`
}

func main() {
	var rtpFlag float64
	flag.Float64Var(&rtpFlag, "rtp", 0, "target rtp in (0, 1]")
	flag.Parse()
	gen := NewGenerator(rtpFlag)
	mux := http.NewServeMux()
	mux.HandleFunc("/get", func(w http.ResponseWriter, req *http.Request) {
		val := gen.Sample()
		w.Header().Set("Content-Type", "application/json")
		res := response{val}
		resBytes, _ := json.Marshal(res)
		w.Write(resBytes)
		fmt.Println(val)
	})

	server := &http.Server{
		Addr:    ":64333",
		Handler: mux,
	}
	log.Println("Start")
	server.ListenAndServe()
}
