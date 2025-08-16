# Migration API Documentation

API ini menyediakan endpoint untuk mengelola migration database secara manual.

## Base URL

```
/migration
```

## Endpoints

### 1. Run Migration

**POST** `/migration/run`

Menjalankan migration database untuk membuat semua tabel yang diperlukan.

**Response Success (200):**

```json
{
  "success": true,
  "message": "Migration database berhasil",
  "data": {
    "tables": [
      "users",
      "roles",
      "products",
      "categories",
      "cart_items",
      "devices",
      "packages",
      "carts",
      "billings",
      "members",
      "transactions",
      "transaction_details"
    ]
  }
}
```

**Response Error (500):**

```json
{
  "success": false,
  "message": "Migration database gagal",
  "error": "error message"
}
```

### 2. Get Migration Status

**GET** `/migration/status`

Mengecek status migration database dan melihat tabel mana yang sudah ada dan belum ada.

**Response Success (200):**

```json
{
  "success": true,
  "data": {
    "total_tables": 12,
    "existing_tables": ["users", "roles", "products"],
    "missing_tables": [
      "categories",
      "cart_items",
      "devices",
      "packages",
      "carts",
      "billings",
      "members",
      "transactions",
      "transaction_details"
    ],
    "is_complete": false
  }
}
```



## Penggunaan

### Menjalankan Migration Pertama Kali

```bash
curl -X POST http://localhost:3000/migration/run
```

### Mengecek Status Migration

```bash
curl -X GET http://localhost:3000/migration/status
```



## Catatan Penting

1. **Auto Migration**: API ini menggunakan GORM AutoMigrate yang akan membuat tabel berdasarkan struct entity yang sudah didefinisikan.

2. **Safe Operations**: Semua operasi migration bersifat aman dan tidak akan menghapus data yang sudah ada.

3. **Foreign Key**: Migration akan mengatur foreign key constraints secara otomatis berdasarkan relasi yang didefinisikan di entity.

4. **Logging**: Semua operasi migration akan di-log ke console untuk monitoring.

## Entity yang Dimigrasikan

- `User` - Tabel users
- `Role` - Tabel roles
- `Product` - Tabel products
- `Category` - Tabel categories
- `CartItem` - Tabel cart_items
- `Device` - Tabel devices
- `Package` - Tabel packages
- `Cart` - Tabel carts
- `Billing` - Tabel billings
- `Member` - Tabel members
- `Transaction` - Tabel transactions
- `TransactionDetail` - Tabel transaction_details
