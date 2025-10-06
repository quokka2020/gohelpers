package util

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var env_filename = flag.String("env-file", "", "The environment-file to load")

type EnvFile struct {
	initialized bool
	FileName    string
	Content     map[string]string
}

var env_file *EnvFile = &EnvFile{
	FileName: *env_filename,
	initialized: false,
	Content: map[string]string{},
}

func Env_Filename() string {
	return *env_filename
}

func (f *EnvFile) init() {
	if !f.initialized {
		if f.FileName == "" {
			f.initialized = true
			return
		}
		var err error
		f.Content, err = godotenv.Read(f.FileName)
		if err != nil {
			log.Fatalf("Failed to read %s gave err:%v", f.FileName, err)
		}
		log.Printf("Just read env-file %s", f.FileName)
		f.initialized = true
	}
}

func (f *EnvFile) Get(key string) (string, bool) {
	f.init()
	val, found := f.Content[key]
	return val, found
}

func GetEnvFile() (string, bool) {
	return env_file.FileName, env_file.initialized
}

// Get env var or default
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if value, ok := env_file.Get(key); ok {
		return value
	}
	return fallback
}

// Get env var or default
func GetEnvInt(key string, fallback int) int {
	value := GetEnv(key, "")
	if value != "" {
		i, err := strconv.Atoi(value)
		if err == nil {
			return i
		}
	}
	return fallback
}

// Get env var or default
func GetEnvBool(key string, fallback bool) bool {
	value := GetEnv(key, "")
	if value != "" {
		return strings.ToLower(value) == "true"
	}
	return fallback
}
