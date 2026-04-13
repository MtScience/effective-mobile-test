package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"subscriptions/internal/logging"

	"github.com/gorilla/mux"

	domain "subscriptions/internal/domain/subscription"
	usecase "subscriptions/internal/usecase/subscription"
)

type SubscriptionHandler struct {
	usecase usecase.Usecase
}

func NewSubscriptionHandler(usecase usecase.Usecase) *SubscriptionHandler {
	return &SubscriptionHandler{usecase}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info("received subscription creation request")

	var req subscriptionModificationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	created, err := h.usecase.Create(r.Context(), usecase.CreateInput{
		UserID:         req.UserID,
		Service:        req.Service,
		Price:          req.Price,
		SubscribedOn:   req.SubscribedOn,
		UnsubscribedOn: req.UnsubscribedOn,
	})
	if err != nil {
		logging.Logger.Error("error during subscription creation", "err", err)

		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newSubscriptionDTO(created))
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info("received subscription retrieval request")

	id, err := getIDFromRequest(r)
	if err != nil {
		logging.Logger.Error("couldn't extract subscription ID from request")

		writeError(w, http.StatusBadRequest, err)
		return
	}

	sub, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		logging.Logger.Error("error during subscription retrieval", "err", err)

		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newSubscriptionDTO(sub))
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info("received subscription update request")

	id, err := getIDFromRequest(r)
	if err != nil {
		logging.Logger.Error("couldn't extract subscription ID from request")

		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req subscriptionModificationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	updated, err := h.usecase.Update(r.Context(), id, usecase.UpdateInput{
		UserID:         req.UserID,
		Service:        req.Service,
		Price:          req.Price,
		SubscribedOn:   req.SubscribedOn,
		UnsubscribedOn: req.UnsubscribedOn,
	})
	if err != nil {
		logging.Logger.Error("error during subscription update", "err", err)

		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newSubscriptionDTO(updated))
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info("received subscription deletion request")

	id, err := getIDFromRequest(r)
	if err != nil {
		logging.Logger.Error("couldn't extract subscription ID from request")

		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		logging.Logger.Error("error during subscription deletion", "err", err)

		writeUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info("received subscription list request")

	subs, err := h.usecase.List(r.Context())
	if err != nil {
		logging.Logger.Error("error during subscription retrieval", "err", err)

		writeUsecaseError(w, err)
		return
	}

	response := make([]subscriptionDTO, 0, len(subs))
	for i := range subs {
		response = append(response, newSubscriptionDTO(&subs[i]))
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *SubscriptionHandler) TotalPrice(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info("received total price request")

	var req priceRequestDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	price, err := h.usecase.TotalPrice(r.Context(), usecase.PriceRequestInput{
		UserID:    req.UserID,
		Service:   req.Service,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		logging.Logger.Error("error during subscription price calculation", "err", err)

		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, priceDTO{price})
}

func getIDFromRequest(r *http.Request) (int64, error) {
	rawID := mux.Vars(r)["id"]
	if rawID == "" {
		return 0, errors.New("missing subscription id")
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid subscription id")
	}

	if id <= 0 {
		return 0, errors.New("invalid subscription id")
	}

	return id, nil
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		logging.Logger.Error("error decoding JSON", "err", err)
		return err
	}

	return nil
}

func writeUsecaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, usecase.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}
