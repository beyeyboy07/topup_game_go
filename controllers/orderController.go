package controllers

import (
	"fmt"
	"net/http"
	"time"
	"topup_games_go/models"
	"topup_games_go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderController menyimpan dependency
// yang dibutuhkan Order Controller.

type OrderController struct {
	DB *gorm.DB
}

type OrderRequest struct {
	ProductID uint   `json:"product_id" example:"1"`
	PlayerID  string `json:"player_id" example:"123456789"`
	ServerID  string `json:"server_id" example:"1234"`
}

// CreateOrder membuat order baru.
//
// POST /api/orders
// @Summary Membuat order
// @Description Membuat order top-up baru berdasarkan product dan ID player.
// @Tags orders
// @Accept json
// @Produce json
// @Param request body OrderRequest true "Data order"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/orders [post]
func (controller *OrderController) CreateOrder(data *gin.Context) {

	// ========================================
	// 1. Request body
	// ========================================

	var request OrderRequest

	// Decode JSON.
	err := data.ShouldBindJSON(&request)

	if err != nil {

		utils.Error(
			data.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)

	}

	// ========================================
	// 2. Validasi
	// ========================================

	if request.ProductID == 0 ||
		request.PlayerID == "" {

		utils.Error(
			data.Writer,
			http.StatusBadRequest,
			"Product ID and player ID are required",
		)

		return
	}

	// ========================================
	// 3. Cari Product
	// ========================================

	var product models.Product

	result := controller.DB.Find(
		&product,
		request.ProductID,
	)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusNotFound,
			"Product not found",
		)
		return
	}

	// ========================================
	// 4. Buat Order Number
	// ========================================
	//
	// Contoh:
	// ORD-20260823-193012

	orderNumber := fmt.Sprintf(
		"ORD-%s",
		time.Now().Format("20060102150405"),
	)

	// ========================================
	// 5. Tentukan waktu expired
	// ========================================
	//
	// Order akan expired setelah 30 menit.

	expiredAt := time.Now().Add(
		30 * time.Minute,
	)

	// ========================================
	// 6. Buat Order
	// ========================================
	order := models.Order{
		UserID:      1,
		ProductID:   product.ID,
		PlayerID:    request.PlayerID,
		ServerID:    request.ServerID,
		OrderNumber: orderNumber,
		ProductName: product.Name,
		BuyPrice:    product.BuyPrice,
		SellPrice:   product.SellPrice,
		Status:      "PENDING",
		ExpiredAt:   &expiredAt,
	}

	// ========================================
	// 7. Simpan ke database
	// ========================================
	result = controller.DB.Create(&order)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)
		return
	}

	// ========================================
	// 8. Response
	// ========================================

	response := map[string]interface{}{
		"id":           order.ID,
		"order_number": order.OrderNumber,
		"product_id":   order.ProductID,
		"product_name": order.ProductName,
		"player_id":    order.PlayerID,
		"server_id":    order.ServerID,
		"buy_price":    order.BuyPrice,
		"sell_price":   order.SellPrice,
		"status":       order.Status,
		"expired_at":   order.ExpiredAt,
	}

	utils.Success(
		data.Writer,
		http.StatusOK,
		"Order created successfully",
		response,
	)

}
