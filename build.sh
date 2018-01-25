env GOOS=darwin GOARCH=amd64 go build -o dist/mac/scaffoldeer scaffoldeer.go
env GOOS=linux GOARCH=amd64 go build  -o dist/linux/scaffoldeer scaffoldeer.go
env GOOS=windows GOARCH=amd64 go build -o dist/windows/scaffoldeer.exe scaffoldeer.go
