build: build_fe build_be

build_fe: t/
	yarn build

build_be: main.go
	go build

dev: build_be
	yarn dev & DEV=y CL="http://localhost:3000" ./idk