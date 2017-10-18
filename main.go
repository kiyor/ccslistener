/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 10-02-2017

* Last Modified : Wed 18 Oct 2017 08:07:13 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var listen *string = flag.String("l", ":8886", "listen interface")

func handler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Println(k, v)
		fmt.Fprintln(w, k, v)
	}
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	fmt.Println(string(b))
	fmt.Fprintln(w, string(b))
}

// request example curl -IL 'http://127.0.0.1/hub.docker.com?local=kiyor/stf:latest&remote=ccr.ccs.tencentyun.com/kiyor/stf:latest'
func handlerDockerhub(w http.ResponseWriter, r *http.Request) {
	local := r.URL.Query().Get("local")
	remote := r.URL.Query().Get("remote")
	go func(local, remote string) {
		cmd := fmt.Sprintf("docker pull %s && docker tag %s %s && docker push %s", local, local, remote, remote)
		c := exec.Command("/bin/sh", "-c", cmd)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err != nil {
			log.Println(err.Error())
		}
	}(local, remote)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/hub.docker.com", handlerDockerhub)
	http.ListenAndServe(*listen, r)
}
