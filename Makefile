gen:
	@buf generate .

push:
	@buf push

build-plugin:
	@GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/protoc-gen-auth-policy protoc-gen-auth-policy/main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/protoc-gen-auth-policy protoc-gen-auth-policy/main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/protoc-gen-auth-policy.exe protoc-gen-auth-policy/main.go
	@GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64/protoc-gen-auth-policy protoc-gen-auth-policy/main.go
