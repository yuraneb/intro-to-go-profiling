package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/profile"
)

// Useless is not actually useless, it prevents potential compiler optimizations
var Useless uint64

func main() {

	iter := flag.Int("iter", 1000000, "number of iterations for main loop")

	flag.Parse()

	// using Dave Cheney's helpful wrapper
	defer profile.Start(profile.CPUProfile).Stop()

	rand.Seed(time.Now().Unix())

	var temp uint64

	for j := 1; j <= *iter; j++ {
		temp = rand.Uint64()
		Useless = temp
	}

	fmt.Printf("Final useless value is : %v", Useless)

}
