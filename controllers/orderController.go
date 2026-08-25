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

// CreateOrder membuat order baru.
//
// POST /api/orders

func (controller *OrderController) CreateOrder(data *gin.Context) {

	// ========================================
	// 1. Request body
	// ========================================

	var request struct {
		ProductID uint   `json:"product_id"`
		PlayerID  string `json:"player_id"`
		ServerID  string `json:"server_id"`
	}

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

// GetOrderByID mengambil detail order berdasarkan ID.
//
// GET /api/orders/:id

func (controller *OrderController) GetOrderByID(data *gin.Context) {

	// Ambil ID dari URL.
	id := data.Param("id")

	// Siapkan object Order.
	var order models.Order

	// Cari order berdasarkan ID.
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

func (controller *OrderController) GetOrders(data *gin.Context) {

	// Menampung semua order.
	var orders []models.Order

	// Ambil data dari database.
	result := controller.DB.Find(
		&orders,
	)

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
	var request struct {
		Status string `json:"status"`
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
