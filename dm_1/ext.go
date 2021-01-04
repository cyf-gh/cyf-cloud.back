package dm_1

var (
	DMExts = map[string][]string {
		"image": []string{
			".png", ".jpeg", ".jpg", ".gif", ".bmp", ".psd", ".ai",
			".eps", ".tif", ".tiff", ".webp", ".raw", ".svg", ".svgz",
		},
		"video": []string{
			".flv", ".mp4", ".mov", ".wmv", ".avi", ".flv", ".mkv", ".swf",
		},
		"music": []string{
			".mp3", ".flac", ".wav", ".ape",
		},
		"backup": []string {
			".rar", ".zip", ".tar", ".bk", ".backup",
		},
		"document": []string {
			".doc", ".docx", ".pdf", ".ppt", ".pptx", ".xls", ".xlsx",
			".xla", ".md", ".txt", ".text",
		},
	}
)