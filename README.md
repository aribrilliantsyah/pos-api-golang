# Proyek Golang: API Point Of Sales (POS)

Proyek ini adalah API berbasis Golang untuk mengelola sistem point of sale (POS) dengan fitur lengkap. Menggunakan PostgreSQL sebagai database dan menyertakan berbagai modul pengelolaan bisnis.

## Persyaratan Sistem

- Go (Golang) versi 1.20 atau lebih tinggi
- PostgreSQL versi 14 atau lebih tinggi

Pastikan Anda telah menginstal dan mengkonfigurasi kedua software tersebut sebelum menjalankan proyek ini.

## Fitur Utama

### 🔐 Autentikasi
- Sistem keamanan menggunakan JSON Web Token (JWT)
- Kontrol akses berbasis peran (role-based access)

### 📦 Modul-Modul Inti

#### Manajemen Kategori
- Operasi CRUD untuk kategori produk
- Kategorisasi produk yang efisien

#### Manajemen Produk
- Operasi CRUD lengkap untuk produk
- Pengelolaan detail dan harga produk

#### Manajemen Pengguna
- Peran pengguna: Admin dan Kasir
- Operasi CRUD untuk akun pengguna
- Izin berbasis peran

#### Manajemen Pelanggan
- Operasi CRUD untuk data pelanggan
- Pengelolaan profil pelanggan

#### Manajemen Inventori
- Pelacakan stok produk
- Pengelolaan inventori masuk dan keluar
- Riwayat stok

#### Pemrosesan Transaksi
- Sistem manajemen pesanan
- Kemampuan pemrosesan pengembalian
- Riwayat transaksi

#### Modul Pelaporan
- Laporan penjualan komprehensif
- Visualisasi dan analisis data
- Fungsi ekspor

## Konfigurasi

Sebelum menjalankan proyek, Anda perlu menyiapkan file konfigurasi. Buat file bernama `config.json` di direktori `util/config` dengan konten berikut:

```json
{
  "DB_DRIVER": "postgres",
  "DB_SOURCE": "postgresql://postgres:postgres@localhost:5333/db_book?sslmode=disable",
  "POSTGRES_USER": "postgres",
  "POSTGRES_PASSWORD": "postgres",
  "POSTGRES_DB": "contact_db",
  "SERVER_ADDRESS": "8080",
  "JWT_SECRET": "nM4t0fw80-qY3jd1N1CRPbRfrB6JiX-D-UZl6uMMmb8"
}
```

Pastikan untuk menyesuaikan nilai-nilai ini sesuai dengan pengaturan spesifik Anda, terutama detail koneksi database dan JWT secret.

## Menjalankan Proyek

Untuk memulai proyek:

1. Pastikan Anda telah menginstal Golang di sistem Anda
2. Navigasikan ke direktori root proyek
3. Jalankan perintah berikut:

   ```
   go run main.go
   ```

   Ini akan memulai server dan melakukan auto-migrasi skema database.

## Dokumentasi API

Dokumentasi API tersedia melalui Swagger UI. Setelah server berjalan, Anda dapat mengakses dokumentasi Swagger di:

```
http://localhost:8080/docs/index.html
```

```
http://localhost:8080/scalar/docs
```

Ganti `8080` dengan nomor port yang sebenarnya jika Anda telah mengubahnya dalam konfigurasi.

## Tech Stack

- Backend: Golang
- Database: PostgreSQL
- Autentikasi: JWT (JSON Web Token)
- SQL Tools: SQLC
- Documentation: OAS, Swagger, Scalar