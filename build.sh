env GOOS=darwin GOARCH=amd64 go build -o dist/scaffoldeer_mac src/scaffoldeer.go
env GOOS=linux GOARCH=amd64 go build  -o dist/scaffoldeer_linux src/scaffoldeer.go
env GOOS=windows GOARCH=amd64 go build -o dist/scaffoldeer.exe src/scaffoldeer.go
