# ExamDemon v1.0 — by Bibi318 [08/2026]

> ExamTopics scraper nhẹ, nhanh, không cần Docker.  
> Chỉ 1 file binary, chạy ngay trên Windows/Linux/macOS.
<img width="1869" height="1298" alt="image" src="https://github.com/user-attachments/assets/a2fb7e53-0a9f-45da-bc08-03c333030f55" />

---

## Cài đặt nhanh

**Yêu cầu:** [Go >= 1.21](https://go.dev/dl/)

```bash
git clone <repo>
cd examdemon
go build -o examdemon.exe .
```

Hoặc chạy trực tiếp không cần build:
```bash
go run . -p amazon -s scs-c03
```

---

## Sử dụng

```
examdemon.exe [options]

  -p string    Nhà cung cấp chứng chỉ (default: google)
  -s string    Chuỗi tìm kiếm tên đề thi (bắt buộc)
  -o string    Đường dẫn file đầu ra (default: output.md)
  -f string    Định dạng: md | json | txt | html (default: md)
  -c           Bao gồm bình luận cộng đồng
  -w int       Số goroutine đồng thời (default: 8)
  -exams       Liệt kê tất cả đề thi của nhà cung cấp
```

---

## Ví dụ

```bash
# Tải AWS Security SCS-C03
examdemon.exe -p amazon -s scs-c03 -f html -o scs-c03.html

# Tải Cisco 200-301 dạng JSON
examdemon.exe -p cisco -s 200-301 -f json -o ccna.json

# Tải với bình luận cộng đồng
examdemon.exe -p comptia -s security+ -c -o sec-plus.md

# Xem tất cả đề thi của Amazon
examdemon.exe -p amazon -exams

# Nhiều goroutine hơn để tải nhanh hơn
examdemon.exe -p google -s cloud-digital -w 15 -f html
```

---

## Nhà cung cấp hỗ trợ

| Flag `-p`        | Nhà cung cấp       |
|------------------|--------------------|
| `amazon`         | AWS Certifications |
| `cisco`          | Cisco / CCNA       |
| `comptia`        | CompTIA            |
| `microsoft`      | Microsoft / Azure  |
| `google`         | Google Cloud       |
| `isc2`           | ISC2 / CISSP       |
| `isaca`          | ISACA / CISM       |
| `fortinet`       | Fortinet NSE       |
| `juniper`        | Juniper            |
| `ec-council`     | CEH / EC-Council   |
| `oracle`         | Oracle             |
| `paloaltonetworks` | Palo Alto PCNSA  |

---

## Chức năng (examtopics-downloader)

| Tính năng         | ExamDemon        |
|-------------------|----------------  |
| Ngôn ngữ          | Go               |
| Docker            | **Không cần**    |
| Setup             | `go run .`       |
| Binary size       | **~9MB (-s -w)** |
| Dependencies      | **1** (goquery)  |
| PDF export        | Không (nhẹ hơn)  |
| GitHub cache      | Không (đơn giản) |
| Định dạng output  | md/json/html/txt |
| Progress display  | inline counter   |

---

## Build tối ưu (binary nhỏ nhất)

```bash
# Windows
go build -ldflags="-s -w" -o examdemon.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o examdemon .

# macOS ARM (M1/M2)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o examdemon .
```

---

*ExamDemon v1.0 — Tạo bởi Bibo318, tháng 8/2026*
