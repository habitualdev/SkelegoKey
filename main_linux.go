package main

import (
	"github.com/cssivision/reverseproxy"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func main() {

	var local_port string

	dest_addr := os.Args[1]

	if len(os.Args) > 1 {
		local_port = os.Args[2]
	} else {
		local_port = "8989"
	}

	handle_proxy := func(w http.ResponseWriter, r *http.Request) {
		path, err := url.Parse(dest_addr)
		if err != nil {
			panic(err)
			return
		}
		proxy := reverseproxy.NewReverseProxy(path)
		proxy.ServeHTTP(w, r)
	}

	http.HandleFunc("/", handle_proxy)

	http_check, _ := regexp.Compile("http[s]?://.*")
	if http_check.Match([]byte(dest_addr)) {
		println("Connecting to " + dest_addr)
		http.ListenAndServe("127.0.0.1:"+local_port, nil)
	} else {
		println("Not a web address, try again.")
	}
}
