## OTP Auth Service (Go + Gin + Postgres + Redis)

این سرویس، لاگین با کد یک‌بارمصرف (OTP) را پیاده‌سازی می‌کند و همچنین لیست و جزییات کاربر را با احراز هویت JWT ارائه می‌دهد. OTP و محدودسازی نرخ (Rate Limit) روی Redis انجام می‌شود و اطلاعات کاربران در Postgres ذخیره می‌گردد.

### ویژگی‌ها
- ایجاد و ذخیره OTP در Redis با انقضا 2 دقیقه (`otp:<phone>`)
- محدودسازی نرخ درخواست OTP برای هر شماره: حداکثر 3 درخواست در 10 دقیقه (`otp:rate:<phone>`)
- صدور JWT پس از تایید OTP
- Swagger UI در مسیر `/swagger/index.html`

---

## اجرا در حالت Local

1) پیش‌نیازها:
- نصب Go (نسخه‌ای سازگار با Dockerfile: 1.24.x)
- اجرای Postgres و Redis لوکال یا تنظیم آنها در `.env`

2) فایل `.env` بسازید (نمونه):
```env
# App
APP_PORT=8080

# Postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=changeme
POSTGRES_DB=postgres_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=supersecret
JWT_EXPIRE_IN=60
```

3) اجرای سرویس:
```bash
go run ./cmd
```

سرویس روی `http://localhost:8080` بالا می‌آید. Swagger در `http://localhost:8080/swagger/index.html`.

---

## اجرا با Docker / Docker Compose

1) فایل `.env` در ریشه پروژه بسازید (برای کانتینر اپ، هاست‌های پایگاه‌داده باید به سرویس‌های Compose اشاره کنند):
```env
APP_PORT=8080

# Postgres داخل Compose با نام سرویس "db" بالا می‌آید
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=changeme
POSTGRES_DB=postgres_db

# Redis داخل Compose با نام سرویس "redis" بالا می‌آید
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=supersecret
JWT_EXPIRE_IN=60
```

2) اجرای Compose:
```bash
docker-compose up --build
```

سرویس اپ روی پورت `8080` میزبان در دسترس است. Postgres روی `5432` و Redis روی `6379` نیز اکسپوز می‌شوند.

---

## API ها و نمونه درخواست/پاسخ

پیشوند مسیرها مطابق `cmd/main.go`:
- Auth: `/api/v1`
- Users (محافظت‌شده با JWT): `/api/v2`

Swagger: `GET /swagger/index.html`

### 1) درخواست OTP
- مسیر: `POST /api/v1/login`
- بدنه:
```json
{ "phone": 989124567890 }
```
- پاسخ موفق (200):
```json
{ "message": "OTP sent successfully", "otp": "123456" }
```
- محدودسازی نرخ: در صورت عبور از 3 درخواست در 10 دقیقه، پاسخ 429:
```json
{ "error": "too many OTP requests, try again later" }
```

نمونه curl:
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"phone":989124567890}'
```

### 2) تایید OTP و دریافت توکن
- مسیر: `POST /api/v1/verify`
- بدنه:
```json
{ "phone": 989124567890, "otp": "123456" }
```
- پاسخ موفق (200):
```json
{ "token": "<JWT_TOKEN>" }
```

نمونه curl:
```bash
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{"phone":989124567890, "otp":"123456"}'
```

### 3) دریافت کاربر با شناسه (Protected)
- مسیر: `GET /api/v2/users/{id}`
- هدر: `Authorization: Bearer <JWT_TOKEN>`
- پاسخ موفق (200):
```json
{ "id": 1, "phone": "989124567890", "created_at": "2025-09-17T16:00:00Z" }
```

نمونه curl:
```bash
curl -X GET http://localhost:8080/api/v2/users/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### 4) لیست کاربران (Protected)
- مسیر: `GET /api/v2/users`
- پارامترها: `page` (پیش‌فرض 1) و `limit` (پیش‌فرض 10)
- هدر: `Authorization: Bearer <JWT_TOKEN>`
- پاسخ موفق (200):
```json
{
  "page": 1,
  "limit": 10,
  "total": 1,
  "users": [
    { "id": 1, "phone": "989124567890", "created_at": "2025-09-17T16:00:00Z" }
  ]
}
```

نمونه curl:
```bash
curl -X GET "http://localhost:8080/api/v2/users?page=1&limit=10" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

## معماری و انتخاب دیتابیس

- Postgres:
  - نگهداری موجودیت‌های پایدار (جدول `users`) با فیلدهای `id`, `phone_number` و زمان‌ها.
  - مهاجرت خودکار مدل‌ها در `internal/database/db.go`.

- Redis:
  - نگهداری داده‌های گذرا و زمان‌مند.
  - OTP با کلید `otp:<phone>` و TTL دو دقیقه ذخیره می‌شود (`SaveOTP`).
  - Rate Limit با کلید `otp:rate:<phone>` و TTL ده دقیقه؛ با هر درخواست `INCR` می‌شود و فقط بار اول `EXPIRE` ست می‌گردد.

- احراز هویت:
  - تایید OTP → صدور JWT با `internal/auth/jwt/jwt.go`.
  - دسترسی به مسیرهای `/api/v2/**` از طریق میان‌افزار `middleware.AuthModdleware()` محافظت می‌شود.

---

## محدودسازی نرخ (Rate Limit)

- پیاده‌سازی در `internal/auth/auth_repository.go`:
  - `IncrementOtpRequestCount(phone, window, limit)`
  - کلید: `otp:rate:<phone>`
  - عملیات: `INCR` تعداد را زیاد می‌کند؛ اگر مقدار 1 شد، `EXPIRE` با مدت `window` ست می‌شود.
  - اگر مقدار از `limit` عبور کند، `ErrRateLimitExceeded` برگردانده می‌شود.

- مصرف در `internal/auth/auth_service.go`، ابتدای `LoginUser`:
  - `IncrementOtpRequestCount(phone, 10*time.Minute, 3)`
  - در صورت خطای حد، پیام کاربرپسند برمی‌گرداند تا هندلر 429 دهد.

- هندلینگ در `internal/auth/auth_handler.go`:
  - در `Login`، اگر پیام «too many OTP requests...» بود، `HTTP 429` برمی‌گرداند.
  - Swagger برای 429 مستند شده است.

---

## Swagger

- تولید مستندات:
```bash
swag init -g ./cmd/main.go
```

- فایل‌های خروجی در `./docs` ایجاد می‌شود و در `cmd/main.go` با `import _ "deca-task/docs"` بارگذاری می‌شود.

اگر `swag` نصب نیست:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

---

## ساختار پروژه (High-level)

```
cmd/
  main.go                 # راه‌اندازی سرور و مسیرها
internal/
  auth/
    auth_handler.go       # هندلرهای Auth (login/verify)
    auth_service.go       # منطق Auth + rate limit
    auth_repository.go    # ذخیره/خواندن OTP و rate limit روی Redis، کاربر روی DB
    dto/                  # DTOهای Auth
    jwt/                  # صدور توکن JWT
  user/
    user_handler.go       # هندلرهای User (نیازمند JWT)
    user_service.go
    user_repository.go
    dto/                  # DTOهای User
  models/                 # مدل‌های GORM (User)
  database/               # اتصال Postgres/Redis
  config/                 # بارگذاری تنظیمات از .env
docs/                     # خروجی‌های Swagger
```

---
