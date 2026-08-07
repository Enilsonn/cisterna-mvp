package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func loadEnv() {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		candidate := filepath.Join(wd, ".env")
		if _, err := os.Stat(candidate); err == nil {
			file, err := os.Open(candidate)
			if err == nil {
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())
					if line == "" || strings.HasPrefix(line, "#") {
						continue
					}
					parts := strings.SplitN(line, "=", 2)
					if len(parts) != 2 {
						continue
					}
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
						value = strings.Trim(value, "\"")
					}
					if _, exists := os.LookupEnv(key); !exists {
						_ = os.Setenv(key, value)
					}
				}
				_ = file.Close()
			}
			break
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
}
