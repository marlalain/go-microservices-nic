package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return products
}

var ErrorProductNotFound = fmt.Errorf("product not found")

func GetProduct(id int) (*Product, int, error) {
	for i, p := range products {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrorProductNotFound
}

func AddProduct(p *Product) {
	p.ID = GetNextID()
	products = append(products, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := GetProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	products[pos] = p
	return nil
}

func GetNextID() int {
	lp := products[len(products)-1]
	return lp.ID + 1
}

func DeleteProduct(id int) error {
	_, pos, err := GetProduct(id)
	if err != nil {
		return err
	}

	products = append(products[:pos], products[pos+1:]...)

	return nil
}

var products = Products{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
