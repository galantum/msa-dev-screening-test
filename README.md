## 📄 Cara Menjalankan
Untuk menjalankan program ini:
1. Buka terminal.
2. Jalankan perintah berikut:
  ```bash
  go run main.go
  ```

### 🔧 Konfigurasi
Untuk mengubah nilai **limit concurrency** dan direktori output, buka kode di **line 111 - 112**:
```go
concurrentLimit := 2
outputDir := "/home/yourname/csv"
```
- `concurrentLimit`: Mengatur batas jumlah concurrency yang berjalan secara bersamaan. Ubah nilai sesuai kebutuhan.
- `outputDir`: Menentukan path folder tempat hasil file CSV akan disimpan. Ubah nilai menjadi path yang diinginkan.
