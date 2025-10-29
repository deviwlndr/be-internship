package controller

import (
	"be-internship/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ✅ Fungsi untuk menambahkan kategori baru
func InsertCategory(db *mongo.Database, col string, kategori model.Kategori) (primitive.ObjectID, error) {
	// Membuat dokumen BSON untuk disimpan ke MongoDB
	categoryData := bson.M{
		"nama_kategori": kategori.NamaKategori,
	}

	// Menyisipkan dokumen ke koleksi
	result, err := db.Collection(col).InsertOne(context.Background(), categoryData)
	if err != nil {
		fmt.Printf("InsertCategory error: %v\n", err)
		return primitive.NilObjectID, err
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

