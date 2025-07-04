# User API Spec

## Base URL

`/users`

---

## 📌 Create User

**POST** `/users`

### Headers

- `Content-Type: application/json`

### Request Body

```json
{
  "name": "John Doe",
  "password": "secret123",
  "role": "admin"
}
```

### Validation Rules

| Field    | Type   | Validation                          |
| -------- | ------ | ----------------------------------- |
| name     | string | required                            |
| password | string | required                            |
| role     | string | required, must be `admin` or `user` |

### Success Response

**Status**: `201 Created`

```json
{
  "id": 1,
  "name": "John Doe",
  "role": "admin"
}
```

### Validation Error

**Status**: `400 Bad Request`

```json
{
  "message": "Validation failed",
  "errors": {
    "Password": "Field Password required"
  }
}
```

---

## 📌 Get All Users

**GET** `/users`

### Success Response

**Status**: `200 OK`

```json
[
  {
    "id": 1,
    "name": "John Doe",
    "role": "admin"
  },
  {
    "id": 2,
    "name": "Jane Smith",
    "role": "user"
  }
]
```

---

## 📌 Get User by ID

**GET** `/users/:id`

### Success Response

**Status**: `200 OK`

```json
{
  "id": 1,
  "name": "John Doe",
  "role": "admin"
}
```

### Not Found

**Status**: `404 Not Found`

```json
{
  "message": "User not found"
}
```

---

## 📌 Update User

**PUT** `/users/:id`

### Request Body

```json
{
  "name": "New Name",
  "password": "newpassword123",
  "role": "user"
}
```

### Success Response

**Status**: `200 OK`

```json
{
  "id": 1,
  "name": "New Name",
  "role": "user"
}
```

---

## 📌 Delete User

**DELETE** `/users/:id`

### Success Response

**Status**: `200 OK`

```json
{
  "message": "User deleted successfully"
}
```

---

## 🔐 Notes

- Passwords are hashed using bcrypt before storage.
- Passwords are not returned in any API response.
- Role must be either `admin` or `user`.
