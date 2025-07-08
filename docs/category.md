# 📦 Category API Documentation

Base URL: `/category`

---

## 📌 Create Category

**POST** `/category`

### Request Body

```json
{
  "name": "Electronics"
}
```

### Validation

| Field | Type   | Rules            |
| ----- | ------ | ---------------- |
| name  | string | required, min: 4 |

### Success Response

**Status**: `201 Created`

```json
{
  "message": "Category created successfully",
  "status": 201
}
```

### Validation Error

**Status**: `400 Bad Request`

```json
{
  "message": "Validation error",
  "errors": {
    "Name": "min"
  }
}
```

---

## 📌 Get All Categories

**GET** `/category`

### Query Params

- `page` (default: 1)
- `limit` (default: 10)
- `keyword` (optional): filter by name

### Success Response

**Status**: `200 OK`

```json
{
  "status": 200,
  "message": "Categories found",
  "data": [
    {
      "id": 1,
      "name": "Electronics"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "totalPage": 1
  }
}
```

---

## 📌 Get Category by ID

**GET** `/category/:id`

### Success Response

**Status**: `200 OK`

```json
{
  "status": 200,
  "message": "Category found",
  "data": {
    "id": 1,
    "name": "Electronics"
  }
}
```

### Not Found

**Status**: `404 Not Found`

```json
{
  "message": "Category not found"
}
```

---

## 📌 Update Category

**PATCH** `/category/:id`

### Request Body

```json
{
  "name": "Home Appliances"
}
```

### Success Response

**Status**: `200 OK`

```json
{
  "message": "Category updated successfully",
  "status": 200
}
```

---

## 📌 Delete Category

**DELETE** `/category/:id`

### Success Response

**Status**: `200 OK`

```json
{
  "message": "Category deleted successfully"
}
```

---

## 🔐 Notes

- All endpoints are protected with authentication and module access middleware.
- Validation uses `validator.v10`, errors are formatted using `pkg.FormatValidationError`.
- The `products` field in `Category` entity is only included if explicitly preloaded.
