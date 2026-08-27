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

type ProductRequest struct {
	GameID    uint   `json:"game_id" example:"1"`
	Name      string `json:"name" example:"86 Diamonds"`
	Code      string `json:"code" example:"ML86"`
	BuyPrice  int64  `json:"buy_price" example:"10000"`
	SellPrice int64  `json:"sell_price" example:"12000"`
}

// CreateProduct membuat product baru.
//
// POST /api/products
// @Summary Membuat product
// @Description Membuat product baru untuk game yang tersedia.
// @Tags products
// @Accept json
// @Produce json
// @Param request body ProductRequest true "Data product"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/products [post]
func (controller *ProductController) CreateProduct(
	c *gin.Context,
) {

	// Request body.
	var request ProductRequest

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
// @Summary Mengambil semua product
// @Description Mengambil daftar seluruh product.
// @Tags products
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/products [get]
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
// @Summary Mengambil product berdasarkan ID
// @Description Mengambil satu product berdasarkan ID.
// @Tags products
// @Produce json
// @Param id path int true "ID product"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/products/{id} [get]

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

// @Summary Mengubah product
// @Description Mengubah data product berdasarkan ID.
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "ID product"
// @Param request body ProductRequest true "Data product"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/products/{id} [put]
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

	var request ProductRequest

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
// @Summary Menghapus product
// @Description Menghapus product secara soft delete berdasarkan ID.
// @Tags products
// @Produce json
// @Param id path int true "ID product"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/products/{id} [delete]
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
