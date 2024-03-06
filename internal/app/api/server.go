package api

import (
	"btcwallet/internal/model"
	"btcwallet/internal/store"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nekitkas/router"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

const exchangeUrl = "http://api-cryptopia.adca.sh/v1/prices/ticker"

type ctxKey int8

type Response struct {
	Data interface{} `json:"data"`
}

type server struct {
	router *router.Router
	logger *log.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: router.New(),
		logger: log.Default(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID, s.logRequest)

	s.router.GET("/transactions", s.get())
	s.router.GET("/balance", s.balance())
	s.router.POST("/transfer", s.transfer())
	s.router.POST("/add", s.add())
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{w, http.StatusOK}
		s.logger.Printf("started %s %s\nremote_addr:%s  request_id:%s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			r.Context().Value(ctxKeyRequestID),
		)
		start := time.Now()
		next.ServeHTTP(rw, r)
		s.logger.Printf("completed in %s with %d %s\nremote_addr:%s  request_id:%s",
			time.Now().Sub(start),
			rw.code,
			http.StatusText(rw.code),
			r.RemoteAddr,
			r.Context().Value(ctxKeyRequestID),
		)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) decode(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

func getExchangeRate(url string) (float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error making request %w", err)
	}

	var exchangeRates model.ExchangeRates
	if err = json.NewDecoder(resp.Body).Decode(&exchangeRates); err != nil {
		return 0, fmt.Errorf("error decoding response body %w", err)
	}

	var exchangeRate string
	for _, rate := range exchangeRates.Data {
		if rate.Symbol == "BTC/EUR" {
			exchangeRate = rate.Value
			break
		}
	}

	exchangeRateFloat, err := strconv.ParseFloat(exchangeRate, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing exchange rate %w", err)
	}

	return exchangeRateFloat, nil
}
