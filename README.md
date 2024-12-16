
## 📄 Cara Menjalankan Aplikasi Microservice

Proyek ini adalah aplikasi API **microservice** sederhana yang mencakup **Service User**, **Service Order**, dan **Service Gateway** menggunakan **Go**. Ikuti panduan berikut untuk menjalankan layanan-layanan tersebut serta melakukan pengujian unit.

---

### 🛠️ Prasyarat

Sebelum memulai, pastikan Anda telah menyiapkan:

1. **Go** versi terbaru  
   👉 [Download di sini](https://go.dev/dl/)
2. Terminal CLI, seperti **Bash**, **Zsh**, atau yang serupa.
3. **Postman** untuk menguji endpoint API.

---

### ⚙️ Konfigurasi

1. Salin file ``.env.example`` menjadi ``.env``.

2. Lengkapi parameter ``SECRET_KEY`` di dalam file ``.env`` dengan nilai yang diinginkan.

3. Perbarui nilai ``secretKey`` pada file ``user/main.go`` di baris 37 agar sesuai dengan nilai ``SECRET_KEY`` yang telah diatur di file ``.env``.

---

### 🚀 Cara Menjalankan Service

### 1. Service User

#### Langkah Menjalankan:

1. Buka terminal.
2. Jalankan perintah berikut:
   ```bash
   go run user/main.go
   ```

#### Langkah Pengujian:

Gunakan **Postman** untuk mengakses endpoint berikut:

- **Login**  
  **POST** `/user/login`  
  Isi *body* (raw, JSON) seperti contoh berikut:

  ```json
  {
      "email": "admin@example.com",
      "password": "admin"
  }
  ```

  **Respon**: Token akan dikembalikan.

Gunakan **Postman** dengan **token** yang diperoleh dari login. 

Header: Tambahkan token ke **Authorization** dengan tipe **Bearer Token**.  

- **Mengakses Profil Pengguna**  
  **GET** `/user/profile/{username}`  
  Gunakan salah satu `username` berikut: `admin`, `staff`, `manager`.

- **Memperbarui Profil Pengguna**  
  **PUT** `/user/profile/{username}`  
  Isi *body* (raw, JSON) seperti berikut:

  ```json
  {
      "username": "admin",
      "email": "admin@gmail.com",
      "age": 17
  }
  ```

---

### 2. Service Order

#### Langkah Menjalankan:

1. Buka terminal baru.
2. Jalankan perintah berikut:
   ```bash
   go run order/main.go
   ```

#### Langkah Pengujian:

Gunakan **Postman** dengan token yang diperoleh dari login. Tambahkan ke **Authorization** dengan tipe **Bearer Token** untuk mengakses endpoint berikut:

- **List Semua Pesanan**  
  **GET** `/order`

- **Detail Pesanan Spesifik**  
  **GET** `/order/{id}`

- **Membuat Pesanan Baru**  
  **POST** `/order`  
  Isi *body* (raw, JSON) seperti berikut:

  ```json
  {
      "items": [
          { "id": 1, "name": "Laptop", "price": 5000000 },
          { "id": 4, "name": "SSD 500G", "price": 1000000 }
      ]
  }
  ```

- **Memperbarui Pesanan Spesifik**  
  **PUT** `/order/{id}`  
  Isi *body* (raw, JSON) seperti berikut:

  ```json
  {
      "items": [
          { "id": 1, "name": "Laptop", "price": 5000000 },
          { "id": 4, "name": "SSD 500G", "price": 1000000 }
      ]
  }
  ```

- **Menghapus Pesanan Spesifik**  
  **DELETE** `/order/{id}`

---

### 3. Service Gateway

#### Langkah Menjalankan:

1. Buka terminal baru.
2. Jalankan perintah berikut:
   ```bash
   go run gateway/main.go
   ```

---

### 🧪 Menjalankan Unit Test

Untuk memastikan setiap bagian aplikasi berjalan dengan baik, jalankan unit test dengan langkah berikut:

1. Buka terminal.
2. Jalankan perintah berikut:
   ```bash
   go test ./microservice/gateway -v
   ```
   Hasil pengujian akan ditampilkan langsung di terminal.

---

### ✨ Tips dan Catatan Tambahan

- Pastikan semua dependensi sudah terinstal dengan baik. Gunakan perintah berikut untuk memastikan:
  ```bash
  go mod tidy
  ```
- Jika terjadi error, periksa kembali struktur direktori proyek dan kode sumber.
- Gunakan **Postman** atau alat pengujian API lainnya untuk eksplorasi lebih lanjut.

Selamat mencoba dan semoga berhasil! 🚀
