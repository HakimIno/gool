# Gool - Go Project Scaffolding Tool

Gool เป็น CLI tool ที่ทันสมัย ยืดหยุ่น และครอบคลุมสำหรับการสร้างโปรเจค Go พร้วมโครงสร้างที่ปรับแต่งได้ ไฟล์ที่ตั้งค่าล่วงหน้า และ dependencies สร้างแอปพลิเคชัน Go ที่พร้อมใช้งานในระดับ production ได้ในไม่กี่วินาที!

## 🚀 คุณสมบัติ

### คุณสมบัติหลัก
- **Web Frameworks หลากหลาย**: เลือกจาก Gin, Echo, Fiber, หรือ Revel
- **รองรับฐานข้อมูล**: PostgreSQL, MySQL, SQLite, MongoDB, Redis, หรือ in-memory store
- **ORM/Database Access**: GORM, sqlx, raw SQL, หรือไม่ใช้
- **Architecture Patterns**: Simple, Clean Architecture, Hexagonal (Ports & Adapters), MVC, หรือ Custom
- **รูปแบบ Configuration**: รองรับ YAML, JSON, หรือ TOML พร้อม environment-specific configs
- **Authentication**: JWT, OAuth2, หรือ Basic Auth พร้อม templates ที่พร้อมใช้
- **Middleware**: CORS, Rate Limiting, Logging, และ Authentication middleware
- **Testing**: Unit test และ integration test templates
- **Logging & Monitoring**: Standard log, Logrus, หรือ Zap พร้อม Prometheus metrics

### คุณสมบัติเพิ่มเติม
- **WebSocket Support**: Templates สำหรับแอปพลิเคชัน real-time
- **Error Handling**: การจัดการ error แบบ centralized พร้อม custom error types
- **Caching**: รองรับ In-memory, Redis, หรือ Memcached
- **Message Queues**: RabbitMQ, Kafka, หรือ NATS สำหรับ async tasks
- **Security**: HTTPS, secure headers (HSTS, CSP), และ CSRF protection
- **API Documentation**: การสร้าง OpenAPI/Swagger
- **Docker & Deployment**: Dockerfile และ docker-compose.yml
- **CI/CD**: Templates สำหรับ GitHub Actions และ GitLab CI
- **Health Checks**: Built-in /health endpoint
- **Internationalization**: รองรับหลายภาษา
- **Cloud Integration**: การตั้งค่าสำหรับ deployment บน AWS, GCP, และ Azure

## 📦 การติดตั้ง

### วิธีที่ 1: ติดตั้งจาก source code
```bash
# Clone repository
git clone https://github.com/HakimIno/gool.git
cd gool

# Build binary
go build -o gool main.go

# ติดตั้งไปยัง system PATH (macOS/Linux)
sudo mv gool /usr/local/bin/

# หรือสำหรับ Windows
# ย้าย gool.exe ไปยัง directory ที่อยู่ใน PATH
```

### วิธีที่ 2: Go Install (เมื่อ publish แล้ว)
```bash
go install github.com/HakimIno/gool@latest
```

