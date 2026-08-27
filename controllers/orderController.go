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

type OrderStatusRequest struct {
	Status string `json:"status" example:"PROCESSING"`
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
// @Security BearerAuth
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
		UserID:      data.MustGet("user_id").(uint),
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

// GetOrderByID mengambil detail order berdasarkan ID.
//
// GET /api/orders/:id
// @Summary Mengambil order berdasarkan ID
// @Description Mengambil satu order berdasarkan ID.
// @Tags orders
// @Produce json
// @Param id path int true "ID order"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/orders/{id} [get]
func (controller *OrderController) GetOrderByID(data *gin.Context) {

	// Ambil ID dari URL.
	id := data.Param("id")

	// Siapkan object Order.
	var order models.Order

	// Cari order berdasarkan ID.
	query := controller.DB
	if data.GetString("user_role") != "admin" {
		query = query.Where("user_id = ?", data.MustGet("user_id").(uint))
	}
	result := query.First(&order, id)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusNotFound,
			"Order not found",
		)

		return
	}

	// Response.

	response := map[string]interface{}{
		"id":                      order.ID,
		"order_number":            order.OrderNumber,
		"product_id":              order.ProductID,
		"product_name":            order.ProductName,
		"player_id":               order.PlayerID,
		"server_id":               order.ServerID,
		"buy_price":               order.BuyPrice,
		"sell_price":              order.SellPrice,
		"status":                  order.Status,
		"provider_transaction_id": order.ProviderTransactionID,
		"expired_at":              order.ExpiredAt,
		"paid_at":                 order.PaidAt,
		"created_at":              order.CreatedAt,
	}

	utils.Success(
		data.Writer,
		http.StatusOK,
		"Order retrieved successfully",
		response,
	)

}

// GetOrders mengambil seluruh order.
//
// @Summary Mengambil semua order
// @Description Mengambil daftar seluruh order.
// @Tags orders
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/orders [get]
func (controller *OrderController) GetOrders(data *gin.Context) {

	// Menampung semua order.
	var orders []models.Order

	// Ambil data dari database.
	query := controller.DB
	if data.GetString("user_role") != "admin" {
		query = query.Where("user_id = ?", data.MustGet("user_id").(uint))
	}
	result := query.Find(&orders)

	if result.Error != nil {

		utils.Error(
			data.Writer,
			http.StatusInternalServerError,
			"Failed to get orders",
		)

		return

	}

	// Menyiapkan response.
	responses := make(
		[]map[string]interface{},
		0,
		len(orders),
	)

	// Convert Order menjadi response.
	for _, order := range orders {

		responses = append(
			responses,
			map[string]interface{}{
				"id":                      order.ID,
				"order_number":            order.OrderNumber,
				"product_id":              order.ProductID,
				"product_name":            order.ProductName,
				"player_id":               order.PlayerID,
				"server_id":               order.ServerID,
				"buy_price":               order.BuyPrice,
				"sell_price":              order.SellPrice,
				"status":                  order.Status,
				"provider_transaction_id": order.ProviderTransactionID,
				"expired_at":              order.ExpiredAt,
				"paid_at":                 order.PaidAt,
				"created_at":              order.CreatedAt,
			},
		)

	}

	// Response
	utils.Success(
		data.Writer,
		http.StatusOK,
		"Orders retrieved successfully",
		responses,
	)
}

// UpdateOrderStatus mengubah status order.
//
// PUT /api/orders/:id/status
// @Summary Mengubah status order
// @Description Mengubah status order berdasarkan ID.
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID order"
// @Param request body OrderStatusRequest true "Status order"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/orders/{id}/status [put]
func (controller *OrderController) UpdateOrderStatus(
	data *gin.Context,
) {

	// Ambil ID dari URL.
	id := data.Param("id")

	// Cari order.
	var order models.Order

	result := controller.DB.First(
		&order,
		id,
	)

	if result.Error != nil {

		utils.Error(
			data.Writer,
			http.StatusNotFound,
			"Order not found",
		)

		return
	}

	// Request body.
	var request OrderStatusRequest

	// Decode JSON.
	if err := data.ShouldBindJSON(&request); err != nil {

		utils.Error(
			data.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)

		return
	}

	// Validasi status.
	if request.Status != "PENDING" &&
		request.Status != "PROCESSING" &&
		request.Status != "SUCCESS" &&
		request.Status != "FAILED" {

		utils.Error(
			data.Writer,
			http.StatusBadRequest,
			"Invalid order status",
		)

		return
	}

	// Update status.
	order.Status = request.Status

	// Simpan perubahan.
	result = controller.DB.Save(&order)

	if result.Error != nil {

		utils.Error(
			data.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)

		return
	}

	// Response.
	response := map[string]interface{}{
		"id":           order.ID,
		"order_number": order.OrderNumber,
		"status":       order.Status,
	}

	utils.Success(
		data.Writer,
		http.StatusOK,
		"Order status updated successfully",
		response,
	)
}
