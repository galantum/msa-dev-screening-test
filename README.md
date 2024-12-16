## 📄 Cara Menjalankan Aplikasi Chat

Aplikasi ini adalah layanan chat sederhana menggunakan **Go** dan **gRPC**. Ikuti langkah-langkah berikut untuk menjalankan server, klien, dan unit test.

---

### 🛠️ Prasyarat

- **Go** versi terbaru. [Download di sini](https://go.dev/dl/)
- Terminal CLI (Command Line Interface) seperti **Bash**, **Zsh**, atau sejenisnya.

---

### 🚀 Menjalankan Server

1. Buka terminal.
2. Jalankan perintah berikut:
   ```bash
   go run main.go
   ```

---

### 🚀 Menjalankan Klien

1. Buka terminal baru.
2. Jalankan perintah berikut:
   ```bash
   go run main_client.go
   ```

**Catatan:** Anda dapat membuka beberapa terminal untuk menjalankan beberapa klien sekaligus.

---

### 📝 Mendapatkan Riwayat Chat

Setelah klien berhasil terhubung ke server dan muncul pesan:
```
Enter message:
```

Ketikkan perintah berikut untuk melihat riwayat chat:
```
/history
```

Riwayat chat sebelumnya akan ditampilkan di terminal klien.

---

### 🧪 Menjalankan Unit Test

Untuk menjalankan pengujian unit:
1. Buka terminal.
2. Ketikkan perintah berikut:
   ```bash
   go test ./grpc-chat-service/server -v
   ```
Hasil pengujian akan ditampilkan di terminal.

---

### ✨ Tips Tambahan

- Pastikan semua dependensi telah terinstal dengan benar.
- Jika terjadi error, periksa file dan direktori proyek.
