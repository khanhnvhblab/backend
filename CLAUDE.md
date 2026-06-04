# CLAUDE.md — Backend (Golang)

> Tài liệu thiết kế đầy đủ: [DESIGN.md](DESIGN.md)

---

## Tech Stack

| Component | Technology |
|---|---|
| Language | Golang 1.22+ |
| Framework | Gin |
| Database | MongoDB (`mongodb://localhost:27017`, db: `todolist`) |
| Auth | JWT — access token (1h) + refresh token (7d, httpOnly cookie) |
| Password | bcrypt cost 12 |

---

## Project Structure

```
backend/
├── main.go
├── config/config.go               # Đọc env vars
├── internal/
│   ├── handler/                   # Nhận request → gọi service → trả response
│   ├── service/                   # Business logic
│   ├── repository/                # MongoDB queries
│   ├── model/                     # MongoDB document structs
│   ├── dto/                       # Request / Response structs
│   ├── middleware/                # auth_middleware.go, cors_middleware.go
│   └── router/router.go           # Đăng ký routes
└── tests/
    ├── handler/
    ├── service/
    └── repository/
```

---

## Conventions

- Pattern bắt buộc: `handler → service → repository` — không bỏ qua tầng nào
- Đặt tên file: `<domain>_handler.go`, `<domain>_service.go`, `<domain>_repository.go`
- **Mọi MongoDB query phải filter theo `user_id`** — data isolation là bắt buộc
- Validate input bằng Gin binding tag (`binding:"required,email"`, v.v.)
- Không comment code hiển nhiên — chỉ khi logic phức tạp hoặc có lý do đặc biệt

### Response format thống nhất

```json
// Thành công có object
{ "data": { ... } }

// Thành công có danh sách
{ "data": [...], "total": 42, "page": 1, "limit": 20 }

// Lỗi
{ "error": "message mô tả lỗi" }
```

---

## Data Models

### users
```json
{ "_id", "email" (unique), "password" (bcrypt), "name", "created_at", "updated_at" }
```

### categories
```json
{ "_id", "user_id", "name", "color" (hex, default #6366f1), "created_at", "updated_at" }
```

### todos
```json
{
  "_id", "user_id", "category_id" (nullable),
  "title", "description",
  "status": "pending | in_progress | done",
  "priority": "low | medium | high",
  "deadline" (nullable), "completed_at" (nullable),
  "created_at", "updated_at"
}
```

### MongoDB Indexes
| Collection | Index | Type |
|---|---|---|
| `users` | `email` | unique |
| `todos` | `(user_id, status)` | compound |
| `todos` | `deadline` | single |
| `categories` | `user_id` | single |

---

## API — Base URL: `/api/v1`

### Auth (không cần JWT)
| Method | Path | Description |
|---|---|---|
| POST | `/auth/register` | Đăng ký |
| POST | `/auth/login` | Đăng nhập → trả access + refresh token |
| POST | `/auth/refresh` | Lấy access token mới từ refresh token |
| POST | `/auth/logout` | Đăng xuất (cần JWT) |

### Todos (cần JWT)
| Method | Path | Description |
|---|---|---|
| GET | `/todos` | Danh sách (filter, search, sort, paginate) |
| POST | `/todos` | Tạo mới |
| GET | `/todos/:id` | Chi tiết |
| PUT | `/todos/:id` | Cập nhật toàn bộ |
| PATCH | `/todos/:id/status` | Chỉ cập nhật status |
| DELETE | `/todos/:id` | Xóa |

Query params cho `GET /todos`: `status`, `priority`, `category_id`, `search`, `sort_by` (`created_at`|`deadline`), `order` (`asc`|`desc`), `page` (default 1), `limit` (default 20, max 100).

### Categories (cần JWT)
| Method | Path | Description |
|---|---|---|
| GET | `/categories` | Danh sách của user |
| POST | `/categories` | Tạo mới |
| PUT | `/categories/:id` | Cập nhật |
| DELETE | `/categories/:id` | Xóa |

### Dashboard (cần JWT)
| Method | Path | Description |
|---|---|---|
| GET | `/dashboard/stats` | Tổng: total, pending, in_progress, done, overdue, due_soon |

> `overdue`: deadline đã qua + status != done  
> `due_soon`: deadline trong 7 ngày tới + status != done

---

## Security

- Rate limiting: 100 req/min per IP trên `/auth/*`
- CORS: chỉ cho phép origin của frontend
- JWT gắn `user_id` vào Gin context — mọi handler lấy từ đó, không nhận từ request body

---

## Environment Variables

```env
APP_ENV=development
APP_PORT=8080

JWT_SECRET=<strong-random-secret>
JWT_ACCESS_TTL=3600
JWT_REFRESH_TTL=604800

MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=todolist
```

> **KHÔNG commit `.env`** — chỉ commit `.env.example`
