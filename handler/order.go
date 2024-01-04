package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/RohanKhatua/go-chi/model"
	"github.com/RohanKhatua/go-chi/repository/order"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Order struct {
	Repo *order.RedisRepo
}

func (h *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	order := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	err := h.Repo.Insert(r.Context(), order)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("error marshalling order")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64

	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := o.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: cursor,
		Size:   size,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}

	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error marshalling order")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func (h *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	o, err := h.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(o); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existingOrder, err := o.Repo.FindByID(r.Context(), orderID)

	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"

	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:
		if existingOrder.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		existingOrder.ShippedAt = &now

	case completedStatus:
		if existingOrder.CompletedAt != nil || existingOrder.ShippedAt == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		existingOrder.CompletedAt = &now

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(existingOrder); err != nil {
		fmt.Println("error marshalling order")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderId, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.DeleteByID(r.Context(), orderId)

	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
