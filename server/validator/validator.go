package validator

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"net/http"
)

const (
	HeaderChallenge = "X-POS-Challenge"
	HeaderSolution  = "X-POS-Solution"
)

type Interface interface {
	InitChallenge(w http.ResponseWriter, r *http.Request) (bool, error)
	CheckChallenge(r *http.Request) error
}

type cachier interface {
	Set(key string)
	Exist(key string) bool
	Delete(key string)
}

func NewValidator(cache cachier, binZerosCount int) Interface {
	return &validator{
		cache:         cache,
		binZerosCount: binZerosCount,
	}
}

type validator struct {
	cache         cachier
	binZerosCount int
}

func (v *validator) InitChallenge(w http.ResponseWriter, r *http.Request) (bool, error) {
	if r.Method == http.MethodHead {
		key := make([]byte, 32)

		if _, err := rand.Read(key); err != nil {
			return false, fmt.Errorf("unable generate challange string")
		}
		challenge := base64.StdEncoding.EncodeToString(key)
		w.Header().Set(HeaderChallenge, challenge)
		v.cache.Set(challenge)

		return true, nil
	}
	return false, nil
}

func (v *validator) CheckChallenge(r *http.Request) error {
	challenge := r.Header.Get(HeaderChallenge)
	solution := r.Header.Get(HeaderSolution)

	if challenge == "" || !v.cache.Exist(challenge) {
		return fmt.Errorf("challenge is not exist")
	}
	defer v.cache.Delete(challenge)
	if solution == "" {
		return fmt.Errorf("challenge is empty")
	}
	if !strings.Contains(solution, challenge) {
		return fmt.Errorf("solution is not valid")
	}
	sum := sha256.Sum256([]byte(solution))
	if countFirstZeroBites(sum) < v.binZerosCount {
		return fmt.Errorf("solution is not valid")
	}
	return nil
}

func countFirstZeroBites(bytes [sha256.Size]byte) int {
	k := 0
	for i := 0; i < len(bytes); i++ {
		for j := 0; j < 8; j++ {
			if bytes[i]>>(7-j)&1 == 0 {
				k++
			} else {
				return k
			}
		}
	}
	return k
}
