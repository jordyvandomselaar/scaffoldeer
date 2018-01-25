env GOOS=darwin GOARCH=amd64 go build -o dist/scaffoldeer_mac scaffoldeer.go
env GOOS=linux GOARCH=amd64 go build  -o dist/scaffoldeer_linux scaffoldeer.go
env GOOS=windows GOARCH=amd64 go build -o dist/scaffoldeer.exe scaffoldeer.go
