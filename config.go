package main

var (
	ignoredFolders = map[string]bool{
		".git": true,
	}
	allowedExtensions = map[string]bool{
		".asp":        true,
		".aspx":       true,
		".c":          true,
		".c++":        true,
		".cgi":        true,
		".cpp":        true,
		".css":        true,
		".csv":        true,
		".go":         true,
		".h":          true,
		".hpp":        true,
		".hs":         true,
		".html":       true,
		".hxx":        true,
		".java":       true,
		".js":         true,
		".jsp":        true,
		".jspx":       true,
		".markdown":   true,
		".md":         true,
		".nim":        true,
		".php":        true,
		".php4":       true,
		".php5":       true,
		".pl":         true,
		".properties": true,
		".py":         true,
		".rb":         true,
		".rhtml":      true,
		".rss":        true,
		".shtml":      true,
		".svg":        true,
		".txt":        true,
		".xml":        true,
	}
)
