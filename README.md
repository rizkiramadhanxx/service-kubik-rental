# Kubik Rental Backend

Aplikasi backend untuk sistem rental PlayStation yang dibangun dengan Go dan Fiber framework.

## Fitur

- Authentication & Authorization
- User Management
- Product & Category Management
- Cart & Billing System
- Device Management
- Member Management
- Package Management
- Transaction Management
- ADB Integration untuk PlayStation

## Persyaratan Sistem

- Go 1.24.4 atau lebih tinggi
- Windows 10/11 (untuk build .exe)
- SQLite database

## Cara Build

### Build untuk Windows (.exe)

Untuk membangun aplikasi menjadi file executable `backend.exe`, jalankan perintah berikut di terminal:

```bash
go build -o backend.exe .
```

### Build untuk Platform Lain

```bash
# Linux
go build -o backend .

# macOS
go build -o backend .

# Windows (tanpa ekstensi)
go build -o backend .
```

## Cara Menjalankan

### Setelah Build

```bash
# Windows
./backend.exe

# Linux/macOS
./backend
```

### Development Mode

```bash
go run main.go
```

## Konfigurasi

1. Buat file `.env` di root directory
2. Sesuaikan konfigurasi database dan environment variables

## Struktur Proyek

```
backend/
├── config/          # Konfigurasi database dan environment
├── docs/            # Dokumentasi API
├── dto/             # Data Transfer Objects
├── embed/           # Platform tools (ADB)
├── entity/          # Model database
├── feature/         # Fitur-fitur aplikasi
├── helpers/         # Helper functions
├── middleware/      # Middleware authentication
├── pkg/             # Package utilities
├── scheduler/       # Background jobs
└── main.go          # Entry point aplikasi
```

## API Endpoints

- `/health-check` - Health check endpoint
- `/auth/*` - Authentication routes
- `/user/*` - User management
- `/product/*` - Product management
- `/category/*` - Category management
- `/cart/*` - Cart management
- `/billing/*` - Billing system
- `/device/*` - Device management
- `/member/*` - Member management
- `/package/*` - Package management
- `/transaction/*` - Transaction management

## Database

Aplikasi menggunakan SQLite database dengan auto-migration. Database akan dibuat otomatis saat pertama kali menjalankan aplikasi.

## Dependencies

- Fiber v2 - Web framework
- GORM - ORM untuk database
- JWT - Authentication
- SQLite - Database
- Validator - Input validation
- Godotenv - Environment configuration

## Lisensi

Proyek ini dikembangkan untuk Kubik Rental.
