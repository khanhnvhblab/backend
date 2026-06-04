# Backend Design — TodoList

## 1. Tech Stack

| Component | Technology |
|---|---|
| Language | Golang 1.22+ |
| Framework | Gin |
| Database | MongoDB |
| Auth | JWT (access + refresh token) |
| Password | bcrypt (cost 12) |

---

## 2. Project Structure

```
backend/
├── main.go
├── go.mod
├── go.sum
├── .env
├── config/
│   └── config.go                   # Đọc env vars, khởi tạo config struct
├── internal/
│   ├── handler/                    # Nhận request, gọi service, trả response
│   │   ├── auth_handler.go
│   │   ├── todo_handler.go
│   │   ├── category_handler.go
│   │   └── dashboard_handler.go
│   ├── service/                    # Business logic
│   │   ├── auth_service.go
│   │   ├── todo_service.go
│   │   └── category_service.go
│   ├── repository/                 # Tất cả truy vấn MongoDB
│   │   ├── user_repository.go
│   │   ├── todo_repository.go
│   │   └── category_repository.go
│   ├── model/                      # MongoDB document structs
│   │   ├── user.go
│   │   ├── todo.go
│   │   └── category.go
│   ├── dto/                        # Request / Response structs
│   │   ├── auth_dto.go
│   │   ├── todo_dto.go
│   │   └── category_dto.go
│   ├── middleware/
│   │   ├── auth_middleware.go      # Xác thực JWT, gắn user_id vào context
│   │   └── cors_middleware.go
│   └── router/
│       └── router.go               # Đăng ký tất cả routes
└── tests/
    ├── handler/
    ├── service/
    └── repository/
```

---

## 3. Data Models (MongoDB)

### Collection: `users`
```json
{
  "_id": "ObjectId",
  "email": "string (unique, required)",
  "password": "string (bcrypt hashed)",
  "name": "string (required)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Collection: `categories`
```json
{
  "_id": "ObjectId",
  "user_id": "ObjectId (ref: users)",
  "name": "string (required)",
  "color": "string (hex, default: #6366f1)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Collection: `todos`
```json
{
  "_id": "ObjectId",
  "user_id": "ObjectId (ref: users)",
  "category_id": "ObjectId (ref: categories, nullable)",
  "title": "string (required)",
  "description": "string",
  "status": "enum: pending | in_progress | done",
  "priority": "enum: low | medium | high",
  "deadline": "datetime (nullable)",
  "completed_at": "datetime (nullable)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Indexes
| Collection | Index | Type |
|---|---|---|
| `users` | `email` | unique |
| `todos` | `(user_id, status)` | compound |
| `todos` | `deadline` | single |
| `categories` | `user_id` | single |

---

## 4. API Specification

**Base URL:** `/api/v1`  
**Auth:** Header `Authorization: Bearer <access_token>`

**Response format thống nhất:**
```json
// Thành công có data
{ "data": { ... } }

// Thành công có danh sách
{ "data": [...], "total": 42, "page": 1, "limit": 20 }

// Lỗi
{ "error": "message mô tả lỗi" }
```

---

### 4.1 Auth Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/auth/register` | No | Đăng ký tài khoản mới |
| POST | `/auth/login` | No | Đăng nhập |
| POST | `/auth/refresh` | No | Lấy access token mới |
| POST | `/auth/logout` | Yes | Đăng xuất |

#### POST /auth/register
```json
// Request
{
  "email": "user@example.com",
  "password": "Abc@12345",
  "name": "Nguyen Van A"
}

// Response 201
{
  "data": {
    "id": "664abc...",
    "email": "user@example.com",
    "name": "Nguyen Van A"
  }
}
```

#### POST /auth/login
```json
// Request
{
  "email": "user@example.com",
  "password": "Abc@12345"
}

// Response 200
{
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "expires_in": 3600
  }
}
```

#### POST /auth/refresh
```json
// Request
{ "refresh_token": "eyJhbGci..." }

// Response 200
{
  "data": {
    "access_token": "eyJhbGci...",
    "expires_in": 3600
  }
}
```

---

### 4.2 Todo Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/todos` | Yes | Danh sách todo (filter + search + sort + paginate) |
| POST | `/todos` | Yes | Tạo todo mới |
| GET | `/todos/:id` | Yes | Chi tiết todo |
| PUT | `/todos/:id` | Yes | Cập nhật toàn bộ todo |
| PATCH | `/todos/:id/status` | Yes | Chỉ cập nhật trạng thái |
| DELETE | `/todos/:id` | Yes | Xóa todo |

#### GET /todos — Query Parameters
| Param | Type | Default | Description |
|---|---|---|---|
| `status` | string | — | `pending` \| `in_progress` \| `done` |
| `priority` | string | — | `low` \| `medium` \| `high` |
| `category_id` | string | — | Filter theo category |
| `search` | string | — | Tìm trong `title` (case-insensitive) |
| `sort_by` | string | `created_at` | `created_at` \| `deadline` |
| `order` | string | `desc` | `asc` \| `desc` |
| `page` | int | `1` | Số trang |
| `limit` | int | `20` | Số item/trang (max: 100) |

#### POST /todos
```json
// Request
{
  "title": "Hoàn thiện báo cáo Q2",
  "description": "Bao gồm số liệu tháng 4-6",
  "priority": "high",
  "category_id": "664abc123...",
  "deadline": "2026-06-30T17:00:00Z"
}

// Response 201
{ "data": { "id": "...", "title": "...", "status": "pending", ... } }
```

#### PATCH /todos/:id/status
```json
// Request
{ "status": "done" }

// Response 200
{ "data": { "id": "...", "status": "done", "completed_at": "2026-06-02T10:00:00Z" } }
```

---

### 4.3 Category Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/categories` | Yes | Danh sách category của user |
| POST | `/categories` | Yes | Tạo category mới |
| PUT | `/categories/:id` | Yes | Cập nhật category |
| DELETE | `/categories/:id` | Yes | Xóa category |

#### POST /categories
```json
// Request
{ "name": "Công việc", "color": "#ef4444" }

// Response 201
{ "data": { "id": "...", "name": "Công việc", "color": "#ef4444" } }
```

---

### 4.4 Dashboard Endpoint

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/dashboard/stats` | Yes | Thống kê tổng quan |

#### GET /dashboard/stats
```json
// Response 200
{
  "data": {
    "total": 42,
    "pending": 15,
    "in_progress": 8,
    "done": 19,
    "overdue": 3,
    "due_soon": 5
  }
}
```
> `overdue`: deadline đã qua, status != done  
> `due_soon`: deadline trong 7 ngày tới, status != done

---

## 5. Security

| Hạng mục | Giải pháp |
|---|---|
| Password | bcrypt cost factor 12 |
| Access token | JWT, TTL 1 giờ |
| Refresh token | JWT, TTL 7 ngày, lưu trong httpOnly cookie |
| CORS | Chỉ cho phép origin của frontend |
| Rate limiting | 100 req/min per IP trên `/auth/*` |
| Input validation | Gin binding tags (`binding:"required,email"`) |
| Data isolation | Mọi query đều filter theo `user_id` lấy từ JWT context |

---

## 6. Environment Variables

```env
APP_ENV=development
APP_PORT=8080

JWT_SECRET=<strong-random-secret>
JWT_ACCESS_TTL=3600
JWT_REFRESH_TTL=604800

MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=todolist
```
