cd /D "%~dp0"
set GOOS=linux
set GOARCH=amd64
go build -o compiled/calc_lag_monitor