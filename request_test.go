// request_test.go
package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	go main()
	// cmd := exec.Command("dig", "@localhost -p 8090 google.com")
	out, _ := exec.Command("dig", "@localhost", "-p8090", "google.com").Output()

	fmt.Println(string(out))
	// send DNS packet to this server
}
