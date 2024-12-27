package main

import "github.com/juju/ratelimit"

func main() {
	bucket := ratelimit.NewBucketWithRate(100, 100)
	for i := 0; i < 1000; i++ {
		bucket.Wait(1)
	}
}
