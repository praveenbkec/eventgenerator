package main

import "testing"

// BenchmarkGetEvent-12    	      51	  20834417 ns/op
// PASS
func BenchmarkGetEvent(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		GetEvent("12345")
	}
}
//BenchmarkListEvent-12    	      54	  22045760 ns/op
func BenchmarkListEvent(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		ListEvents()
	}
}