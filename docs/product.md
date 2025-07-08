# 📦 Product API Documentation

Base URL: `/product`

---

## 📌 Create Product

**POST** `/product`

### Request Body

```json
{
  "name": "PS5 Console",
  "sku": "PS5001",
  "price": 7500000,
  "stock": 5,
  "category_id": 1
}
```

### Validation

| Field       | Type   | Rules                    |
| ----------- | ------ | ------------------------ |
| name        | string | required, min: 3         |
| sku         | string | required                 |
| price       | float  | required, greater than 0 |
| stock       | int    | required, >= 0           |
| category_id | uint   | required                 |

### Success Response

**Status**: `201 Created`

```json
{
  "message": "Product created successfully",
  "status": 201
}
```

### Validation Error

**Status**: `400 Bad Request`

```json
{
  "message": "Validation error",
  "errors": {
    "Price": "gt"
  }
}
```

---

## 📌 Get All Products

**GET** `/product`

### Query Params

- `page` (default: 1)
- `limit` (default: 10)
- `keyword` (optional): filter by name

### Success Response

**Status**: `200 OK`

```json
{
  "status": 200,
  "message": "Products found",
  "data": [
    {
      "id": 1,
      "name": "PS5 Console",
      "sku": "PS5001",
      "price": 7500000,
      "stock": 5,
      "category_id": 1
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

## 📌 Get Product by ID

**GET** `/product/:id`

### Success Response

**Status**: `200 OK`

```json
{
  "status": 200,
  "message": "Product found",
  "data": {
    "id": 1,
    "name": "PS5 Console",
    "sku": "PS5001",
    "price": 7500000,
    "stock": 5,
    "category_id": 1
  }
}
```

### Not Found

**Status**: `404 Not Found`

```json
{
  "message": "Product not found"
}
```

---

## 📌 Update Product

**PATCH** `/product/:id`

### Request Body

```json
{
  "name": "PS5 Slim",
  "sku": "PS5002",
  "price": 7800000,
  "stock": 7,
  "category_id": 1
}
```

### Success Response

**Status**: `200 OK`

```json
{
  "message": "Product updated successfully",
  "status": 200
}
```

---

## 📌 Delete Product

**DELETE** `/product/:id`

### Success Response

**Status**: `200 OK`

```json
{
  "message": "Product deleted successfully"
}
```

---

## 🔐 Notes

- All endpoints are protected with authentication and module access middleware.
- Validation uses `validator.v10`.
- Responses follow the global `dto.Response` structure.
- `category_id` must refer to an existing category.
- The `category` relation is not returned unless explicitly preloaded (e.g., for admin panel).
