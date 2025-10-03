set app=listtokens

set HTTP_PROXY=http://127.0.0.1:58080
set HTTPS_PROXY=http://127.0.0.1:58080

go mod init %app%
go mod tidy

SET GOOS=windows
SET GOARCH=amd64
SET CGO_ENABLED=0
go build -buildvcs=false -ldflags="-s -w -buildid=" -trimpath -o bin/%app%.exe

SET GOOS=linux
SET GOARCH=amd64
SET CGO_ENABLED=0
go build -buildvcs=false -ldflags="-s -w -buildid=" -trimpath -o bin/%app%_%GOOS%_%GOARCH%
