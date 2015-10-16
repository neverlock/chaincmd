package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/neverlock/utility/httplog"
	"gopkg.in/gcfg.v1"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Config struct {
	Copycat struct {
		Cmd  string
		Arg  string
		Bind string
	}
}

var cfg Config

func main() {

	err := gcfg.ReadFileInto(&cfg, "chaincmd.gcfg")
	if err != nil {
		log.Fatalf("Failed to parse gcfg data: %s", err)
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/run", runcmd).Methods("GET").Queries("ID", "{ID}")
	http.Handle("/", rtr)
	rtr.NotFoundHandler = http.HandlerFunc(notFound)
	bind := cfg.Copycat.Bind
	fmt.Printf("listening on %s...\n", bind)
	http.ListenAndServe(bind, nil)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Neverlock Chaincmd")
	w.Write([]byte("404 End point not found"))
}

func runcmd(w http.ResponseWriter, r *http.Request) {
	httplog.HttpLogln(r)
	params := mux.Vars(r)
	Arg := strings.Replace(cfg.Copycat.Arg, "ID", params["ID"], -1)
	cmd := exec.Command("/bin/sh", "-c", cfg.Copycat.Cmd, Arg)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	}

	log.Println("STDOUT:", out.String())

	ret := fmt.Sprintf("{\"Cmd\":\"%s\",\"Arg\":\"%s\"}", cfg.Copycat.Cmd, Arg)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(ret))
}
