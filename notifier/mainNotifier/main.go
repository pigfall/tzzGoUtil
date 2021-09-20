package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pigfall/tzzGoUtil/output"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type Msg struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func main() {
	var port string
	flag.StringVar(&port, "port", "10101", "port")
	flag.Parse()
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			output.Errf("read request err:%v\n", err)
			return
		}
		msg := &Msg{}
		err = json.Unmarshal(bytes, msg)
		if err != nil {
			output.Errf("unmarshal req msg failed: %v\n", err)
			return
		}
		err = exec.Command("notify", "-t", msg.Title, "-m", msg.Message).Start()
		if err != nil {
			output.Errf("start cmd notify failed: %v\n", err)
			return
		}
	})
	fmt.Printf("start server at 0.0.0.0:%s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		fmt.Printf("serve over: %v\n", err)
	} else {
		fmt.Println("serve over")
	}
}
