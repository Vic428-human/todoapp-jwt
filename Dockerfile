FROM golang:1.25-alpine

WORKDIR /app

# 先複製依賴檔案
COPY go.* ./
RUN go mod download

# 再複製程式碼
COPY . .

# 編譯入口檔案 (cmd/api/main.go)
RUN go build -o main ./cmd/api/main.go

EXPOSE 8080

CMD ["./main"]
