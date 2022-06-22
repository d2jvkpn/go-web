package misc

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"runtime"
)

/*
web browser address
  http://localhost:5060/debug/pprof/

get profiles and view in browser
  $ go tool pprof -http=:8080 http://localhost:5060/debug/pprof/allocs?seconds=30
  $ go tool pprof http://localhost:5060/debug/pprof/block?seconds=30
  $ go tool pprof http://localhost:5060/debug/pprof/goroutine?seconds=30
  $ go tool pprof http://localhost:5060/debug/pprof/heap?seconds=30
  $ go tool pprof http://localhost:5060/debug/pprof/mutex?seconds=30
  $ go tool pprof http://localhost:5060/debug/pprof/profile?seconds=30
  $ go tool pprof http://localhost:5060/debug/pprof/threadcreate?seconds=30

download profile file and convert to svg image
  $ wget -O profile.out localhost:5060/debug/pprof/profile?seconds=30
  $ go tool pprof -svg profile.out > profile.svg

get pprof in 30 seconds and save to svg image
  $ go tool pprof -svg http://localhost:5060/debug/pprof/allocs?seconds=30 > allocs.svg

get trace in 5 seconds
  $ wget -O trace.out http://localhost:5060/debug/pprof/trace?seconds=5
  $ go tool trace trace.out

get cmdline and symbol binary data
  $ wget -O cmdline.out http://localhost:5060/debug/pprof/cmdline
  $ wget -O symbol.out http://localhost:5060/debug/pprof/symbol
*/

// create new Pprof and run server
func LoadPprof(mux *http.ServeMux) {
	mux.HandleFunc("/debug/healthy", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte{})
	})

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)

	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	mux.HandleFunc("/debug/runtime/status", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")

		memStats := new(runtime.MemStats)
		runtime.ReadMemStats(memStats)
		num := runtime.NumGoroutine()

		json.NewEncoder(res).Encode(map[string]interface{}{
			"numGoroutine": num,
			"memStats":     memStats,
		})
	})

	return
}
