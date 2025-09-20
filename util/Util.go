package util

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func IsEmpty(s string) bool {
	return len(s) == 0
}

func NotEmpty(s string) bool {
	return len(s) > 0
}

// Get env var or default
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Get env var or default
func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err == nil {
			return i
		}
	}
	return fallback
}

// Get env var or default
func GetEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return strings.ToLower(value) == "true"
	}
	return fallback
}

func Contains(s []string, str string) bool {
	for /*nr*/ _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func FileExists(filename string) bool {
	if IsEmpty(filename) {
		return false
	}
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

type HostPort struct {
	Host string
	Port int
}

func ToHostPort(in string) (HostPort, error) {
	host, portStr, err := net.SplitHostPort(in)
	if err == nil {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			if strings.Contains(host,":") {
				return HostPort{Host: fmt.Sprintf("[%s]",host), Port: port}, nil
			}

			return HostPort{Host: host, Port: port}, nil
		}
		return HostPort{}, errors.New("failed to parse port")
	}
	return HostPort{}, err
}

func UriToHostPort(in string) (HostPort, error) {
	url, urlerr := url.ParseRequestURI(in)
	if urlerr!=nil {
		return HostPort{}, urlerr
	}
	if strings.Index(url.Host, ":") > 0 {
		return ToHostPort(url.Host)
	}
	switch url.Scheme {
	case "http": return HostPort{
		Host: url.Host,
		Port: 80,
	},nil
	case "ws": return HostPort{
		Host: url.Host,
		Port: 80,
	},nil
	default:
		return HostPort{}, errors.New("failed to parse port")
	}
}

func (hp HostPort) String() string {
	return fmt.Sprintf("%s:%d", hp.Host, hp.Port)
}

func MaxInt(l int, r int) int {
	if l < r {
		return r
	}
	return l
}

func SingleHeader(r *http.Request, name string) (string, bool) {
	val := r.Header.Get(name)
	return val, NotEmpty(val)
}

func RemoteAddr(r *http.Request) string {
	if val, found := SingleHeader(r, "X-Forwarded-For"); found {
		return val
	}
	return r.RemoteAddr
}

func MinUint(l, r uint) uint {
	if l > r {
		return r
	}
	return l
}

func MaxUint(l, r uint) uint {
	if l < r {
		return r
	}
	return l
}