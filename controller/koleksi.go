package controller

import (
	"be-internship/config"
	"be-internship/model"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertKoleksi(c *fiber.Ctx) error {
	// Bind data menu dari request body
	var koleksi model.Koleksi
	if err := c.BodyParser(&koleksi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	// Validasi: Periksa jika NoRegistrasi kosong
	if koleksi.NoRegistrasi == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No regist is required and cannot be empty",
		})
	}

	// Validasi: Periksa jika NoInventaris kosong atau nol
	if koleksi.NoInventaris == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No inv is required and cannot be zero",
		})
	}

	// Validasi: Periksa jika NamaBenda kosong
	if koleksi.NamaBenda == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Description is required and cannot be empty",
		})
	}

	// Validasi: Periksa jika NamaKategori kosong
	if koleksi.NamaKategori.NamaKategori == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Description is required and cannot be empty",
		})
	}

	// Validasi: Periksa jika Bahan kosong
	if koleksi.Bahan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Menu categories is required and cannot be empty",
		})
	}

	// Validasi: Periksa jika Ukuran kosong atau nol
	if koleksi.Ukuran == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ukuran is required and cannot be zero",
		})
	}

	// Validasi: Periksa jika TahunPerolehan kosong atau nol
	if koleksi.TahunPerolehan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "TahunPerolehan is required and cannot be zero",
		})
	}

	// Validasi: Periksa jika AsalPerolehan kosong
	if koleksi.AsalPerolehan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "AsalPerolehan is required and cannot be empty",
		})
	}

	// Validasi: Periksa jika Keterangan kosong
	if koleksi.Keterangan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Keterangan is required and cannot be empty",
		})
	}

	// Validasi: Periksa jika TempatPenyimpanan kosong
	if koleksi.TempatPenyimpanan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "TempatPenyimpanan is required and cannot be empty",
		})
	}

	// Proses upload gambar
	// file, err := c.FormFile("Image")
	// if err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
	// 		"status":  http.StatusBadRequest,
	// 		"message": "Image file is required: " + err.Error(),
	// 	})
	// }
	// imageURL, err := UploadImageToGitHub(file, koleksi.NamaBenda)
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  http.StatusInternalServerError,
	// 		"message": err.Error(),
	// 	})
	// }

	// koleksi.Foto = imageURL // Tambahkan ID unik dan waktu pembuatan
	koleksi.ID = primitive.NewObjectID()
	koleksi.CreatedAt = time.Now()

	// Connect ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil koleksi koleksi
	koleksiCollection := config.Ulbimongoconn.Collection("koleksi")

	// Masukkan data menu ke MongoDB
	insertedID, err := koleksiCollection.InsertOne(ctx, koleksi)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert menu",
		})
	}
	// Response sukses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      http.StatusOK,
		"message":     "Koleksi data saved successfully.",
		"inserted_id": insertedID,
		// "image_url":   imageURL,
	})
}

// func UploadImageToGitHub(file *multipart.FileHeader, namaKoleksi string) (string, error) {
// 	githubToken := os.Getenv("GH_ACCESS_TOKEN")
// 	repoOwner := "ghaidafasya24"
// 	repoName := "images-koleksi-museum"
// 	filePath := fmt.Sprintf("koleksi/%d_%s.jpg", time.Now().Unix(), namaKoleksi)

// 	fileContent, err := file.Open()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to open image file: %w", err)
// 	}
// 	defer fileContent.Close()

// 	imageData, err := ioutil.ReadAll(fileContent)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read image file: %w", err)
// 	}

// 	encodedImage := base64.StdEncoding.EncodeToString(imageData)
// 	payload := map[string]string{
// 		"message": fmt.Sprintf("Add image for product %s", namaKoleksi),
// 		"content": encodedImage,
// 	}
// 	payloadBytes, _ := json.Marshal(payload)

// 	req, _ := http.NewRequest("PUT", fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, filePath), bytes.NewReader(payloadBytes))
// 	req.Header.Set("Authorization", "Bearer "+githubToken)
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to upload image to GitHub: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated {
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		return "", fmt.Errorf("GitHub API error: %s", body)
// 	}

// 	var result map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return "", fmt.Errorf("failed to parse GitHub API response: %w", err)
// 	}

// 	content, ok := result["content"].(map[string]interface{})
// 	if !ok || content["download_url"] == nil {
// 		return "", fmt.Errorf("GitHub API response missing download_url")
// 	}

// 	return content["download_url"].(string), nil
// }
