package util

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var _ = flag.String("envfile", "", "The environment-file to load")

type EnvFile struct {
	sync.Mutex
	initialized bool
	FileName    string
	Content     map[string]string
}

var env_file *EnvFile = &EnvFile{
	FileName: "",
	initialized: false,
	Content: map[string]string{},
}

func Env_Filename() string {
	return env_file.FileName
}

func (f *EnvFile) init() {
	f.Lock()
	defer f.Unlock()
	if !f.initialized {
		var init_verbose = false
		for i,arg := range os.Args {
			if arg == "-envfile" {
				if len(os.Args) > i+1 {
					f.FileName=os.Args[i+1]
				}
			}
		}
		f.initialized = true

		if f.FileName == "" {
			if init_verbose {
				log.Printf("No envfile used")
			}
			return
		}
		if init_verbose {
			log.Printf("Env Filename is [%s]",f.FileName)
		}
		
		var err error
		f.Content, err = godotenv.Read(f.FileName)
		if err != nil {
			log.Fatalf("Failed to read %s gave err:%v", f.FileName, err)
		}
		// log.Printf("Just read env-file %s", f.FileName)
		f.initialized = true
	}
}

func (f *EnvFile) Get(key string) (string, bool) {
	f.init()
	val, found := f.Content[key]
	// log.Printf("Read from env-file [%s]=[%s] %t", key,val,found)
	return val, found
}

func GetEnvFile() (string, bool) {
	return env_file.FileName, env_file.initialized
}

// Get env var or default
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		// log.Printf("Read env [%s]=[%s]", key,value)
		return value
	}
	if value, ok := env_file.Get(key); ok {
		// log.Printf("Read env-file [%s]=[%s]", key,value)
		return value
	}
	// log.Printf("Read fallback [%s]=[%s]", key,fallback)
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
