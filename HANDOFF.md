# HANDOFF: boilerplate ready-to-use

Repo: `svc-boilerplate-fiber-golang` (Go Fiber v2, layered domain).
Remote: https://github.com/IamKenra/boilerplate-golang.git
Baseline: commit `66cdcca` (sudah ada: Health() lengkap, add-domain.sh fixed + repo wiring,
default PORT, .env.example, .gitignore, CORS dari env, DB pool dari env + Ping, test postgre & message).

Belum ada yang diverifikasi dengan toolchain (sesi sebelumnya read-only).
Langkah pertama: `go build ./... && go vet ./... && go test ./...`.

## 1. Security
- `utils/jwt/jwt.go` getJWTSecret(): hapus fallback `default-secret-for-development-only`.
  Validasi `JWT_SECRET` saat startup di `main.go` (fail-fast kalau kosong).
- `utils/middleware/auth.go` JWTProtected(): ganti parse manual dengan `jwt.ValidateAccessToken`.
  Bug: baca `claims["uuid"]`, padahal token tulis `user_uuid`, jadi Locals("uuid") selalu nil.
  Set Locals dari `*jwt.Claims`: `user_uuid`, `username`, `authority`.
- `main.go`: tambah `app.Use(recover.New())` (`github.com/gofiber/fiber/v2/middleware/recover`).
- `utils/database/builder.go`: komentar di atas QueryBuilder: key map (kolom, operator, ORDER_BY/GROUP_BY)
  digabung langsung ke SQL, jadi wajib whitelist dan tidak boleh berasal dari input user.
- `git rm --cached .DS_Store utils/.DS_Store` (sudah ter-track, `.gitignore` tidak berlaku untuk file itu).

## 2. Cleanup sisa project asal dan duplikasi
- Hapus `entity/example-entity.go` (AssetMaintenance, typo tag `from:`).
- Hapus `utils/auth/validator.go` (duplikat JWTProtected). Grep caller dulu.
- Hapus `GetDatabaseDSN()` di `utils/config/config.go` (tidak dipakai).
- `domain/masterdata/repository/define.go`: hapus `POSTGRE_JENIS/STATUS/LOKASI/VENDOR/DEPARTEMEN_TABLE`,
  perbaiki typo "exaple".
- `main.go`: hapus `type metricsJobUsecase` dan import `_ "github.com/lib/pq"` (sudah di utils/database).
- `builder.go`: `MYSQL_DEBUG` jadi `DB_DEBUG`, tambahkan `DB_DEBUG=false` ke `.env.example`.
- `utils/message/message.go` isReleaseMode(): cek `APP_ENV` saja. Hapus `GO_ENV`/`GIN_MODE`
  (sisa Gin) dan sesuaikan `message_test.go` + `.env.example`.

## 3. Runtime
- `utils/database/transaction.go:38`: `tx.Commit()` diabaikan, jadi commit gagal dianggap sukses (risiko data loss).
  Ganti dengan `return tx.Commit()`.
- `domain/masterdata/repository/postgre-utils.go`: snowflake node dari env `NODE_ID` (default 1, valid 0-1023),
  karena multi-replica dengan node sama bisa menghasilkan ID duplikat. Tambahkan ke `.env.example`.
  Ubah juga template di `add-domain.sh` kalau ikut generate file ini.
- GenerateUUID: pakai `uuid.NewString()` (`google/uuid` sudah terinstall), tanpa round-trip DB.
  Ubah juga template `add-domain.sh`.
- `main.go` graceful shutdown: `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)`,
  Listen di goroutine, `<-ctx.Done()`, `app.ShutdownWithTimeout(10*time.Second)`, `postgreDB.Close()`.
- Cek `add-domain.sh` masih menemukan anchor `postgreDB := ConnectPostgre()` setelah main.go diubah.

## 4. Oracle
- `utils/database/oracle.go`: tambah `//go:build oracle` di baris pertama.
  godror butuh CGO + Oracle Instant Client dan akan merusak build `CGO_ENABLED=0` / alpine.
- README: build Oracle pakai `go build -tags oracle`. `utils/media` butuh binary `cwebp`/`dwebp` untuk WebP.

## 5. Verifikasi
- Tambah `utils/database/builder_test.go` (tanpa DB):
  - select Postgre + where AND + LIMIT, cek query `... WHERE uuid = $1 LIMIT $2 OFFSET $3` dan urutan Value.
  - update tanpa where harus return error `unsafe update with no where`.
  - insert Postgre, cek placeholder `($1, $2, $3)`.
- Jalankan:
  - `gofmt -l .` (harus kosong)
  - `go build ./... && go vet ./... && go test ./...`
  - `CGO_ENABLED=0 go build ./...` (harus lolos setelah langkah 4)
  - `go build -tags oracle ./...` (opsional, butuh CGO)
  - salin repo ke `/tmp`, jalankan `./add-domain.sh demo`, lalu `go build ./...` di salinan itu.

## Keputusan (sengaja tidak dikerjakan)
- `utils/response` dan `utils/message` keduanya dipertahankan. Satukan saat ada domain nyata.
- `utils/media` dipertahankan, cukup didokumentasikan.
- Tidak ada Dockerfile, Makefile, viper, DI framework, atau swagger. Tambahkan saat ada kebutuhan deploy nyata.
- `binding:"required"` di valueobject adalah tag Gin dan diabaikan Fiber. Tambahkan validasi input saat endpoint CRUD dibuat.
