all:
	go build -o output/timetracker -ldflags="-s -w" cmd/timetracker.go

windows:
	GOOS=windows GOARCH=amd64 go build -o output/timetracker.exe   -ldflags="-s -w" cmd/timetracker.go
