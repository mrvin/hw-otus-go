package hw10programoptimization

import (
	"bytes"
	"io"

	"github.com/mrvin/hw-otus-go/hw10programoptimization/user"
	"github.com/ugorji/go/codec"
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(content, []byte{'\n'})

	jh := new(codec.JsonHandle)
	dotDomain := []byte("." + domain)
	var user user.User
	resultDomainStat := make(DomainStat)
	for _, line := range lines {
		err = codec.NewDecoderBytes(line, jh).Decode(&user)
		if err == nil {
			byteEmail := []byte(user.Email)
			if bytes.HasSuffix(byteEmail, dotDomain) {
				resultDomainStat[string(bytes.ToLower(bytes.SplitN(byteEmail, []byte{'@'}, 2)[1]))]++
			}
		}
	}

	return resultDomainStat, err
}

var numGoroutin = 4

type result struct {
	res string
	err error
}

func GetDomainStatLongLivedGorout(r io.Reader, domain string) (DomainStat, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(content, []byte{'\n'})
	numLines := len(lines)

	if numLines < numGoroutin {
		numGoroutin = numLines
	}
	perGorout := numLines / numGoroutin
	remGorout := numLines % numGoroutin
	resultCh := make(chan *result, numGoroutin)

	jh := new(codec.JsonHandle)
	dotDomain := []byte("." + domain)
	from, to := 0, 0
	for i := 0; i < numGoroutin; i++ {
		to += perGorout
		if i < remGorout {
			to++
		}
		go func(partLines [][]byte) {
			var user user.User
			for _, line := range partLines {
				err = codec.NewDecoderBytes(line, jh).Decode(&user)
				if err == nil {
					byteEmail := []byte(user.Email)
					if bytes.HasSuffix(byteEmail, dotDomain) {
						resultCh <- &result{string(bytes.ToLower(bytes.SplitN(byteEmail, []byte{'@'}, 2)[1])), nil}
						continue
					}
				}
				resultCh <- &result{"", err}
			}
		}(lines[from:to])
		from = to
	}

	resultDomainStat := make(DomainStat)
	for i := 0; i < numLines; i++ {
		res := <-resultCh
		resultDomainStat[res.res]++
		if res.err != nil {
			err = res.err
		}
	}
	delete(resultDomainStat, "")
	close(resultCh)

	return resultDomainStat, err
}
