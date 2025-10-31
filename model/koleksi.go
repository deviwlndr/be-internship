package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struct untuk data koleksi museum
type Koleksi struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NoRegistrasi      string             `json:"no_reg,omitempty" bson:"no_reg,omitempty"`
	NoInventaris      string             `json:"no_inv,omitempty" bson:"no_inv,omitempty"`
	NamaBenda         string             `json:"nama_benda,omitempty" bson:"nama_benda,omitempty"`
	NamaKategori      Kategori           `json:"kategori,omitempty" bson:"kategori,omitempty"` // misal: Etnografi, Arkeologi, Numismatika, Biologi
	Bahan             string             `json:"bahan,omitempty" bson:"bahan,omitempty"`       // bahan utama koleksi (kayu, logam, kain, dll)
	Ukuran            string             `json:"ukuran,omitempty" bson:"ukuran,omitempty"`     // ukuran fisik (cm, meter, dll)
	TahunPerolehan    string             `json:"tahun_perolehan,omitempty" bson:"tahun_perolehan,omitempty"`
	AsalPerolehan     string             `json:"asal_perolehan,omitempty" bson:"asal_perolehan,omitempty"`         // asal usul atau sumber koleksi
	Keterangan        string             `json:"ket,omitempty" bson:"ket,omitempty"`                               // baik, rusak ringan, rusak berat, dsb
	TempatPenyimpanan string             `json:"tempat_penyimpanan,omitempty" bson:"tempat_penyimpanan,omitempty"` // lokasi penyimpanan (misal: Gudang A, Lemari 3)
	// Foto              string             `json:"foto,omitempty" bson:"foto,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"` // Field baru untuk waktu pembuatan menu                            // URL atau path gambar koleksi
	// Status            string `json:"status,omitempty" bson:"status,omitempty"`                         // aktif, dipinjam, diperbaiki, dipamerkan
}
