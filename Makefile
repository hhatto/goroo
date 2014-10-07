test:
	go test -v

bench:
	go test -v -parallel=2 -bench=. -benchmem=true
