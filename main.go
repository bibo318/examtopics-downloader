package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	AppName   = "ExamDemon"
	Version   = "v1.0"
	Creator   = "Demongod"
	BuildDate = "08/2026"
)

func banner() {
	fmt.Println("\033[1;31m╔══════════════════════════════════════════╗")
	fmt.Printf( "║  %-10s %-6s by %-10s [%s]  ║\n", AppName, Version, Creator, BuildDate)
	fmt.Println("║  ExamTopics Scraper — nhẹ, nhanh, không  ║")
	fmt.Println("║  cần Docker. Chỉ cần: go run .           ║")
	fmt.Println("╚══════════════════════════════════════════╝\033[0m")
	fmt.Println()
}

func main() {
	provider     := flag.String("p", "google", "Nhà cung cấp chứng chỉ (amazon, cisco, comptia...)")
	grepStr      := flag.String("s", "", "Chuỗi tìm kiếm tên đề thi (vd: scs-c03, devops)")
	outputPath   := flag.String("o", "output.md", "Đường dẫn file đầu ra")
	format       := flag.String("f", "md", "Định dạng: md | json | txt | html")
	withComments := flag.Bool("c", false, "Bao gồm bình luận cộng đồng")
	listExams    := flag.Bool("exams", false, "Liệt kê tất cả đề thi của nhà cung cấp rồi thoát")
	workers      := flag.Int("w", 8, "Số goroutine đồng thời (mặc định 8)")
	flag.Parse()

	banner()

	s := NewScraper(*workers)

	if *listExams {
		fmt.Printf("[*] Lấy danh sách đề thi của '%s'...\n\n", *provider)
		exams := s.GetExams(*provider)
		if len(exams) == 0 {
			fmt.Println("[!] Không tìm thấy đề thi nào.")
			os.Exit(1)
		}
		for _, e := range exams {
			fmt.Println(e)
		}
		fmt.Printf("\nTổng: %d đề thi\n", len(exams))
		return
	}

	if *grepStr == "" {
		fmt.Fprintln(os.Stderr, "[LỖI] Thiếu chuỗi tìm kiếm. Dùng -s <tên_đề_thi>")
		fmt.Fprintln(os.Stderr, "Ví dụ: go run . -p amazon -s scs-c03")
		fmt.Fprintln(os.Stderr, "       go run . -p cisco  -s 200-301 -f json")
		os.Exit(1)
	}

	fmt.Printf("  Nhà cung cấp : %s\n", *provider)
	fmt.Printf("  Tìm kiếm     : %s\n", *grepStr)
	fmt.Printf("  Định dạng    : %s\n", *format)
	fmt.Printf("  Workers      : %d\n\n", *workers)

	questions := s.ScrapeAll(*provider, *grepStr)

	if len(questions) == 0 {
		fmt.Println("\n[!] Không tìm thấy câu hỏi nào. Kiểm tra lại -p và -s.")
		os.Exit(1)
	}

	fmt.Printf("\n[+] Thu thập được %d câu hỏi.\n", len(questions))

	base := strings.TrimSuffix(*outputPath, ".md")
	var (
		outFile string
		err     error
	)

	switch *format {
	case "json":
		outFile = base + ".json"
		err = writeJSON(outFile, questions)
	case "html":
		outFile = base + ".html"
		err = writeHTML(outFile, questions, *withComments)
	case "txt":
		outFile = base + ".txt"
		err = writeTXT(outFile, questions, *withComments)
	default:
		outFile = base + ".md"
		err = writeMD(outFile, questions, *withComments)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "[LỖI] Ghi file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[✓] Đã lưu: %s\n", outFile)
}
