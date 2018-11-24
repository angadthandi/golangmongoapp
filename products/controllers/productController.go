package controllers

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/products/models"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

type NewProduct struct {
	ProductName string `json:"productname"`
	ProductCode string `json:"productcode"`
}

type RecieveProduct struct {
	Product NewProduct `json:"product"`
}

func GetProducts(
	dbRef *mongo.Database,
) []models.ProductFields {
	log.Debug("controllers GetProducts")

	products := models.DBGetProducts(dbRef)
	log.Debugf("controllers products: %v", products)

	return products
}

func CreateProduct(
	dbRef *mongo.Database,
	jsonMsg json.RawMessage,
) []models.ProductFields {
	var msg RecieveProduct
	err := json.Unmarshal(jsonMsg, &msg)
	if err != nil {
		log.Errorf("CreateProduct Unable to unmarshal json: %v",
			err)
		return nil
	}

	models.DBCreateProduct(
		dbRef,
		msg.Product.ProductName,
		msg.Product.ProductCode,
	)

	return GetProducts(dbRef)
}
