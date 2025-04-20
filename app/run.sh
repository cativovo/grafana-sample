cd /app

go mod tidy
go tool air --build.cmd "go build -o ./tmp/server ./cmd/server" --build.bin "./tmp/server"
