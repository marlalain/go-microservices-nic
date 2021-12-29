package handlers

import (
	"fmt"
	"go-microservices-nic/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	} else if r.Method == http.MethodPost {
		p.postProduct(rw, r)
		return
	} else if r.Method == http.MethodPut {
		id, err := p.extractIdFromUrl(rw, r)
		if err != nil {
			return
		}
		p.l.Printf("Trying to update product #%d", id)
		p.putProduct(id, rw, r)
		return
	} else if r.Method == http.MethodDelete {
		id, err := p.extractIdFromUrl(rw, r)
		if err != nil {
			return
		}
		p.l.Printf("Trying to delete product #%d", id)
		p.deleteProduct(id, rw)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) extractIdFromUrl(rw http.ResponseWriter, r *http.Request) (int, error) {
	regex := regexp.MustCompile(`/products/([0-9]+)`)
	g := regex.FindAllStringSubmatch(r.URL.Path, -1)

	if len(g) != 1 || len(g[0]) != 2 {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return -1, fmt.Errorf("could not parse ID")
	}

	idString := g[0][1]
	id, _ := strconv.Atoi(idString)

	return id, nil
}

func (p *Products) getProducts(rw http.ResponseWriter, _ *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	p.l.Println("server: Returning list of products...")
}

func (p *Products) postProduct(rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	data.AddProduct(prod)
	p.l.Printf("server: Saving product...")
}

func (p *Products) putProduct(id int, rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Product could not found", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
	p.l.Printf("server: Updating product...")
}

func (p *Products) deleteProduct(id int, rw http.ResponseWriter) {
	err := data.DeleteProduct(id)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Product could not found", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
	p.l.Printf("server: Deleting product...")
}
