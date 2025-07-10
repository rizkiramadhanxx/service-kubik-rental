
# 📦 Package API Documentation

**Base URL:** `/package`

Semua endpoint di bawah ini **memerlukan autentikasi JWT** dan akses module `package`.

---

## 📌 Create Package

**POST** `/package`

### Request Body

```json
{
  "name": "Paket 1 Jam",
  "duration": 60,
  "price": 10000
}
```

### Response

```json
{
  "status": 201,
  "message": "Package created successfully",
  "data": {
    "id": 1,
    "name": "Paket 1 Jam",
    "duration": 60,
    "price": 10000,
    "created_at": 1720000000
  }
}
```

---

## 📌 Get All Packages (Pagination + Search)

**GET** `/package`

### Query Parameters

| Parameter | Type   | Description                      |
|-----------|--------|----------------------------------|
| `page`    | int    | Page number (default: 1)         |
| `limit`   | int    | Items per page (default: 10)     |
| `keyword` | string | Search by name (optional)        |

### Example

```
GET /package?page=1&limit=5&keyword=1%20Jam
```

### Response

```json
{
  "status": 200,
  "message": "Packages retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Paket 1 Jam",
      "duration": 60,
      "price": 10000,
      "created_at": 1720000000,
      "updated_at": 1720000000
    }
  ],
  "meta": {
    "page": 1,
    "limit": 5,
    "total": 1,
    "total_page": 1
  }
}
```

---

## 📌 Get Package by ID

**GET** `/package/:id`

### Example

```
GET /package/1
```

### Response

```json
{
  "status": 200,
  "message": "Package found",
  "data": {
    "id": 1,
    "name": "Paket 1 Jam",
    "duration": 60,
    "price": 10000,
    "created_at": 1720000000,
    "updated_at": 1720000000
  }
}
```

### Error (Not Found)

```json
{
  "status": 404,
  "message": "Package not found"
}
```

---

## 📌 Update Package

**PATCH** `/package/:id`

### Request Body

```json
{
  "name": "Paket 2 Jam",
  "duration": 120,
  "price": 18000
}
```

### Response

```json
{
  "status": 200,
  "message": "Package updated successfully",
  "data": {
    "id": 1,
    "name": "Paket 2 Jam",
    "duration": 120,
    "price": 18000,
    "created_at": 1720000000,
    "updated_at": 1720000500
  }
}
```

---

## 📌 Delete Package

**DELETE** `/package/:id`

### Example

```
DELETE /package/1
```

### Response

```json
{
  "status": 200,
  "message": "Package deleted successfully"
}
```

### Error (Not Found)

```json
{
  "message": "Package not found"
}
```

---

## 🛡️ Auth & Middleware

Semua endpoint di atas:

- **Memerlukan JWT Bearer Token** melalui header:
  
  ```
  Authorization: Bearer <token>
  ```

- **Dibatasi oleh middleware akses modul**:
  
  ```
  middleware.RequireModuleAccess("package")
  ```
