package aries

import (
	"strings"
)

var mobileUserAgents = []string{
	"Android",
	"iPhone",
	"iPad",
	"iPod",
	"IEMobile",
	"WPDesktop",
	"BlackBerry",
}

func isMobile(userAgent string) bool {
	for _, s := range mobileUserAgents {
		if strings.Contains(userAgent, s) {
			return true
		}
	}
	return false
}
