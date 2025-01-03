package main

import (
	"github.com/zhangheng0027/ratelimit-plus"
)

func main() {

	bucket := ratelimit.NewBucketWithRate(100, 100)
	for i := 0; i < 1000; i++ {
		bucket.Wait(1)
	}
}
