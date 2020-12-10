package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println(ValidPasswords(os.Stdin))
}

func ValidPasswords(r io.Reader) (byCount, byPos int, err error) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		var (
			a, b     int
			char     rune
			password string
		)
		if _, err := fmt.Sscanf(s.Text(), "%d-%d %c: %s", &a, &b, &char, &password); err != nil {
			return 0, 0, err
		}
		found := 0
		for _, r := range password {
			if r == char {
				found++
			}
		}
		if a <= found && found <= b {
			byCount++
		}
		rpassword := []rune(password)
		if (rpassword[a-1] == char) != (rpassword[b-1] == char) {
			byPos++
		}
	}
	return byCount, byPos, s.Err()
}
