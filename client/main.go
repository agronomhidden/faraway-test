package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/Netflix/go-env"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	HeaderChallenge = "X-POS-Challenge"
	HeaderSolution  = "X-POS-Solution"
)

type conf struct {
	Url              string `env:"FARAWAY_CLIENT_URL,default=http://localhost:8080"`
	HashcashBinZeros int    `env:"FARAWAY_CLIENT_BIN_ZEROS,default=24"`
}

var configuration conf

func init() {
	_, _ = env.UnmarshalFromEnviron(&configuration)
}

func main() {
	for {
		call()
		time.Sleep(2 * time.Second)
	}
}

func call() {
	r, err := http.DefaultClient.Head(configuration.Url)
	if err != nil {
		log.Println(err)
		return
	}
	challenge := r.Header.Get(HeaderChallenge)
	if challenge == "" {
		log.Println("challenge is empty")
		return
	}
	t := time.Now()
	solution := calculate(challenge)
	delay := time.Now().Sub(t).Seconds()
	log.Printf("solution was found in %0.4fs", delay)

	req, err := http.NewRequest(http.MethodGet, configuration.Url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set(HeaderChallenge, challenge)
	req.Header.Set(HeaderSolution, solution)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	quote, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("RESULT:", string(quote))
}

func calculate(input string) string {
	i := 1
	for {
		sol := fmt.Sprintf("%s%x", input, i)

		sum := sha256.Sum256([]byte(sol))
		bitesCount := countFirstZeroBites(sum)
		if bitesCount >= configuration.HashcashBinZeros {
			return sol
		}
		i++
	}
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
