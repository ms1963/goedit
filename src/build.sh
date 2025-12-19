
# Initialize
go mod init goedit
go mod tidy

# Build
go build -o goedit

# Windows
go build -o goedit.exe

# Cross-compile
#GOOS=windows GOARCH=amd64 go build -o goedit.exe
#GOOS=linux GOARCH=amd64 go build -o goedit
#GOOS=darwin GOARCH=amd64 go build -o goedit
GOOS=darwin GOARCH=arm64 go build -o goedit
