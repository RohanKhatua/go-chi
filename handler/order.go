package handler

import (
	"fmt"
	"net/http"
)

type Order struct {
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Order Create")
}

func (o *Order) List (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Order List")
}

func (o *Order) GetByID (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Order Get")
}

func (o *Order) UpdateByID (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Order Update")
}

func (o *Order) DeleteByID (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Order Delete")
}