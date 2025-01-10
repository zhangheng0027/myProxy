package test

import (
	"fmt"
	"github.com/robfig/cron/v3"
)
import "testing"

func TestCron(t *testing.T) {
	c := cron.New()
	c.AddFunc("@every 1s", aaa)
	c.Start()
	select {}
}

func aaa() {
	fmt.Println("aaa")
}
