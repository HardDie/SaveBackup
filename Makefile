all: windows linux darwin

windows: bin
	GOOS=windows GOARCH=amd64 go build -o bin/windows_amd64.exe .
	GOOS=windows GOARCH=386 go build -o bin/windows_386.exe .

linux: bin
	GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64 .
	GOOS=linux GOARCH=386 go build -o bin/linux_386 .
	GOOS=linux GOARCH=arm go build -o bin/linux_arm .
	GOOS=linux GOARCH=arm64 go build -o bin/linux_arm64 .

darwin: bin
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin_arm64 .

bin:
	@mkdir bin
