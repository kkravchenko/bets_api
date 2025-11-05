// controller/bet_controller.go
package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
	"tz.api/internal/entity"
	"tz.api/internal/errors"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"tz.api/internal/store"
)

type BetController struct {
	Store *store.BetStore
}

func NewBetController(store *store.BetStore) *BetController {
	return &BetController{Store: store}
}

type BetRequest struct {
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	CrashPoint float64 `json:"crash_point"`
}

func (br *BetRequest) Validate() error {
	if br.UserID == "" {
		return httpError("user_id is required", http.StatusBadRequest)
	}
	if br.Amount <= 0 || br.Amount > 10000 {
		return httpError("amount must be > 0 and <= 10000", http.StatusBadRequest)
	}
	if br.CrashPoint < 1.0 || br.CrashPoint > 100.0 {
		return httpError("crash_point must be between 1.0 and 100.0", http.StatusBadRequest)
	}
	return nil
}

func httpError(msg string, code int) error {
	return fmt.Errorf(msg, code)
}

func (c *BetController) CreateBetHandler(w http.ResponseWriter, r *http.Request) {
	var req BetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		if e, ok := err.(*errors.ValidationError); ok {
			http.Error(w, e.Msg, e.Code)
		} else {
			http.Error(w, "Validation failed", http.StatusBadRequest)
		}
		return
	}

	bet := entity.Bet{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		Amount:     req.Amount,
		CrashPoint: req.CrashPoint,
		CreatedAt:  time.Now(),
	}

	c.Store.Create(bet)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bet)
}

func (c *BetController) GetBetsHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	sortBy := r.URL.Query().Get("sort")
	userID := r.URL.Query().Get("user_id")

	bets := c.Store.GetAll()

	// Фильтр по user_id
	if userID != "" {
		filtered := []entity.Bet{}
		for _, b := range bets {
			if b.UserID == userID {
				filtered = append(filtered, b)
			}
		}
		bets = filtered
	}

	// Сортировка
	sort.Slice(bets, func(i, j int) bool {
		switch sortBy {
		case "amount_desc":
			return bets[i].Amount > bets[j].Amount
		case "amount_asc":
			return bets[i].Amount < bets[j].Amount
		case "date_desc":
			return bets[i].CreatedAt.After(bets[j].CreatedAt)
		default: // date_asc
			return bets[i].CreatedAt.Before(bets[j].CreatedAt)
		}
	})

	// Пагинация
	total := len(bets)
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	resp := map[string]interface{}{
		"data":        bets[start:end],
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": (total + limit - 1) / limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *BetController) GetBetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	bet, found := c.Store.GetByID(id)
	if !found {
		http.Error(w, "Bet not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(bet)
}

func (c *BetController) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
