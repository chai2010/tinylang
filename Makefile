all:
	go run ./tiny sum.tiny
	go run ./casl sum.casl
	go run ./vm sum

.PHONY: tiny
tiny:
	go run ./tiny sum.tiny

.PHONY: casl
casl:
	go run ./casl sum.casl

.PHONY: vm
vm:
	go run ./vm -debug sum

clean:
