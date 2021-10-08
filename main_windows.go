package main

import (
	"context"
	"github.com/cssivision/reverseproxy"
	"github.com/tadvi/winc"
	"net/http"
	"net/url"
)

func main() {

	connection_state := 0

	//http_check, _ := regexp.Compile("http[s]?://.*")

	var server http.Server
	var target_addr string
	ctx, cancel := context.WithCancel(context.Background())

	mux := http.NewServeMux()

	handle_proxy := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				return
			}
		}()
		path, err := url.Parse(target_addr)
		if err != nil {
			panic(err)
			return
		}
		proxy := reverseproxy.NewReverseProxy(path)
		proxy.ServeHTTP(w, r)
	}

	mux.HandleFunc("/", handle_proxy)

	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(350, 180) // (width, height)
	mainWindow.SetText("SkelegoKey Local Proxy")

	edt := winc.NewEdit(mainWindow)
	edt1 := winc.NewEdit(mainWindow)
	edt.SetSize(250, 20)
	edt1.SetSize(150, 20)
	edt.SetPos(10, 20)
	edt1.SetPos(10, 60)
	// Most Controls have default size unless SetSize is called.
	edt.SetText("Target Address")
	edt1.SetText("Local Proxy Port")
	btn := winc.NewPushButton(mainWindow)
	btn.SetText("Start Proxy")
	btn.SetPos(40, 85)   // (x, y)
	btn.SetSize(100, 40) // (width, height)

	btn.OnClick().Bind(func(e *winc.Event) {
		if connection_state == 0 {
			go func(ctx context.Context) {
				target_addr = edt.Text()
				server = http.Server{Addr: "127.0.0.1:" + edt1.Text(), Handler: mux}
				go server.ListenAndServe()
				btn.SetText("Shutdown")
				edt.Hide()
				edt1.Hide()
				connection_state = 1
			}(ctx)
		} else if connection_state == 1 {
			server.Shutdown(ctx)
			edt.Show()
			edt1.Show()
			btn.SetText("Connect")
			connection_state = 0
			cancel()
		}

	})
	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}
