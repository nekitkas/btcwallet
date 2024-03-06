package api

import (
	"btcwallet/internal/model"
	"errors"
	"fmt"
	"math"
	"net/http"
)

func (s *server) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactions, err := s.store.Transaction().Get()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, transactions)
	}
}

func (s *server) balance() http.HandlerFunc {
	type responseBody struct {
		EUR float64 `json:"EUR"`
		BTC float64 `json:"BTC"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		amount, err := s.store.Transaction().Balance()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		rate, err := getExchangeRate(exchangeUrl)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		eur := amount * rate
		s.respond(w, r, http.StatusOK, Response{Data: responseBody{
			EUR: eur,
			BTC: amount,
		}})
	}
}

func (s *server) transfer() http.HandlerFunc {
	type requestBody struct {
		Amount float64 `json:"amount"`
	}

	type response struct {
		Spent bool    `json:"spent"`
		BTC   float64 `json:"BTC"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		exchangeRate, err := getExchangeRate(exchangeUrl)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		var body requestBody
		if err = s.decode(r, &body); err != nil {
			s.error(w, r, http.StatusInternalServerError, fmt.Errorf("expected numeric value: %w", err))
			return
		}
		transAmountBTC := body.Amount / exchangeRate

		if err = s.store.Transaction().Transfer(transAmountBTC); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, Response{Data: response{
			Spent: true,
			BTC:   transAmountBTC,
		}})
	}
}

func (s *server) add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transaction := model.NewTransaction()
		if err := s.decode(r, transaction); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if transaction.Amount <= 0 {
			s.error(w, r, http.StatusInternalServerError, errors.New("amount should be more than 0"))
			return
		}

		exchangeRate, err := getExchangeRate(exchangeUrl)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		transaction.Amount = math.Round((transaction.Amount/exchangeRate)*1e6) / 1e6
		if err = s.store.Transaction().Add(transaction, transaction.Amount); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, Response{Data: transaction})
	}
}
