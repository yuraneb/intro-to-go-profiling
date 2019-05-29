package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"runtime"

	"sync"
	"time"

	"github.com/pkg/profile"
)

// Useless is not actually useless, it prevents potential compiler optimizations
var Useless uint64

type Options struct {
	iter     int
	requests int
	time     int
}

func Init() {
	runtime.SetCPUProfileRate(10000)
}

const URL = "http://api.timezonedb.com/v2.1/get-time-zone"

func main() {

	opts := Options{}

	opts.iter = *flag.Int("iter", 10000000, "number of loop iterations")
	opts.requests = *flag.Int("requests", 10, "number of GET requests")
	opts.time = *flag.Int("time", 20, "http timeout in seconds")

	flag.Parse()

	// using Dave Cheney's helpful wrapper
	defer profile.Start(profile.CPUProfile).Stop()

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(opts.time)*time.Second)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		networkThread(ctx, opts.requests)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cpuThread(opts.iter)
	}()

	wg.Wait()

	fmt.Println("All done.")

}

// networkThread targets a remote API
func networkThread(ctx context.Context, reqs int) {
	startTime := time.Now()
	newclient := http.Client{}

	for i := 0; i < reqs; i++ {

		newRequest, _ := http.NewRequest("GET", URL, nil)
		newRequest = newRequest.WithContext(ctx)
		resp, _ := newclient.Do(newRequest)

		body := make([]byte, 0)
		body, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("Response %v : %v\n", i, string(body))

		time.Sleep(time.Second * 1)

	}

	endTime := time.Now().Sub(startTime)
	fmt.Printf("Total running time for networkThread : %v", endTime)

}

// cpuThread does some cpu/memory-only work
func cpuThread(iter int) {
	rand.Seed(time.Now().Unix())
	var temp uint64

	startTime := time.Now()

	for j := 1; j <= iter; j++ {
		temp = rand.Uint64()
		Useless = temp
	}

	fmt.Printf("Final useless value is : %v\n", Useless)
	endTime := time.Now().Sub(startTime)
	fmt.Printf("Total running time for cpuThread : %v", endTime)
}
