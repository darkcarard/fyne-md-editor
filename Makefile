LINUX_BINARY_NAME=MarkDownEditor.tar.xz
APP_NAME=MarkDownEditor
VERSION=1.0.0

## build: build binary and package app
linux-build:
	rm -rf ${LINUX_BINARY_NAME}
	fyne package -os linux -icon icon.png -appVersion ${VERSION} -name ${APP_NAME}

## build: build binary and package app
win-build:
	fyne package -os windows -icon icon.png -appVersion ${VERSION} -name ${APP_NAME}

mac-build
	fyne package -os darwin -icon icon.png -appVersion ${VERSION} -name ${APP_NAME}
## run: builds and runs the application
run:
	go run .

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf ${BINARY_NAME}
	@echo "Cleaned!"

## test: runs all tests
test:
	go test -v ./...