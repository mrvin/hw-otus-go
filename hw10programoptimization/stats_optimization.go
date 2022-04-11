package hw10programoptimization

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/mrvin/hw-otus-go/hw10programoptimization/user"
	"github.com/ugorji/go/codec"
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	content, err := ioutil.ReadAll(r)
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
