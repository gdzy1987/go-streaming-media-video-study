package main

import (
	"fmt"
	"testing"
)

func testPrint1(t *testing.T) {
	// t.SkipNow()
	res := Print1to20()
	fmt.Println("hey")
	if res != 210 {
		t.Errorf("Wrong result of Print1to20")
	}
}

// 不执行
func testPrint2(t *testing.T) {
	fmt.Println("test")
}

func TestAll(t *testing.T) {
	t.Run("TestPrint1", testPrint1)
	t.Run("TestPrint2", testPrint2)
}

// 顺序执行
func TestPrintSub(t *testing.T) {
	t.Run("a1", func(t *testing.T) { fmt.Println("a1") })
	t.Run("a2", func(t *testing.T) { fmt.Println("a2") })
	t.Run("a3", func(t *testing.T) { fmt.Println("a3") })
}

// TestMain
func TestMain(m *testing.M) {
	fmt.Println("test main first")
	m.Run()
}

// go test -bench=.
func BenchmarkAll(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Print1to20()
	}
}
