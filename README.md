# Go Web Native

Aplikasi web sederhana untuk mengelola kategori dan produk menggunakan Go, `net/http`, template HTML, dan MySQL.

## Prasyarat

- Go 1.24 atau lebih baru
- MySQL yang dapat diakses dari komputer lokal, misalnya container Docker dengan port `3306` dipublikasikan

## Menjalankan MySQL

Jika container `my-mysql` sudah dibuat:

```bash
docker start my-mysql
```

Pastikan container memublikasikan port MySQL ke host, umumnya `127.0.0.1:3306`.

## Konfigurasi Database

Aplikasi memakai database `go_products`. Atur DSN di `config/database.go` agar sesuai dengan kredensial MySQL Anda. Untuk MySQL di Docker dengan port yang dipublikasikan, format DSN-nya seperti berikut:

```text
USER:PASSWORD@tcp(127.0.0.1:3306)/go_products?parseTime=true
```

Ganti `USER` dan `PASSWORD` dengan akun MySQL Anda. Pada format DSN driver MySQL ini, karakter `@` pada password dapat mengganggu pemisahan kredensial dan alamat; gunakan password tanpa `@` atau sesuaikan cara konfigurasi koneksinya.

## Migrasi Database

Buat database dan tabel berikut melalui klien MySQL. `CREATE TABLE IF NOT EXISTS` membuat langkah ini aman dijalankan kembali jika tabel sudah ada.

```sql
CREATE DATABASE IF NOT EXISTS go_products
  DEFAULT CHARACTER SET = 'utf8mb4';

USE go_products;

CREATE TABLE IF NOT EXISTS categories (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS products (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  category_id INT NOT NULL,
  stock BIGINT NOT NULL,
  description TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX idx_products_category_id (category_id),
  CONSTRAINT fk_products_category
    FOREIGN KEY (category_id) REFERENCES categories(id)
) ENGINE=InnoDB;
```

## Data Contoh (Seeder)

Tambahkan kategori terlebih dahulu. Perintah ini melewati kategori yang namanya sudah ada.

```sql
INSERT INTO categories (name, created_at, updated_at)
SELECT seed.name, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
FROM (
  SELECT 'Elektronik' AS name
  UNION ALL SELECT 'Peralatan Rumah'
  UNION ALL SELECT 'Aksesoris'
) AS seed
WHERE NOT EXISTS (
  SELECT 1 FROM categories existing WHERE existing.name = seed.name
);
```

Lalu tambahkan produk contoh. Produk hanya dimasukkan jika namanya belum ada.

```sql
INSERT INTO products (name, category_id, stock, description, created_at, updated_at)
SELECT seed.name, categories.id, seed.stock, seed.description,
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
FROM (
  SELECT 'Keyboard Mekanis' AS name, 'Elektronik' AS category_name,
         12 AS stock, 'Keyboard untuk penggunaan harian.' AS description
  UNION ALL SELECT 'Lampu Meja LED', 'Peralatan Rumah',
         8, 'Lampu meja dengan pencahayaan LED.'
  UNION ALL SELECT 'Mouse Wireless', 'Aksesoris',
         15, 'Mouse nirkabel untuk komputer.'
) AS seed
JOIN categories ON categories.name = seed.category_name
WHERE NOT EXISTS (
  SELECT 1 FROM products existing WHERE existing.name = seed.name
);
```

## Instalasi dan Menjalankan

Jalankan perintah dari direktori root project:

```bash
go mod download
go run .
```

Server berjalan di <http://localhost:8080>.

Jika Go 1.24 dipasang di direktori lokal `~/.local/go1.24.0`, tambahkan ke `PATH` pada terminal sebelum menjalankan perintah:

```bash
export PATH="$HOME/.local/go1.24.0/bin:$PATH"
go version
go mod download
go run .
```

## Halaman Aplikasi

- `/` - halaman utama
- `/categories` - daftar kategori
- `/categories/add` - tambah kategori
- `/categories/edit?id=ID` - ubah kategori
- `/categories/delete?id=ID` - hapus kategori
- `/products` - daftar produk
- `/products/add` - tambah produk
- `/products/detail?id=ID` - detail produk
- `/products/edit?id=ID` - ubah produk
- `/products/delete?id=ID` - hapus produk

Form tambah dan ubah menggunakan `GET` untuk menampilkan halaman dan `POST` untuk mengirim perubahan.