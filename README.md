# cuahang-be — High Performance Go Backend

Dịch vụ Backend REST API cho dự án mẫu DevOps **`cuahang`**, viết bằng ngôn ngữ **Golang 1.23** và **GORM**, tuân thủ cấu trúc thư mục chuẩn **Standard Go Project Layout**.

## Cấu trúc thư mục (Clean Architecture):
```text
cuahang-be/
├── cmd/
│   └── server/
│       └── main.go              # Entrypoint khởi chạy HTTP Server
├── internal/
│   ├── config/
│   │   └── config.go            # Quản lý cấu hình & biến môi trường
│   ├── database/
│   │   ├── db.go                # Khởi tạo kết nối PostgreSQL & AutoMigrate
│   │   └── seed.go              # Khởi tạo dữ liệu mẫu ban đầu
│   ├── models/
│   │   ├── product.go           # Model cấu trúc bảng Product
│   │   └── order.go             # Model cấu trúc bảng Order
│   ├── handlers/
│   │   ├── health.go            # Kubernetes Probe (/healthz, /api/info)
│   │   ├── product.go           # CRUD sản phẩm (/api/products)
│   │   └── order.go             # Đặt hàng (/api/orders)
│   └── middleware/
│       └── cors.go              # Middleware CORS cho Web Frontend
├── test/
│   └── pricing_test.go          # Unit tests cho CI pipeline
├── Dockerfile                   # Multi-stage build (Image ~20 MB)
├── Jenkinsfile                  # CI/CD pipeline chuẩn DevSecOps
└── go.mod                       # Quản lý Golang dependencies
```

## Chạy thử nghiệm với Docker:
```bash
# 1. Build Docker image
docker build -t cuahang-be:latest .

# 2. Chạy container kết nối vào Database VM (192.168.1.253)
docker run -p 5000:5000 \
  -e DATABASE_URL="host=192.168.1.253 user=app_cuahang_dev password=DevPass@2026 dbname=cuahang_dev port=5432 sslmode=disable" \
  cuahang-be:latest
```
