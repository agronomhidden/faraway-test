package main

import (
	"fmt"
	"github.com/Netflix/go-env"
	"log"
	"log/slog"
	"net/http"
	"os"
	"server/cache"
	"server/quoteService"
	"server/validator"
)

type conf struct {
	Port             int `env:"FARAWAY_SERVER_PORT,default=8080"`
	HashcashBinZeros int `env:"FARAWAY_SERVER_BIN_ZEROS,default=24"`
}

var configuration conf

func init() {
	_, _ = env.UnmarshalFromEnviron(&configuration)
}

func main() {
	c := cache.NewCache()
	v := validator.NewValidator(c, configuration.HashcashBinZeros)
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", configuration.Port),
		Handler: NewHandler(v),
	}
	log.Fatal(s.ListenAndServe())
}

func NewHandler(v validator.Interface) *Handler {
	return &Handler{
		service:   quoteService.NewQuoteService(),
		log:       slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		validator: v,
	}
}

type Handler struct {
	service   quoteService.QuoteHandler
	log       *slog.Logger
	validator validator.Interface
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	init, err := h.validator.InitChallenge(w, r)
	if err != nil {
		h.log.Error("unable to init challenge", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if init {
		return
	}
	if err := h.validator.CheckChallenge(r); err != nil {
		h.log.Error("access denied", slog.Any("error", err))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	result, err := h.service.GetRandom()
	if err != nil {
		h.log.Error("unable to call quote service", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(result)); err != nil {
		h.log.Error("unable to send response", slog.Any("error", err))
	}
}
