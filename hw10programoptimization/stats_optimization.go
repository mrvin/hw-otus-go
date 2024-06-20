package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

const sizeScannerBuffer = 256
const itemsInitializDomainStat = 512

type UserEmail struct {
	Email string
}

func getDomain(email, firstLevelDomain string) string {
	indexA := strings.LastIndex(email, "@")
	if indexA == -1 {
		return ""
	}

	domain := strings.ToLower(email[indexA+1:])

	lenDomain := len(domain)
	lenFirstLevelDomain := len(firstLevelDomain)

	if lenFirstLevelDomain+1 > lenDomain {
		return ""
	}
	for i := 1; i <= lenFirstLevelDomain; i++ {
		if domain[lenDomain-i] != firstLevelDomain[lenFirstLevelDomain-i] {
			return ""
		}
		if domain[lenDomain-(lenFirstLevelDomain+1)] != byte('.') {
			return ""
		}
	}

	return domain
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	input := bufio.NewScanner(r)
	buf := make([]byte, sizeScannerBuffer)
	input.Buffer(buf, sizeScannerBuffer)

	var user UserEmail
	var err error
	var allDomain string
	result := make(DomainStat, itemsInitializDomainStat)
	for input.Scan() {
		if err = easyjson.Unmarshal(input.Bytes(), &user); err != nil {
			return result, fmt.Errorf("unmarshal line: %w", err)
		}
		allDomain = getDomain(user.Email, domain)
		if allDomain != "" {
			result[allDomain]++
		}
	}
	if err = input.Err(); err != nil {
		return result, fmt.Errorf("read line: %w", err)
	}

	return result, nil
}
