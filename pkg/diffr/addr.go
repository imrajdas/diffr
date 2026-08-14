package diffr

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func listenAndDisplay(address string, port int) (listenAddr, displayURL string, err error) {
	scheme := "http"
	host := "localhost"
	raw := strings.TrimSpace(address)
	if raw == "" {
		raw = "http://localhost"
	}

	switch {
	case strings.Contains(raw, "://"):
		u, perr := url.Parse(raw)
		if perr != nil {
			return "", "", fmt.Errorf("address: %w", perr)
		}
		if u.Scheme != "" {
			scheme = u.Scheme
		}
		if h := u.Hostname(); h != "" {
			host = h
		} else if u.Host != "" {
			host = u.Host
		}
	case strings.Contains(raw, "/") && !strings.Contains(raw, ":"):
		return "", "", fmt.Errorf("address: invalid host %q", raw)
	default:
		if h, p, splitErr := net.SplitHostPort(raw); splitErr == nil {
			host = h
			if p != "" {
				n, convErr := strconv.Atoi(p)
				if convErr != nil {
					return "", "", fmt.Errorf("address port: %w", convErr)
				}
				port = n
			}
		} else {
			host = raw
		}
	}

	if host == "" {
		host = "localhost"
	}
	listenAddr = net.JoinHostPort(host, strconv.Itoa(port))
	displayURL = scheme + "://" + listenAddr
	return listenAddr, displayURL, nil
}
