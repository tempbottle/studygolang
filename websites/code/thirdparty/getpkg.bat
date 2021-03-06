@echo off

setlocal

if exist getpkg.bat goto ok
echo getpkg.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

go get -u -v github.com/go-sql-driver/mysql
go get -u -v github.com/studygolang/mux
go get -u -v github.com/gorilla/sessions
go get -u -v github.com/robfig/cron
go get -u -v github.com/qiniu/api
go get -u -v github.com/dchest/captcha

if exist src/golang.org/x/text (
	git clone https://github.com/golang/text src/golang.org/x/text
)
go install golang.org/x/text/...

if exist src/golang.org/x/net (
	git clone https://github.com/golang/net src/golang.org/x/net
)
go install golang.org/x/net/...
#go get -u -v golang.org/x/net

if src/code.google.com/p/cascadia (
	git clone https://github.com/studygolang/cascadia src/code.google.com/p/cascadia
)
go install code.google.com/p/cascadia

go get -u -v github.com/PuerkitoBio/goquery

set GOPATH=%OLDGOPATH%

:end
echo finished