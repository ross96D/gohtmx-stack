package download

import (
	"fmt"
	"testing"
)

func TestHtmxHead(t *testing.T) {
	version, err := LatestHtmxVersion()
	if err != nil {
		fmt.Printf("t: %v\n", err)
		t.FailNow()
	}
	println(version)
}
