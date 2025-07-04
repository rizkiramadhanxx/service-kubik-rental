# Role API Spec

## 🎯 Base URL

`/roles`

---

## ✅ Module Enum (Allowed Modules)

```json
["user", "setting", "role", "transaksi"]
```

`modules` field must be a JSON string of array of strings, containing only allowed module names above.

---

## 📌 Create Role

**POST** `/roles`

### Request Headers

- `Content-Type: application/json`

### Request Body

```json
{
  "name": "admin",
  "modules": "[\"user\", \"setting\", \"role\", \"transaksi\"]"
}
```

### Success Response

**201 Created**

```json
{
  "id": 1,
  "name": "admin",
  "modules": "[\"user\", \"setting\", \"role\", \"transaksi\"]"
}
```

---

## 📌 Get All Roles

**GET** `/roles`

### Success Response

**200 OK**

```json
[
  {
    "id": 1,
    "name": "admin",
    "modules": "[\"user\", \"setting\", \"role\", \"transaksi\"]"
  },
  ...
]
```

---

## 📌 Get Role by ID

**GET** `/roles/:id`

### Success Response

**200 OK**

```json
{
  "id": 1,
  "name": "admin",
  "modules": "[\"user\", \"setting\", \"role\", \"transaksi\"]"
}
```

---

## 📌 Update Role

**PUT** `/roles/:id`

### Request Body

```json
{
  "name": "manager",
  "modules": "[\"user\", \"transaksi\"]"
}
```

### Success Response

**200 OK**

```json
{
  "id": 1,
  "name": "manager",
  "modules": "[\"user\", \"transaksi\"]"
}
```

---

## 📌 Delete Role

**DELETE** `/roles/:id`

### Success Response

**200 OK**

```json
{
  "message": "Role deleted successfully"
}
```

---

## ❌ Validation Errors

- Invalid JSON in `modules`:

```json
{
  "message": "Invalid modules format"
}
```

- Invalid module enum:

```json
{
  "message": "Invalid module: unknown"
}
```
