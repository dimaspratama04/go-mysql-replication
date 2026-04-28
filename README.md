# MySQL Replication RnD — Products API

Stack: **GoFiber** · **MySQL 8** · **GTID Replication** · **Docker Compose**

---

## Quick Start

```bash
# 1. Clone / copy project
cd mysql-replication-rnd

# 2. Generate go.sum (required once)
go mod tidy

# 3. Start semua services
make up

# 4. Cek status
make status
make replication-status
```

Tunggu ~40 detik sampai semua container healthy. Replication setup jalan otomatis.

---

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│                     Docker Network: rnd-network          │
│                                                          │
│   ┌─────────────┐    Write     ┌──────────────────────┐  │
│   │             │◄────────────►│  MySQL Primary       │  │
│   │ GoFiber App │              │  port: 3306          │  │
│   │ port: 8080  │              │  Binary Log + GTID   │  │
│   └─────────────┘              └──────────┬───────────┘  │
│                                           │              │
│                                    Replication           │
│                                    (GTID auto)           │
│                                           │              │
│                                ┌──────────▼───────────┐  │
│                                │  MySQL Replica       │  │
│                                │  port: 3307          │  │
│                                │  read_only=ON        │  │
│                                └──────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

---

## API Reference

```
GET    /health                  Health check
GET    /api/products            List all products
GET    /api/products/:id        Get product by ID
POST   /api/products            Create product
PUT    /api/products/:id        Update product
DELETE /api/products/:id        Soft delete product
```

### Example Payloads

**POST /api/products**

```json
{
  "name": "SSD Samsung 1TB",
  "description": "NVMe M.2 PCIe Gen4",
  "price": 1500000,
  "stock": 50
}
```

**PUT /api/products/:id**

```json
{
  "name": "SSD Samsung 1TB Pro",
  "description": "NVMe M.2 PCIe Gen4 - Updated",
  "price": 1350000,
  "stock": 45
}
```

---

## Log Output Examples

**CREATE:**

```
10:01:23 INF [CREATE] Product created successfully operation=CREATE product_id=6 name=SSD Samsung 1TB price=1500000 stock=50
```

**UPDATE:**

```
10:02:11 INF [UPDATE] Product updated successfully operation=UPDATE product_id=6 before.name=SSD Samsung 1TB before.price=1500000 after.price=1350000
```

**DELETE:**

```
10:03:05 INF [DELETE] Product soft-deleted successfully operation=DELETE type=soft_delete product_id=6 name=SSD Samsung 1TB deleted_at=2026-01-01T10:03:05Z
```

---

## Replication Commands

```bash
# Cek status replication
make replication-status

# Bandingkan data primary vs replica
make check-replica

# Masuk MySQL shell
make shell-primary
make shell-replica
```

**Verify replication manual:**

```sql
-- Di Primary: buat data baru
INSERT INTO products (name, description, price, stock) VALUES ('Test Repl', 'Test', 100000, 1);

-- Di Replica: cek data sudah sync (read-only)
SELECT * FROM products ORDER BY id DESC LIMIT 1;
```

---

## Makefile Commands

| Command                   | Description                          |
| ------------------------- | ------------------------------------ |
| `make up`                 | Build + start semua services         |
| `make down`               | Stop semua services                  |
| `make clean`              | Hapus containers + volumes + images  |
| `make status`             | Container status + API health        |
| `make replication-status` | MySQL replication thread status      |
| `make check-replica`      | Compare row count Primary vs Replica |
| `make test-create`        | Quick test: create product           |
| `make test-list`          | Quick test: list products            |
| `make logs-app`           | Tail app logs                        |
