gen:
	@buf generate .

push:
	@buf push

build-auth-policy:
	@GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/protoc-gen-auth-policy protoc-gen-auth-policy/main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/protoc-gen-auth-policy protoc-gen-auth-policy/main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/protoc-gen-auth-policy.exe protoc-gen-auth-policy/main.go
	@GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64/protoc-gen-auth-policy protoc-gen-auth-policy/main.go

build-database-query:
	@GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/protoc-gen-database-query protoc-gen-database-query/main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/protoc-gen-database-query protoc-gen-database-query/main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/protoc-gen-database-query.exe protoc-gen-database-query/main.go
	@GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64/protoc-gen-database-query protoc-gen-database-query/main.go
