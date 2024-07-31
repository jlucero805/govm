package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strata/interp"
	"strata/parser"
	"strings"
)

func Tmp(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	body, _ := ioutil.ReadAll(r.Body)
	bodyString := fmt.Sprintf("%s", body)
	_, log := interp.TopInterpLogs(parser.ParseAll(bodyString))
	io.WriteString(w, strings.Join(log.Logs, "\n"))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	// Set routing rules
	http.HandleFunc("/", Tmp)

	//Use the default DefaultServeMux.
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
