@echo off
go mod init bruteforce
go get github.com/go-vgo/robotgo
go build -ldflags="-H windowsgui -s -w"
pause