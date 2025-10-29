package route

import (
	"be-internship/config"
	"be-internship/controller"
	"be-internship/model"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes initializes all the application routes
func SetupRoutes(app *fiber.App) {
	// Group API routes
	api := app.Group("/api")

	// User routes
	userRoutes := api.Group("/users")
	userRoutes.Post("/register", controller.Register) // Route untuk registrasi pengguna
	userRoutes.Post("/login", controller.Login)       // Route untuk login pengguna

	// Tambahkan kategori route
	kategoriRoutes := api.Group("/kategori")

	// ✅ [POST] Tambah kategori
	kategoriRoutes.Post("/", func(c *fiber.Ctx) error {
		var data model.Kategori
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request tidak valid"})
		}

		if data.NamaKategori == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama kategori wajib diisi"})
		}

		id, err := controller.InsertCategory(config.Ulbimongoconn, "kategori", data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menambah kategori"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Kategori berhasil ditambahkan",
			"id":      id.Hex(),
		})
	})


	// Menu routes
	// menuRoutes := api.Group("/menu")
	// // menuRoutes.Post("/", controller.JWTAuth, controller.InsertMenu)
	// menuRoutes.Post("/", controller.InsertMenu)    // Insert menu
	// menuRoutes.Get("/", controller.GetAllMenu)     // Route untuk mengambil semua menu
	// menuRoutes.Get("/:id", controller.GetMenuByID) // Route untuk mengambil menu berdasarkan ID
	// // menuRoutes.Put("/:id", controller.JWTAuth, controller.UpdateMenu)
	// menuRoutes.Put("/:id", controller.UpdateMenu) // Route untuk memperbarui menu berdasarkan ID
	// // menuRoutes.Delete("/:id", controller.JWTAuth, controller.DeleteMenu)
	// menuRoutes.Delete("/:id", controller.DeleteMenu) // Route untuk menghapus menu berdasarkan ID
}
