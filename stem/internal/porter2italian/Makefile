profile:
	go test -run=XXX -bench=BenchmarkStemBytes -cpuprofile=cpu.out
	go tool pprof porter2italian.test cpu.out

clean:
	go clean
	rm -f *.out


fuzz:
	go-fuzz-build gitlab.com/gonlp/stem/porter2italian
	go-fuzz -workdir=fuzzdir -bin=porter2italian-fuzz.zip

lint:
	gometalinter --disable=gotype
