# Stage 1: Build binary statically linked
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cài đặt git và ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy toàn bộ mã nguồn
COPY . .

# Đồng bộ dependencies và tải các thư viện cần thiết
RUN go mod tidy && go mod download

# Biên dịch tĩnh từ cmd/server, loại bỏ debug symbols (-w -s)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Stage 2: Production Minimal Image
FROM alpine:3.20 AS runner

WORKDIR /app

# Cài đặt chứng chỉ SSL và wget cho healthcheck
RUN apk add --no-cache ca-certificates tzdata wget

# Tạo người dùng non-root (Tuân thủ 100% Kyverno Admission Control)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy file thực thi duy nhất từ builder
COPY --from=builder /app/server /app/server

# Chuyển sang người dùng non-root
USER appuser

EXPOSE 5000

ENV PORT=5000
ENV GIN_MODE=release

HEALTHCHECK --interval=20s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:5000/healthz || exit 1

ENTRYPOINT ["/app/server"]
