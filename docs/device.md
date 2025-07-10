
# 📱 Device API Documentation

**Base URL:** `/device`

---

## 📌 Get All Devices (with Pagination & Search)

**GET** `/device`

### Query Parameters

| Parameter | Type   | Description                            |
|-----------|--------|----------------------------------------|
| `page`    | int    | Page number (default: 1)               |
| `limit`   | int    | Number of items per page (default: 10) |
| `keyword` | string | Filter by name or IP (optional)        |

### Response

```json
{
  "status": 200,
  "message": "Success get all devices",
  "data": [
    {
      "id": 1,
      "name": "TV 1",
      "ip": "192.168.1.10",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_page": 1
  }
}
```

---

## 📌 Get Device by ID

**GET** `/device/:id`

### Response

```json
{
  "status": 200,
  "message": "Success get device",
  "data": {
    "id": 1,
    "name": "TV 1",
    "ip": "192.168.1.10",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

---

## 📌 Create Device

**POST** `/device`

### Request Body

```json
{
  "ip": "192.168.1.10",
  "name": "TV 1"
}
```

### Response

```json
{
  "message": "Device created successfully",
  "status": 201
}
```

---

## 📌 Update Device

**PATCH** `/device/:id`

### Request Body

```json
{
  "ip": "192.168.1.11",
  "name": "TV Updated"
}
```

### Response

```json
{
  "message": "Device updated successfully",
  "status": 200
}
```

---

## 📌 Delete Device

**DELETE** `/device/:id`

### Response

```json
{
  "message": "Device deleted successfully",
  "status": 200
}
```

---

## 📡 Ping Device

**GET** `/device/ping/:id`

### Response

```json
{
  "status": 200,
  "message": "Success ping device",
  "data": {
    "isReachable": true,
    "isAdbConnected": false
  }
}
```

---

## 📡 Ping Multiple Devices

**POST** `/device/multiple-ping`

### Request Body

```json
{
  "ids": [1, 2, 3]
}
```

### Response

```json
{
  "status": 200,
  "message": "Success ping multiple devices",
  "data": [
    {
      "id": 1,
      "ip": "192.168.1.10",
      "name": "TV 1",
      "isReachable": true,
      "isAdbConnected": false
    },
    {
      "id": 2,
      "isReachable": false,
      "isAdbConnected": false,
      "error": "Device not found"
    }
  ]
}
```

---

## 🔧 Send ADB Action to Device

**GET** `/device/action/:id?action=<ACTION>`

### Available `action` values

- `sleep`
- `wake`
- `reboot`
- `shutdown`
- (Lihat `dto.AllActions`)

### Example

```
GET /device/action/1?action=sleep
```

### Response

```json
{
  "status": 200,
  "message": "Action 'sleep' sent to device",
  "data": {
    "id": 1,
    "ip": "192.168.1.10",
    "action": "sleep"
  }
}
```

---

## 🛡️ Authentication & Middleware

Semua endpoint ini:

- **Memerlukan token JWT** di header `Authorization: Bearer <token>`
- **Beberapa action menggunakan middleware `RequireModuleAccess("device")`**
