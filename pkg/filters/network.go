package filters

import "net"

// Ip 验证过滤IP
func Ip(val string) string {
	if address := net.ParseIP(val); address != nil {
		return address.String()
	}
	return ""
}
