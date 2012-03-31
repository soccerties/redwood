// Redwood is an internet content-filtering program. 
// It is designed to replace and improve on DansGuardian 
// as the core of the Security Appliance internet filter. 
package main

import (
	"code.google.com/p/go-icap"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
)

var testURL = flag.String("test", "", "URL to test instead of running ICAP server")
var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
var cores = flag.Int("cores", runtime.NumCPU(), "number of CPU cores to use")

func main() {
	loadConfiguration()

	runtime.GOMAXPROCS(*cores)

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *testURL != "" {
		runURLTest(*testURL)
		return
	}

	go accessLog()

	startWebServer()
	go http.ListenAndServe(":1345", nil)

	icap.HandleFunc("/reqmod", handleRequest)
	icap.HandleFunc("/respmod", handleResponse)
	err := icap.ListenAndServe(":1344", nil)
	if err != nil {
		log.Println("Error listening on port 1344:", err)
	}
}
