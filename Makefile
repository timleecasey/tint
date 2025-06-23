


all: test takehome


takehome:
	go build takehome.go expiry_pq.go stacked_list.go time_provider.go

test:
	go test -json ./...