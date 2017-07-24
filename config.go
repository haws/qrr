package main

var (
	ignoredFolders = map[string]bool{
		".git": true,
	}
	allowedExtensions = map[string]bool{
		".asm":        true,
		".asp":        true,
		".aspx":       true,
		".awk":        true,
		".bash":       true,
		".bat":        true,
		".c":          true,
		".c++":        true,
		".cgi":        true,
		".class":      true,
		".cpp":        true,
		".cs":         true,
		".css":        true,
		".csv":        true,
		".docx":       true,
		".dtl":        true,
		".erb":        true,
		".factor":     true,
		".go":         true,
		".h":          true,
		".hpp":        true,
		".hs":         true,
		".hss":        true,
		".htm":        true,
		".html":       true,
		".hxx":        true,
		".ini":        true,
		".java":       true,
		".ko":         true,
		".js":         true,
		".jsp":        true,
		".jspx":       true,
		".less":       true,
		".markdown":   true,
		".md":         true,
		".mk":         true,
		".ml":         true,
		".nim":        true,
		".odt":        true,
		".php":        true,
		".php4":       true,
		".php5":       true,
		".pl":         true,
		".pptx":       true,
		".properties": true,
		".py":         true,
		".rb":         true,
		".rc":         true,
		".rhtml":      true,
		".rss":        true,
		".rtf":        true,
		".sass":       true,
		".scss":       true,
		".sh":         true,
		".shtml":      true,
		".sql":        true,
		".svg":        true,
		".swift":      true,
		".tex":        true,
		".txt":        true,
		".vb":         true,
		".xhtml":      true,
		".xlsx":       true,
		".xml":        true,
		".y":          true,
		".yaws":       true,
		".yxx":        true,
	}
)
