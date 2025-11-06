package controller

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"be-internship/config"
	"be-internship/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
 
// Fungsi utama untuk insert koleksi (pakai form-data)
func InsertKoleksi(c *fiber.Ctx) error {
	noReg := c.FormValue("no_reg")
	noInv := c.FormValue("no_inv")
	namaBenda := c.FormValue("nama_benda")
	kategoriID := c.FormValue("kategori_id") // 🔹 ambil ID kategori, bukan nama
	bahan := c.FormValue("bahan")
	ukuran := c.FormValue("ukuran")
	tahunPerolehan := c.FormValue("tahun_perolehan")
	asalPerolehan := c.FormValue("asal_perolehan")
	ket := c.FormValue("ket")
	tempat := c.FormValue("tempat_penyimpanan")

	if noReg == "" || noInv == "" || namaBenda == "" || kategoriID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Field penting tidak boleh kosong (no_reg, no_inv, nama_benda, kategori_id).",
		})
	}

	// 🔹 Cek kategori berdasarkan ID
	objID, err := primitive.ObjectIDFromHex(kategoriID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID kategori tidak valid.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var kategori model.Kategori
	kategoriCollection := config.Ulbimongoconn.Collection("kategori")
	err = kategoriCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&kategori)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Kategori tidak ditemukan.",
		})
	}

	// 🔹 Upload gambar
	file, err := c.FormFile("foto")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File foto wajib diunggah.",
		})
	}

	imageURL, err := uploadImageToGitHub(file, namaBenda)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Gagal upload gambar ke GitHub: %v", err),
		})
	}

	// 🔹 Buat data koleksi
	data := model.Koleksi{
		ID:                primitive.NewObjectID(),
		NoRegistrasi:      noReg,
		NoInventaris:      noInv,
		NamaBenda:         namaBenda,
		Kategori:          kategori, // isi dengan hasil pencarian
		Bahan:             bahan,
		Ukuran:            ukuran,
		TahunPerolehan:    tahunPerolehan,
		AsalPerolehan:     asalPerolehan,
		Keterangan:        ket,
		TempatPenyimpanan: tempat,
		Foto:              imageURL,
		CreatedAt:         time.Now(),
	}

	collection := config.Ulbimongoconn.Collection("koleksi")
	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan ke database: " + err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":   "Koleksi berhasil disimpan.",
		"image_url": imageURL,
	})
}


// =============================================================
// 🟣 Fungsi Upload Gambar ke GitHub
// =============================================================
func uploadImageToGitHub(file *multipart.FileHeader, namaBenda string) (string, error) {
	githubToken := os.Getenv("GH_ACCESS_TOKEN") // Pastikan sudah di-set
	repoOwner := "ghaidafasya24"
	repoName := "images-koleksi-museum"
	filePath := fmt.Sprintf("koleksi/%d_%s.jpg", time.Now().Unix(), namaBenda)

	if githubToken == "" {
		return "", fmt.Errorf("GH_ACCESS_TOKEN belum diatur di environment variable")
	}

	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %w", err)
	}
	defer f.Close()

	imageData, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file: %w", err)
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	payload := map[string]string{
		"message": fmt.Sprintf("Upload image for %s", namaBenda),
		"content": encodedImage,
	}
	payloadBytes, _ := json.Marshal(payload)

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, filePath)

	req, _ := http.NewRequest("PUT", apiURL, bytes.NewReader(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal request ke GitHub API: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	content, ok := result["content"].(map[string]interface{})
	if !ok || content["download_url"] == nil {
		return "", fmt.Errorf("tidak menemukan download_url dari GitHub response")
	}

	return content["download_url"].(string), nil
}
