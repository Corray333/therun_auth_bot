.SILENT:
build:
	cd cmd && go build main.go
run: build
	cd cmd && ./main ../.env