### วิธีที่ 3: Download binary (เมื่อมี releases)
Download binary ล่าสุดจาก [releases page](https://github.com/HakimIno/gool/releases)

## 🛠️ การ Setup โปรเจคใหม่

### สำหรับผู้พัฒนา (การตั้งค่า Git Repository)

หากคุณต้องการสร้าง repository ของตัวเอง:

```bash
# 1. สร้าง repository ใหม่บน GitHub
# 2. Clone หรือ setup local repository
git init
git add .
git commit -m "Initial commit"

# 3. เพิ่ม remote repository
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO_NAME.git

# 4. Push ขึ้น GitHub
git branch -M main
git push -u origin main
```

### การ build และทดสอบ

```bash
# ติดตั้ง dependencies
go mod tidy

# Build binary
go build -o gool main.go

# ทดสอบ command
./gool --help
./gool version
```

## 🎯 การเริ่มต้นใช้งาน

### Interactive Mode (แนะนำ)
```bash
gool init my-awesome-app
```

จะมีคำถามให้คุณตอบเพื่อปรับแต่งโปรเจค:
- ชื่อโปรเจคและ module path
- เลือก Web framework
- ตั้งค่าฐานข้อมูลและ ORM
- เลือก Architecture pattern
- วิธี Authentication
- คุณสมบัติและ middleware เพิ่มเติม

### Non-Interactive Mode
```bash
gool init my-app --framework=gin --orm=gorm --database=postgresql --arch=simple
```

### ตัวเลือกที่มี
```bash
# Framework options
--framework=gin|echo|fiber|revel

# ORM options  
--orm=gorm|sqlx|raw|none

# Database options
--database=postgresql|mysql|sqlite|mongodb|redis|memory

# Architecture options
--arch=simple|clean|hexagonal|mvc|custom
```

## 📂 โครงสร้างโปรเจคที่สร้าง

### Simple Architecture
```
my-app/
├── cmd/                    # Application entrypoints
├── internal/
│   ├── app/               # Application initialization
│   ├── handlers/          # HTTP handlers
│   ├── models/            # Data models
│   ├── services/          # Business logic
│   └── middleware/        # Custom middleware
├── pkg/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection
│   └── logger/            # Logging utilities
├── api/
│   └── routes/            # Route definitions
├── docs/                   # Documentation
├── scripts/               # Build and deployment scripts
├── deployments/           # Deployment configs
├── test/                  # Test files
├── static/                # Static assets (if enabled)
├── locales/               # Internationalization files
├── main.go                # Application entry point
├── go.mod                 # Go modules
├── .env                   # Environment variables
├── .gitignore            # Git ignore rules
├── Dockerfile            # Docker configuration
└── README.md             # Project documentation
```

### Clean Architecture
```
my-app/
├── cmd/
├── internal/
│   ├── entity/           # Business entities
│   ├── usecase/          # Business use cases
│   ├── controller/       # Controllers (handlers)
│   ├── repository/       # Data access layer
│   └── infrastructure/   # External services
├── pkg/
└── ...
```

### Hexagonal Architecture
```
my-app/
├── cmd/
├── internal/
│   ├── core/
│   │   ├── domain/       # Domain models
│   │   ├── ports/        # Interfaces
│   │   └── services/     # Domain services
│   └── adapters/
│       ├── primary/      # Driving adapters (HTTP, CLI)
│       └── secondary/    # Driven adapters (Database, Cache)
├── pkg/
└── ...
```

## 🔧 ตัวอย่างการใช้งาน

### สร้าง simple web API
```bash
gool init my-api \
  --framework=gin \
  --orm=gorm \
  --database=postgresql \
  --arch=simple
```

### สร้าง microservice ด้วย clean architecture
```bash
gool init my-service \
  --framework=echo \
  --orm=sqlx \
  --database=mysql \
  --arch=clean
```

### สร้าง full-stack application
```bash
# Interactive mode จะถามเกี่ยวกับ:
# - Static file serving
# - Template engine
# - WebSocket support
# - Authentication
# - และอื่นๆ...
gool init my-fullstack-app
```

## 🏃‍♂️ รันโปรเจคที่สร้างแล้ว

หลังจากสร้างโปรเจคแล้ว:

```bash
cd my-app

# ติดตั้ง dependencies
go mod tidy

# คัดลอกและตั้งค่า environment
cp .env.example .env
# แก้ไข .env ด้วยข้อมูล database credentials

# รันแอปพลิเคชัน
go run main.go

# หรือ build และรัน
go build -o app main.go
./app
```

### ด้วย Docker
```bash
# เริ่ม services (database, redis, etc.)
docker-compose up -d

# Build และรัน app
docker build -t my-app .
docker run -p 8080:8080 my-app
```

## 🧪 การทดสอบ

โปรเจคที่สร้างจะมี test templates:

```bash
# รัน unit tests
go test ./...

# รัน tests พร้อม coverage
go test -cover ./...

# รัน integration tests
go test -tags=integration ./test/...
```

## 📖 เอกสาร API

หากเปิดใช้งาน Swagger documentation, เยี่ยมชม:
```
http://localhost:8080/swagger/index.html
```

## 🎛️ การตั้งค่า

โปรเจคที่สร้างรองรับการตั้งค่าตาม environment:

```bash
# Development
APP_ENV=development go run main.go

# Production  
APP_ENV=production go run main.go

# Test
APP_ENV=test go test ./...
```

## 🔍 Health Checks

Built-in health check endpoint:
```bash
curl http://localhost:8080/api/v1/health
```

## 🤝 การมีส่วนร่วม

เรายินดีรับการมีส่วนร่วม! โปรดดู [Contributing Guide](CONTRIBUTING.md) สำหรับรายละเอียด

## 📝 ใบอนุญาต

โปรเจคนี้ได้รับอนุญาตภายใต้ MIT License - ดูไฟล์ [LICENSE](LICENSE) สำหรับรายละเอียด

## 🙏 กิตติกรรมประกาศ

- [Cobra](https://github.com/spf13/cobra) สำหรับ CLI framework
- [Viper](https://github.com/spf13/viper) สำหรับ configuration management
- [Survey](https://github.com/AlecAivazis/survey) สำหรับ interactive prompts
- ชุมชน Go สำหรับ libraries และ tools ที่ยอดเยี่ยม

## 📞 การสนับสนุน

- 📚 [เอกสาร](https://github.com/HakimIno/gool/wiki)
- 🐛 [Issue Tracker](https://github.com/HakimIno/gool/issues)
- 💬 [Discussions](https://github.com/HakimIno/gool/discussions)

---

สร้างด้วย ❤️ โดยทีม Gool
