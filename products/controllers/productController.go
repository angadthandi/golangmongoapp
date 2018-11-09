package controllers

import (
	"github.com/angadthandi/golangmongoapp/products/models"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

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
	ProductName string,
	ProductCode string,
) {
	models.DBCreateProduct(
		dbRef,
		ProductName,
		ProductCode,
	)
}
