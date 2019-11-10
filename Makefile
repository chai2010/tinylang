all:
	go run ./tiny sum.tiny
	go run ./casl sum.casl
	go run . -f=sum.comet

.PHONY: tiny
tiny:
	go run ./tiny sum.tiny

.PHONY: casl
casl:
	go run ./casl sum.casl

.PHONY: vm
vm:
	go run . -d sum.coment

clean:
