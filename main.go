package main
import (
	"fmt"
	"log"
	"html"
	"strings"
	"os/exec"
	"net/http"
	"github.com/gorilla/mux"
	"code.google.com/p/gcfg"
	"bytes"
	)
type Config struct {
	Copycat struct {
	Cmd string
	Arg string
	}
}

func httplog(r *http.Request){
        log.Printf("%s - %s - %s - %s - %q",
                r.RemoteAddr,
                r.Proto,
                r.Method,
                r.UserAgent(),
                html.EscapeString(r.URL.Path),
        )
}


func main(){
	rtr := mux.NewRouter()
	rtr.HandleFunc("/run",runcmd).Methods("GET").Queries("ID","{ID}")
        http.Handle("/", rtr)
        bind := ":8080"
        fmt.Printf("listening on %s...\n", bind)
        http.ListenAndServe(bind, nil)
}

func runcmd(w http.ResponseWriter, r *http.Request) {
	httplog(r)
	params := mux.Vars(r)
	fmt.Println(params["ID"])

	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "chaincmd.gcfg")
	if err != nil {
	    log.Fatalf("Failed to parse gcfg data: %s", err)
	}

/*
	fmt.Println(cfg.Copycat.Cmd)
	fmt.Println(cfg.Copycat.Arg)
	fmt.Println(strings.Replace(cfg.Copycat.Arg,"ID",params["ID"],-1))
*/

	Arg := strings.Replace(cfg.Copycat.Arg,"ID",params["ID"],-1)
	cmd := exec.Command(cfg.Copycat.Cmd,Arg)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("STDOUT:", out.String())

	ret := fmt.Sprintf("{\"Cmd\":\"%s\",\"Arg\":\"%s\"}",cfg.Copycat.Cmd,Arg)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte (ret))
}
