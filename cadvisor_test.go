package main

import (
	"fmt"

	"github.com/google/cadvisor/client"

	"github.com/golang/glog"
)

func main() {
	cadvisorClient, err := client.NewClient("http://127.0.0.1:4194")
	if err != nil {
		glog.Errorf("Error on creating cadvisor client: %v", err)
	}
	fmt.Println(cadvisorClient)
}
