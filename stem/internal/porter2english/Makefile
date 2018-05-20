cpu: profile
	go tool pprof porter2english.test cpu.out

mem: profile
	go tool pprof --inuse_objects porter2english.test mem.out

profile:
	go test -run=XXX -bench=BenchmarkStemBytes -cpuprofile=cpu.out -memprofile=mem.out -benchtime=3s


clean:
	go clean
	rm -f *.out


fuzz:
	go-fuzz-build gitlab.com/gonlp/stem/porter2english
	go-fuzz -workdir=fuzzdir -bin=porter2english-fuzz.zip

lint:
	gometalinter --disable=gotype
