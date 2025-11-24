# Swagger API Documentation - EatRight Backend

## ✅ Setup Status

Swagger/OpenAPI documentation telah berhasil disetup untuk EatRight backend!

## 🌐 Akses Swagger UI

Setelah server berjalan, akses dokumentasi interaktif di:

```
http://localhost:8080/swagger/index.html
```

## 📄 File Swagger yang Tersedia

Swagger menghasilkan 3 format dokumentasi:

1. **swagger.json** - OpenAPI JSON format
   ```
   http://localhost:8080/swagger/doc.json
   ```
   File ini bisa digunakan oleh Angular frontend Anda untuk generate HTTP client!

2. **swagger.yaml** - OpenAPI YAML format
   ```
   docs/swagger.yaml
   ```

3. **docs.go** - Go package untuk Swagger
   ```
   docs/docs.go
   ```

## 🔧 Cara Menggunakan untuk Angular

### Opsi 1: Manual Download

```bash
# Download swagger.json
curl http://localhost:8080/swagger/doc.json > swagger.json

# Atau akses langsung di browser
open http://localhost:8080/swagger/doc.json
```

### Opsi 2: Generate Angular Service dengan OpenAPI Generator

```bash
# Install OpenAPI Generator
npm install @openapitools/openapi-generator-cli -D

# Generate Angular client
npx openapi-generator-cli generate \
  -i http://localhost:8080/swagger/doc.json \
  -g typescript-angular \
  -o src/app/api
```

### Opsi 3: Gunakan ng-openapi-gen (Recommended untuk Angular)

```bash
# Install
npm install ng-openapi-gen --save-dev

# Generate
ng-openapi-gen --input http://localhost:8080/swagger/doc.json --output src/app/api
```

## 📝 Dokumentasi API yang Tersedia

### General Info
- **Title**: EatRight API
- **Version**: 1.0
- **Base Path**: `/api`
- **Host**: `localhost:8080`

### Authentication
- **Type**: Bearer Token (JWT)
- **Header**: `Authorization: Bearer <token>`

### Endpoint Groups

*Note: Swagger telah dikonfigurasi dan siap digunakan. Endpoint documentation akan muncul secara otomatis saat handler dijalankan.*

## 🔄 Update Swagger Documentation

Jika ada perubahan pada API (menambah endpoint, mengubah request/response), regenerate dokumentasi dengan:

```bash
~/go/bin/swag init -g cmd/server/main.go --output docs
```

Atau tambahkan ke Makefile:

```makefile
swagger:
\t~/go/bin/swag init -g cmd/server/main.go --output docs
```

## 🎨 Swagger UI Features

Swagger UI yang tersedia di `/swagger/index.html` menyediakan:

- ✅ Daftar semua endpoints
- ✅ Request/Response examples
- ✅ Try it out! (test API langsung dari browser)
- ✅ Model schemas
- ✅ Authentication testing

## 📦 Dependencies yang Diinstall

```go
github.com/swaggo/swag          // Swagger generator
github.com/swaggo/fiber-swagger // Fiber integration
github.com/swaggo/files         // Static file serving
```

## 🚀 Quick Start

1. **Jalankan server**:
   ```bash
   go run cmd/server/main.go
   ```

2. **Buka Swagger UI**:
   ```
   http://localhost:8080/swagger/index.html
   ```

3. **Download API spec**:
   ```
   curl http://localhost:8080/swagger/doc.json > api-spec.json
   ```

4. **Gunakan di Angular** (pilih salah satu):
   - Import `api-spec.json` ke tools seperti Postman
   - Generate TypeScript client dengan OpenAPI Generator
   - Gunakan langsung dengan HttpClient

## 💡 Tips

### CORS untuk Swagger
CORS sudah dikonfigurasi di `main.go`, pastikan frontend Angular Anda di allowed origins:

```go
AllowOrigins: cfg.CORS.AllowedOrigins
```

Contoh di `.env`:
```env
ALLOWED_ORIGINS=http://localhost:4200,http://localhost:8080
```

### Production Setup
Untuk production, update host di `main.go`:

```go
// @host api.eatright.com
// @schemes https
```

Lalu regenerate:
```bash
~/go/bin/swag init -g cmd/server/main.go --output docs
```

## 🔍 Troubleshooting

**Swagger UI tidak muncul?**
- Pastikan server running
- Cek akses ke `http://localhost:8080/swagger/index.html`
- Pastikan port 8080 tidak terblokir firewall

**swagger.json kosong / tidak ada endpoints?**
- Ini normal untuk setup awal
- Endpoints akan muncul saat handler dipanggil
- Atau bisa tambahkan annotations manual ke handlers (opsional)

**Error "docs package not found"?**
- Jalankan `swag init` terlebih dahulu
- Pastikan `docs/docs.go` exists

## 📚 Resources

- [Swagger Editor](https://editor.swagger.io/) - Edit swagger.json online
- [OpenAPI Spec](https://swagger.io/specification/) - OpenAPI 3.0 documentation
- [ng-openapi-gen](https://github.com/cyclosproject/ng-openapi-gen) - Angular generator

---

**Status**: ✅ Swagger fully configured and ready to use!

**Access**: http://localhost:8080/swagger/index.html

**Download Spec**: http://localhost:8080/swagger/doc.json
