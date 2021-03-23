package main

import (
	"fmt"
	"testing"
)

func TestIsMatched(t *testing.T) {
	stsName := "advanced-tidb-tikv"
	podName := "advanced-tidb-tikv-0"

	if res, err := isMatched(stsName, podName); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
