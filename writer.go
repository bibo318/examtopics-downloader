package main

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"
)

func writeMD(path string, questions []Question, withComments bool) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# %s %s — by %s [%s]\n\n", AppName, Version, Creator, BuildDate)
	fmt.Fprintf(f, "> Tổng số câu hỏi: **%d**\n\n---\n\n", len(questions))

	for i, q := range questions {
		fmt.Fprintf(f, "## Câu %d — %s\n\n", i+1, q.Title)
		if q.Topic != "" {
			fmt.Fprintf(f, "%s\n\n", q.Topic)
		}
		if q.Content != "" {
			fmt.Fprintf(f, "%s\n\n", q.Content)
		}
		for _, c := range q.Choices {
			fmt.Fprintf(f, "- %s\n", c)
		}
		fmt.Fprintln(f)
		fmt.Fprintf(f, "**Đáp án: %s**\n\n", q.Answer)
		if q.Time != "" {
			fmt.Fprintf(f, "*%s*\n\n", q.Time)
		}
		fmt.Fprintf(f, "[Xem trên ExamTopics](%s)\n\n", q.URL)
		if withComments && q.Comments != "" {
			fmt.Fprintf(f, "<details><summary>Bình luận cộng đồng</summary>\n\n%s\n\n</details>\n\n", q.Comments)
		}
		fmt.Fprintln(f, "---\n")
	}
	return nil
}

func writeJSON(path string, questions []Question) error {
	b, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func writeTXT(path string, questions []Question, withComments bool) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sep := strings.Repeat("=", 60)
	fmt.Fprintf(f, "%s %s by %s [%s]\n", AppName, Version, Creator, BuildDate)
	fmt.Fprintf(f, "Tổng số câu hỏi: %d\n%s\n\n", len(questions), sep)

	for i, q := range questions {
		fmt.Fprintf(f, "Câu %d: %s\n\n", i+1, q.Title)
		if q.Topic != "" {
			fmt.Fprintf(f, "%s\n\n", q.Topic)
		}
		if q.Content != "" {
			fmt.Fprintf(f, "%s\n\n", q.Content)
		}
		for _, c := range q.Choices {
			fmt.Fprintf(f, "  %s\n", c)
		}
		fmt.Fprintf(f, "\nĐÁP ÁN: %s\n", q.Answer)
		if q.Time != "" {
			fmt.Fprintf(f, "Thời gian: %s\n", q.Time)
		}
		fmt.Fprintf(f, "Link: %s\n", q.URL)
		if withComments && q.Comments != "" {
			fmt.Fprintf(f, "\n--- Bình luận ---\n%s\n", q.Comments)
		}
		fmt.Fprintf(f, "\n%s\n\n", strings.Repeat("-", 60))
	}
	return nil
}

func writeHTML(path string, questions []Question, withComments bool) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprint(f, `<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ExamDemon — Câu hỏi luyện thi</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:system-ui,-apple-system,sans-serif;background:#0d1117;color:#e6edf3;padding:24px 16px}
.container{max-width:860px;margin:0 auto}
header{text-align:center;padding:32px 0 24px;border-bottom:1px solid #30363d;margin-bottom:32px}
header h1{font-size:2em;color:#f78166;letter-spacing:1px}
header p{color:#8b949e;margin-top:8px}
.question{background:#161b22;border:1px solid #30363d;border-radius:10px;padding:24px;margin-bottom:24px}
.question h2{color:#79c0ff;font-size:1.1em;margin-bottom:14px;padding-bottom:10px;border-bottom:1px solid #21262d}
.topic{color:#8b949e;font-size:.88em;margin-bottom:12px;font-style:italic}
.content{margin-bottom:14px;line-height:1.6}
.choices{list-style:none;margin-bottom:16px}
.choices li{background:#21262d;border-radius:5px;padding:8px 12px;margin-bottom:6px;font-size:.95em;line-height:1.5}
.answer{display:inline-block;background:#1a4731;color:#56d364;padding:6px 16px;border-radius:20px;font-weight:700;margin-bottom:10px;font-size:.95em}
.meta{color:#6e7681;font-size:.82em;margin-top:8px}
.meta a{color:#58a6ff;text-decoration:none}
.meta a:hover{text-decoration:underline}
details{margin-top:12px}
details summary{cursor:pointer;color:#8b949e;font-size:.88em;user-select:none}
details p{font-size:.84em;color:#6e7681;margin-top:8px;line-height:1.6;white-space:pre-wrap}
.badge{background:#21262d;border:1px solid #30363d;border-radius:20px;padding:2px 10px;font-size:.75em;color:#8b949e;vertical-align:middle}
.stats{color:#8b949e;font-size:.9em;text-align:center;margin-top:32px;padding-top:16px;border-top:1px solid #30363d}
</style>
</head>
<body>
<div class="container">
<header>
`)
	fmt.Fprintf(f, "  <h1>%s <span class=\"badge\">%s</span></h1>\n", AppName, Version)
	fmt.Fprintf(f, "  <p>Tác giả: <strong>%s</strong> &nbsp;|&nbsp; Phát hành: %s</p>\n", Creator, BuildDate)
	fmt.Fprintf(f, "  <p style=\"margin-top:6px\">Tổng số câu hỏi: <strong>%d</strong></p>\n", len(questions))
	fmt.Fprint(f, "</header>\n\n")

	for i, q := range questions {
		fmt.Fprintf(f, "<div class=\"question\">\n")
		fmt.Fprintf(f, "  <h2>%d. %s</h2>\n", i+1, html.EscapeString(q.Title))
		if q.Topic != "" {
			fmt.Fprintf(f, "  <p class=\"topic\">%s</p>\n", html.EscapeString(q.Topic))
		}
		if q.Content != "" {
			fmt.Fprintf(f, "  <p class=\"content\">%s</p>\n", html.EscapeString(q.Content))
		}
		if len(q.Choices) > 0 {
			fmt.Fprint(f, "  <ul class=\"choices\">\n")
			for _, c := range q.Choices {
				fmt.Fprintf(f, "    <li>%s</li>\n", html.EscapeString(c))
			}
			fmt.Fprint(f, "  </ul>\n")
		}
		fmt.Fprintf(f, "  <div class=\"answer\">Đáp án: %s</div>\n", html.EscapeString(q.Answer))
		fmt.Fprint(f, "  <div class=\"meta\">\n")
		if q.Time != "" {
			fmt.Fprintf(f, "    <span>%s</span> &nbsp;&middot;&nbsp;\n", html.EscapeString(q.Time))
		}
		fmt.Fprintf(f, "    <a href=\"%s\" target=\"_blank\" rel=\"noopener\">Xem trên ExamTopics ↗</a>\n", html.EscapeString(q.URL))
		fmt.Fprint(f, "  </div>\n")
		if withComments && q.Comments != "" {
			fmt.Fprintf(f, "  <details><summary>Bình luận cộng đồng</summary><p>%s</p></details>\n", html.EscapeString(q.Comments))
		}
		fmt.Fprint(f, "</div>\n\n")
	}

	fmt.Fprintf(f, "<p class=\"stats\">Tạo bởi %s %s — %s [%s]</p>\n", AppName, Version, Creator, BuildDate)
	fmt.Fprint(f, "</div>\n</body>\n</html>\n")
	return nil
}
