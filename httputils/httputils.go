package httputils

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	geoip2 "github.com/oschwald/geoip2-golang"
)

var (
	ErrNotConnected = errors.New("Your device is not connected to the internet")
)

// ProxyClient return a http client with "http(s)?_proxy" support
func ProxyClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func LocalIPAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", ErrNotConnected
}

func GeoIP(addr string, geoDB []byte) (country string, latitude float64, longitude float64, err error) {
	db, err := geoip2.FromBytes(geoDB)
	if err != nil {
		return country, latitude, longitude, err
	}
	defer db.Close()

	ip := net.ParseIP(addr)
	record, err := db.City(ip)
	if err != nil {
		return country, latitude, longitude, err
	}
	return record.Country.Names["en"], record.Location.Latitude, record.Location.Longitude, nil
}
