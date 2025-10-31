package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Kategori struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NamaKategori string             `bson:"nama_kategori" json:"nama_kategori"`
	Deskripsi    string             `bson:"deskripsi,omitempty" json:"deskripsi,omitempty"`
}
