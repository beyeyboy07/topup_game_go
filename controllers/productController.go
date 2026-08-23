package controllers

import (
	"net/http"

	"topup_games_go/models"
	"topup_games_go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductController menyimpan dependency
// yang dibutuhkan Product Controller.
type ProductController struct {
	DB *gorm.DB
}

// CreateProduct membuat product baru.
//
// POST /api/products
func (controller *ProductController) CreateProduct(
	c *gin.Context,
) {

	// Request body.
	var request struct {
		GameID    uint   `json:"game_id"`
		Name      string `json:"name"`
		Code      string `json:"code"`
		BuyPrice  int64  `json:"buy_price"`
		SellPrice int64  `json:"sell_price"`
	}

	// Decode JSON.
	if err := c.ShouldBindJSON(&request); err != nil {

		utils.Error(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)

		return
	}

	// Validasi.
	if request.GameID == 0 ||
		request.Name == "" ||
		request.Code == "" {

		utils.Error(
			c.Writer,
			http.StatusBadRequest,
			"Game ID, name and code are required",
		)

		return
	}

	// Pastikan game ada.
	var game models.Game

	result := controller.DB.First(
		&game,
		request.GameID,
	)

	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusNotFound,
			"Game not found",
		)

		return
	}

	// Buat Product.
	product := models.Product{
		GameID:    request.GameID,
		Name:      request.Name,
		Code:      request.Code,
		BuyPrice:  request.BuyPrice,
		SellPrice: request.SellPrice,
	}

	// Simpan ke database.
	result = controller.DB.Create(&product)

	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)

		return
	}

	// Response.
	response := models.ProductResponse{
		ID:        product.ID,
		GameID:    product.GameID,
		Name:      product.Name,
		Code:      product.Code,
		BuyPrice:  product.BuyPrice,
		SellPrice: product.SellPrice,
	}

	utils.Success(
		c.Writer,
		http.StatusCreated,
		"Product created successfully",
		response,
	)
}

// GetProducts mengambil semua product.
//
// GET /api/products
func (controlles *ProductController) GetProducts(c *gin.Context) {

	// Menampung semua product.
	var products []models.Product

	// Ambil product sekaligus data Game.
	//
	// Preload("Game") akan mengambil
	// data game berdasarkan GameID.

	result := controlles.DB.Preload("Game").Find(&products)

	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusInternalServerError,
			"Failed to get products",
		)

		return
	}

	// Menyiapkan response.
	responses := make(
		[]models.ProductResponse,
		0,
		len(products),
	)

	// Convert Product -> ProductResponse.

	for _, product := range products {

		responses = append(
			responses,
			models.ProductResponse{
				ID:        product.ID,
				GameID:    product.GameID,
				Name:      product.Name,
				Code:      product.Code,
				BuyPrice:  product.BuyPrice,
				SellPrice: product.SellPrice,
			},
		)
	}

	// Response.
	utils.Success(
		c.Writer,
		http.StatusOK,
		"Products retrieved successfully",
		responses,
	)

}

// GetProductByID mengambil satu product berdasarkan ID.
//
// GET /api/products/:id

func (controller *ProductController) GetProductByID(c *gin.Context) {

	// Ambil ID dari URL.
	//
	// Contoh:
	// /api/products/1
	//
	// hasil:
	// "1"
	id := c.Param("id")

	var product models.Product

	// Cari product dan sekaligus
	// mengambil data Game.

	result := controller.DB.Preload("Game").First(&product, id)

	// Jika product tidak ditemukan.
	if result.Error != nil {
		utils.Error(
			c.Writer,
			http.StatusNotFound,
			"Product not found",
		)
		return
	}

	// Convert Product menjadi response.
	response := models.ProductResponse{

		ID:        product.ID,
		GameID:    product.GameID,
		Name:      product.Name,
		Code:      product.Code,
		BuyPrice:  product.BuyPrice,
		SellPrice: product.SellPrice,
	}

	// Response.
	utils.Success(
		c.Writer,
		http.StatusOK,
		"Product retrieved successfully",
		response,
	)
}

func (controller *ProductController) UpdateProduct(data *gin.Context) {

	// Ambil ID dari URL.
	id := data.Param("id")

	// Cari product.
	var product models.Product

	result := controller.DB.First(&product, id)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusNotFound,
			"Product not found",
		)

		return
	}

	var request struct {
		GameID    uint   `json:"game_id"`
		Name      string `json:"name"`
		Code      string `json:"code"`
		BuyPrice  int64  `json:"buy_price"`
		SellPrice int64  `json:"sell_price"`
	}

	// Decode JSON.

	if err := data.ShouldBindJSON(&request); err != nil {
		utils.Error(
			data.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)

		return
	}

	// Validasi.
	if request.GameID == 0 ||
		request.Name == "" ||
		request.Code == "" {

		utils.Error(
			data.Writer,
			http.StatusBadRequest,
			"Game ID, name and code are required",
		)

		return
	}

	var game models.Game

	result = controller.DB.First(
		&game,
		request.GameID,
	)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusNotFound,
			"Game not found",
		)

		return
	}

	// Update product.
	product.GameID = request.GameID
	product.Name = request.Name
	product.Code = request.Code
	product.BuyPrice = request.BuyPrice
	product.SellPrice = request.SellPrice

	// Simpan perubahan.
	result = controller.DB.Save(
		&product,
	)

	if result != nil {
		utils.Error(
			data.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)
		return
	}

	// Response.

	response := models.ProductResponse{
		ID:        product.ID,
		GameID:    product.GameID,
		Name:      product.Name,
		Code:      product.Code,
		BuyPrice:  product.BuyPrice,
		SellPrice: product.SellPrice,
	}

	utils.Success(
		data.Writer,
		http.StatusOK,
		"Product updated successfully",
		response,
	)

}

// DeleteProduct menghapus product.
//
// DELETE /api/products/:id
func (controller *ProductController) DeleteProduct(data *gin.Context) {

	id := data.Param("id")

	// Cari product.
	var product models.Product

	result := controller.DB.First(&product, id)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusNotFound,
			"Product not found",
		)

		return
	}

	// Hapus product.
	//
	// Karena menggunakan gorm.Model,
	// GORM akan melakukan soft delete.

	result = controller.DB.Delete(&product)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)

		return
	}

	// Response.
	utils.Success(
		data.Writer,
		http.StatusOK,
		"Product deleted successfully",
		nil,
	)

}
