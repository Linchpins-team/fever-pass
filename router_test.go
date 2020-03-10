// router_test.go
package main

import (
	"fmt"
	"testing"
)

func BenchmarkNewSelfRecord(b *testing.B) {
	body := fmt.Sprintf("temperature=%f&type=%d", 37.2, Ear)
	for i := 0; i < b.N; i++ {
		testHandler("POST", "/", body)
	}
}
