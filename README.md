# ExamDemon v1.0 — by Bibi318 [08/2026]
> Lightweight, fast ExamTopics scraper. No Docker required.  
> Single binary, runs instantly on Windows/Linux/macOS.

![ExamDemon Screenshot](https://github.com/user-attachments/assets/a2fb7e53-0a9f-45da-bc08-03c333030f55)

---

## Quick Install

**Requirements:** [Go >= 1.21](https://go.dev/dl/)

```bash
git clone <repo>
cd examdemon
go build -o examdemon.exe .
```

Or run directly without building:

```bash
go run . -p amazon -s scs-c03
```

---

## Usage
---
examdemon.exe [options]
  -p string    Nhà cung cấp chứng chỉ (default: google)
  -s string    Chuỗi tìm kiếm tên đề thi (bắt buộc)
  -o string    Đường dẫn file đầu ra (default: output.md)
  -f string    Định dạng: md | json | txt | html (default: md)
  -c           Bao gồm bình luận cộng đồng
  -w int       Số goroutine đồng thời (default: 8)
  -exams       Liệt kê tất cả đề thi của nhà cung cấp

## Examples

```bash
# Download AWS Security SCS-C03
examdemon.exe -p amazon -s scs-c03 -f html -o scs-c03.html

# Download Cisco 200-301 as JSON
examdemon.exe -p cisco -s 200-301 -f json -o ccna.json

# Download with community comments
examdemon.exe -p comptia -s security+ -c -o sec-plus.md

# List all Amazon exams
examdemon.exe -p amazon -exams

# Use more goroutines for faster downloads
examdemon.exe -p google -s cloud-digital -w 15 -f html
```
---
## Supported Providers

| Flag`-p`           | Provider           |
| -------------------- | ------------------ |
| `amazon`           | AWS Certifications |
| `cisco`            | Cisco / CCNA       |
| `comptia`          | CompTIA            |
| `microsoft`        | Microsoft / Azure  |
| `google`           | Google Cloud       |
| `isc2`             | ISC2 / CISSP       |
| `isaca`            | ISACA / CISM       |
| `fortinet`         | Fortinet NSE       |
| `juniper`          | Juniper            |
| `ec-council`       | CEH / EC-Council   |
| `oracle`           | Oracle             |
| `paloaltonetworks` | Palo Alto PCNSA    |

---

## Features

| Feature          | ExamDemon              |
| ---------------- | ---------------------- |
| Language         | Go                     |
| Docker           | **Not required** |
| Setup            | `go run .`           |
| Binary size      | **~9MB (-s -w)** |
| Dependencies     | **1** (goquery)  |
| PDF export       | No (keeps it lean)     |
| GitHub cache     | No (keeps it simple)   |
| Output formats   | md / json / html / txt |
| Progress display | Inline counter         |

---

## Optimized Build (Smallest Binary)

```bash
# Windows
go build -ldflags="-s -w" -o examdemon.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o examdemon .

# macOS ARM (M1/M2)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o examdemon .
```

---

*ExamDemon v1.0 — Created by Bibi318, August 2026*
