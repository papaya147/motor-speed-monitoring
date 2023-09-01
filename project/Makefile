MOTORSPEEDSERVER_BINARY=motorSpeedServerApp

build: build_library_server

build_library_server:
	@echo "Building motor speed server binary..."
	cd ../motor-speed-server && env GOOS=linux CGO_ENABLED=0 go build -o ${MOTORSPEEDSERVER_BINARY} ./cmd/api
	@echo "Done!"
