.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o wechatbot ./main.go

.PHONY: docker
docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -installsuffix cgo -ldflags '-w' -o wechatbot ./main.go
	docker build . -t wechatbot:latest
