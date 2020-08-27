run:
	go run ./main.go

test:
	go test ./ -bench Count -v -count 1 -cpuprofile cpu.out -memprofile mem.out -benchtime 10x