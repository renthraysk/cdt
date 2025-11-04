package cdt

import "testing"

func BenchmarkEtag(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {

		for s := range Etags("a,b,c").Tags {
			_ = s
		}
	}
}
