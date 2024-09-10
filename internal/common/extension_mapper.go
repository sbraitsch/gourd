package common

import "github.com/rs/zerolog/log"

var extMap = make(map[string]string)

func init() {
	extMap[".go"] = "go"
	extMap["go"] = ".go"
	extMap[".java"] = "text/x-java"
	extMap["text/x-java"] = ".java"
	extMap[".kt"] = "text/x-kotlin"
	extMap["text/x-kotlin"] = ".kt"
	extMap[".rs"] = "rust"
	extMap["rust"] = ".rs"
	extMap[".js"] = "javascript"
	extMap["javascript"] = ".js"
	extMap[".ts"] = "application/typescript"
	extMap["application/typescript"] = ".ts"
	extMap[".py"] = "python"
	extMap["python"] = ".py"
}

func ResolveExtMode(input string) string {
	value, exists := extMap[input]
	if !exists {
		log.Info().Msg(value)
		return "text/plain"
	}
	return value
}
