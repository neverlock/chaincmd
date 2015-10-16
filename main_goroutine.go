package main

import (
	"bytes"
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Config struct {
	Copycat struct {
		Cmd  string
		Arg  string
		Arg1 string
		Bind string
	}
}

var cfg Config

func httplog(r *http.Request) {
	log.Printf("%s - %s - %s - %s - %q",
		r.RemoteAddr,
		r.Proto,
		r.Method,
		r.UserAgent(),
		html.EscapeString(r.URL.Path),
	)
}

func main() {

	err := gcfg.ReadFileInto(&cfg, "chaincmd.gcfg")
	if err != nil {
		log.Fatalf("Failed to parse gcfg data: %s", err)
	}
	//var outPut chan string = make(chan string)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/run", runcmd).Methods("GET").Queries("ID", "{ID}")
	http.Handle("/", rtr)
	bind := cfg.Copycat.Bind
	fmt.Printf("listening on %s...\n", bind)
	http.ListenAndServe(bind, nil)
}

func execmd(cmd string, arg1 string, arg2 string) {
	exe := exec.Command(cmd, arg1, arg2)
	var out bytes.Buffer
	var stderr bytes.Buffer
	exe.Stdout = &out
	exe.Stderr = &stderr
	err := exe.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	}
	log.Println("STDOUT:", out.String())
}

func runcmd(w http.ResponseWriter, r *http.Request) {
	httplog(r)
	params := mux.Vars(r)
	Arg := strings.Replace(cfg.Copycat.Arg, "ID", params["ID"], -1)
	log.Println("Arg := ", Arg)

	go execmd(cfg.Copycat.Cmd, Arg, cfg.Copycat.Arg1)

	ret := fmt.Sprintf("{\"Cmd\":\"%s\",\"Arg\":\"%s\"}", cfg.Copycat.Cmd, Arg)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(ret))
}
