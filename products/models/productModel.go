package models

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const ProductsTableName = "products"

type ProductFields struct {
	ProductID   string `bson:"_id"`
	ProductName string `bson:"productname"`
	ProductCode string `bson:"productcode"`
}

func DBGetProducts(
	dbRef *mongo.Database,
) []ProductFields {
	log.Debug("models DBGetProducts")
	c := dbRef.Collection(ProductsTableName)
	cur, err := c.Find(context.Background(), nil)
	if err != nil {
		log.Errorf("DBGetProducts Find error: %v", err)
		return nil
	}
	log.Debugf("DBGetProducts cursor: %v", cur)

	defer cur.Close(context.Background())

	var products []ProductFields

	for cur.Next(context.Background()) {
		// elem := bson.NewDocument()
		// err := cur.Decode(elem)

		elem := bson.D{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Errorf("DBGetProducts Decode error: %v", err)
			return nil
		}

		// do something with elem....
		// log.Debugf("DBGetProducts elem: %v", elem)
		// // 	CreatedAt:   elem.Lookup("createdAt").DateTime().UTC(),
		// p := ProductFields{
		// 	// lookup for mongodb bson alias names
		// 	ProductID:   elem.Lookup("_id").ObjectID().Hex(),
		// 	ProductName: elem.Lookup("productname").StringValue(),
		// 	ProductCode: elem.Lookup("productcode").StringValue(),
		// }

		// // append struct items to slice of structs
		// products = append(products, p)

		var p ProductFields

		pMap := elem.Map()
		log.Debugf("DBGetProducts pMap: %v", pMap)

		ProductID, ok := pMap["_id"].(objectid.ObjectID)
		if ok {
			p.ProductID = ProductID.Hex()
		}
		ProductName, ok := pMap["productname"].(string)
		if ok {
			p.ProductName = ProductName
		}
		ProductCode, ok := pMap["productcode"].(string)
		if ok {
			p.ProductCode = ProductCode
		}

		// append struct items to slice of structs
		products = append(products, p)
	}

	if err := cur.Err(); err != nil {
		log.Errorf("DBGetProducts Cursor error: %v", err)
		return nil
	}

	// return slice of product structs
	return products
}

func DBCreateProduct(
	dbRef *mongo.Database,
	ProductName string,
	ProductCode string,
) {
	log.Debug("models CreateProduct")
	// p := ProductFields{
	// 	ProductID:   objectid.New(),

	// var p struct {
	// 	ProductName string
	// 	ProductCode string
	// }
	// p.ProductName = ProductName
	// p.ProductCode = ProductCode

	c := dbRef.Collection(ProductsTableName)

	// res, err := c.InsertOne(
	// 	context.Background(),
	// 	p,
	// )

	res, err := c.InsertOne(
		context.Background(),
		bson.D{
			{Key: "productname", Value: ProductName},
			{Key: "productcode", Value: ProductCode},
		},
	)

	if err != nil {
		log.Errorf("Echo Collection Insert error: %v", err)
	}

	id := res.InsertedID.(objectid.ObjectID)
	log.Errorf("models DBCreateProduct Inserted id: %v", id)
}
