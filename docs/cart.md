# Dokumentasi API Cart Module

## Struktur Entitas

### Cart

- ID: uint
- Status: string
- CreatedAt: time.Time
- UpdatedAt: time.Time
- CartItems: []CartItem

### CartItem

- ID: uint
- CartID: uint (relasi ke Cart)
- ItemType: string ("product" atau "billing")
- ProductID: \*uint (nullable)
- BillingID: \*uint (nullable)
- Product: \*Product (relasi jika item_type == "product")
- Billing: \*Billing (relasi jika item_type == "billing")
- Duration: \*int
- Price: int
- TotalPrice: int
- Qty: int
- CreatedAt: time.Time
- UpdatedAt: time.Time

### Billing

- ID: uint
- DeviceID: uint
- Device: Device
- StartTime: time.Time
- EndTime: time.Time
- PackageID: uint
- Package: Package

### Package

- ID: uint
- Name: string
- Duration: int
- Price: int
- IsLoss: bool

---

## Endpoints

### [POST] /cart

Membuat cart kosong.

- Body: `entity.Cart`
- Response: `201 Created`

### [GET] /cart

Ambil semua cart dengan pagination dan optional status filter.

- Query Params: `page`, `limit`, `status`
- Response: `200 OK` dengan pagination meta

### [GET] /cart/:id

Ambil detail cart lengkap dengan cart items dan preload relasi produk/billing.

- Response: `200 OK` atau `404 Not Found`

### [PATCH] /cart/:id

Update status cart.

- Body: `{ "status": "string" }`
- Response: `200 OK` atau `404 Not Found`

### [DELETE] /cart/:id

Hapus cart.

- Response: `200 OK`

### [POST] /cart/item

Tambah item ke cart.

- Body: `CreateCartItemRequest`
  - CartID: uint
  - ItemType: string ("product" atau "billing")
  - ProductID/BillingID: uint (salah satu tergantung item_type)
  - Qty: int
- Response: `201 Created` atau `400 Bad Request` jika duplikat

### [PATCH] /cart/item/qty

Update kuantitas item dalam cart.

- Body: `UpdateQtyRequest`
  - CartID: uint
  - ItemType: string
  - ProductID/BillingID: uint
  - Action: string ("increment", "decrement", "set")
  - Value: \*int (untuk action "set")
- Response: `200 OK`, `201 Created`, atau `404 Not Found`

---

## Catatan

- Semua relasi dipastikan dengan `Preload(...)`
- Format response menggunakan `dto.Response`
- Penanganan error konsisten: 400 (bad request), 404 (not found), 500 (internal server error)
