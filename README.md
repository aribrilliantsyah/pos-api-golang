# Proyek Golang: API Point Of Sales (POS) Simple

Proyek ini adalah API berbasis Golang untuk mengelola penjulan sederhana. Menggunakan PostgreSQL sebagai database dan menyertakan fitur-fitur seperti pendaftaran pengguna, login, dan autentikasi JWT.

## Persyaratan Sistem

- Go (Golang) versi 1.20 atau lebih tinggi
- PostgreSQL versi 14 atau lebih tinggi

Pastikan Anda telah menginstal dan mengkonfigurasi kedua software tersebut sebelum menjalankan proyek ini.

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

1. Pastikan Anda telah menginstal Golang di sistem Anda.
2. Navigasikan ke direktori root proyek.
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

Ganti `8080` dengan nomor port yang sebenarnya jika Anda telah mengubahnya dalam konfigurasi.

## Menggunakan API

1. **Mendaftarkan Pengguna**:
   Gunakan endpoint register di Swagger UI untuk membuat akun pengguna baru.

2. **Login**:
   Setelah mendaftar, gunakan endpoint login untuk mendapatkan token JWT.

3. **Mengakses Rute Terproteksi**:
   Untuk mengakses rute terproteksi:
   - Untuk setiap endpoint API yang memerlukan autentikasi, cari kolom input header "Authorization".
   - Masukkan token JWT Anda dalam format: `Bearer <token_anda_di_sini>`
   - Kirim permintaan dengan token yang disertakan dalam header.

   Metode ini memungkinkan Anda untuk menambahkan token secara individual untuk setiap panggilan API, memberikan Anda lebih banyak kontrol atas permintaan mana yang menyertakan autentikasi.

## Catatan Keamanan

Ingatlah untuk menjaga `JWT_SECRET` Anda tetap aman dan jangan pernah meng-commit-nya ke version control. Dalam lingkungan produksi, pertimbangkan untuk menggunakan variabel lingkungan atau sistem manajemen rahasia yang aman.

