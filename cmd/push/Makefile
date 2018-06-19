VERSION := 1.2

default:
	gox -osarch="linux/amd64"
	docker build -t webup/push:$(VERSION) .
	docker push webup/push:$(VERSION)
