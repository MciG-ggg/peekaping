package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Test different regex patterns
	patterns := []string{
		// Original pattern
		`^(rediss?://)([^@]*@)?([^:]+)(:\d{1,5})?(/[0-9]*)?$`,
		// My first attempt
		`^(rediss?://)([^@]*@)?(\[?[^:/\]]*\]?|[^:/\]]+)(:\d{1,5})?(/[0-9]*)?$`,
		// Better pattern for IPv6
		`^(rediss?://)([^@]*@)?(\[[^\]]+\]|[^:/\]]+)(:\d{1,5})?(/[0-9]*)?$`,
		// Alternative pattern
		`^(rediss?://)([^@]*@)?(\[[^\]]+\]|[^:/]+)(:\d{1,5})?(/[0-9]*)?$`,
	}

	testCases := []string{
		"redis://localhost:6379",
		"redis://[::1]:6379",
		"redis://[::1]",
		"redis://[2001:db8::1]:6379",
		"redis://user:password@[::1]:6379",
		"redis://user:password@[::1]:6379/0",
		"rediss://[::1]:6379",
	}

	for i, pattern := range patterns {
		fmt.Printf("\nPattern %d: %s\n", i+1, pattern)
		re := regexp.MustCompile(pattern)
		for _, test := range testCases {
			matches := re.MatchString(test)
			fmt.Printf("  %s: %t\n", test, matches)
		}
	}
}
