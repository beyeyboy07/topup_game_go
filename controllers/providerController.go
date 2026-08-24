package controllers

import (
	"net/http"
	"topup_games_go/models"
	"topup_games_go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProviderController menyimpan dependency
// yang dibutuhkan Provider Controller.

type ProviderController struct {
	DB *gorm.DB
}

type ProviderRequest struct {
	Name    string `json:"name" example:"Digiflazz"`
	Code    string `json:"code" example:"DIGIFLAZZ"`
	BaseURL string `json:"base_url" example:"https://api.example.com"`
	ApiKey  string `json:"api_key" example:"secret-key"`
	Status  bool   `json:"status" example:"true"`
}

// CreateProvider membuat provider baru.
// @Summary Membuat provider
// @Description Membuat provider baru.
// @Tags providers
// @Accept json
// @Produce json
// @Param request body ProviderRequest true "Data provider"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/providers [post]
func (controller *ProviderController) CreateProvider(data *gin.Context) {

	// Request body.
	var request ProviderRequest

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
	if request.Name == "" ||
		request.Code == "" ||
		request.BaseURL == "" ||
		request.ApiKey == "" {

		utils.Error(
			data.Writer,
			http.StatusBadRequest,
			"Name, code, base_url and api_key are required",
		)

		return
	}

	// Buat Provider.
	provider := models.Provider{
		Name:    request.Name,
		Code:    request.Code,
		BaseURL: request.BaseURL,
		ApiKey:  request.ApiKey,
		Status:  request.Status,
	}

	// Simpan ke database.
	result := controller.DB.Create(
		&provider,
	)

	if result.Error != nil {

		utils.Error(
			data.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)

		return
	}

	// Response.
	response := models.ProviderResponse{
		ID:      provider.ID,
		Name:    provider.Name,
		Code:    provider.Code,
		BaseURL: provider.BaseURL,
		Status:  provider.Status,
	}

	utils.Success(
		data.Writer,
		http.StatusCreated,
		"Provider created successfully",
		response,
	)

}

// GetProviders mengambil semua provider.
//
// GET /api/providers
// @Summary Mengambil semua provider
// @Description Mengambil daftar provider tanpa mengekspos API key.
// @Tags providers
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/providers [get]
func (controller *ProviderController) GetProviders(data *gin.Context) {

	// Menampung semua provider.
	var provider []models.Provider

	// Ambil data dari database.
	result := controller.DB.Find(&provider)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusInternalServerError,
			"Failed to get providers",
		)
		return
	}

	// Menyiapkan response.

	responses := make(
		[]models.ProviderResponse,
		0,
		len(provider),
	)

	// Convert Provider -> ProviderResponse.
	//
	// ApiKey tidak dimasukkan.

	for _, provider := range provider {

		responses = append(responses, models.ProviderResponse{
			ID:      provider.ID,
			Name:    provider.Name,
			Code:    provider.Code,
			BaseURL: provider.BaseURL,
			Status:  provider.Status,
		})

	}

	// Response.
	utils.Success(
		data.Writer,
		http.StatusOK,
		"Providers retrieved successfully",
		responses,
	)

}

// GetProviderByID mengambil satu provider berdasarkan ID.
//
// GET /api/providers/:id
// @Summary Mengambil provider berdasarkan ID
// @Description Mengambil satu provider berdasarkan ID tanpa API key.
// @Tags providers
// @Produce json
// @Param id path int true "ID provider"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/providers/{id} [get]
func (controller *ProviderController) GetProviderByID(data *gin.Context) {

	// Ambil ID dari URL.
	//
	// Contoh:
	// /api/providers/1
	id := data.Param("id")

	// Siapkan object Provider.
	var provider models.Provider

	// Cari provider berdasarkan ID.
	result := controller.DB.First(
		&provider,
		id,
	)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusNotFound,
			"Provider not found",
		)

		return
	}

	// Convert ke response.

	response := models.ProviderResponse{
		ID:      provider.ID,
		Name:    provider.Name,
		Code:    provider.Code,
		BaseURL: provider.BaseURL,
		Status:  provider.Status,
	}

	// Response.
	utils.Success(
		data.Writer,
		http.StatusOK,
		"Provider retrieved successfully",
		response,
	)

}

// UpdateProvider mengubah data provider.
//
// PUT /api/providers/:id
// @Summary Mengubah provider
// @Description Mengubah data provider berdasarkan ID.
// @Tags providers
// @Accept json
// @Produce json
// @Param id path int true "ID provider"
// @Param request body ProviderRequest true "Data provider"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/providers/{id} [put]
func (controller *ProviderController) UpdateProvider(data *gin.Context) {

	// Ambil ID dari URL.
	id := data.Param("id")

	// Cari provider.
	var provider models.Provider

	result := controller.DB.First(
		&provider,
		id,
	)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusNotFound,
			"Provider not found",
		)
		return
	}

	// Request body.
	var request ProviderRequest

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
	if request.Name == "" ||
		request.Code == "" ||
		request.BaseURL == "" ||
		request.ApiKey == "" {
		utils.Error(
			data.Writer,
			http.StatusBadRequest,
			"Name, code, base_url and api_key are required",
		)

		return
	}

	// Update provider.
	provider.Name = request.Name
	provider.Code = request.Code
	provider.BaseURL = request.BaseURL
	provider.ApiKey = request.ApiKey
	provider.Status = request.Status

	// Simpan perubahan.
	result = controller.DB.Save(
		&provider,
	)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)

		return
	}

	// Response.
	response := models.ProviderResponse{
		ID:      provider.ID,
		Name:    provider.Name,
		Code:    provider.Code,
		BaseURL: provider.BaseURL,
		Status:  provider.Status,
	}

	utils.Success(
		data.Writer,
		http.StatusOK,
		"Provider updated successfully",
		response,
	)

}

// DeleteProvider menghapus provider.
//
// DELETE /api/providers/:id
// @Summary Menghapus provider
// @Description Menghapus provider secara soft delete berdasarkan ID.
// @Tags providers
// @Produce json
// @Param id path int true "ID provider"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/providers/{id} [delete]
func (controller *ProviderController) DeleteProvider(data *gin.Context) {

	// Ambil ID dari URL (Param).
	id := data.Param("id")

	// Cari provider.
	var provider models.Provider

	result := controller.DB.First(
		&provider,
		id,
	)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusNotFound,
			"Provider not found",
		)

		return
	}

	// Hapus provider.
	result = controller.DB.Delete(
		&provider,
	)

	if result.Error != nil {
		utils.Error(
			data.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)
	}

	utils.Success(
		data.Writer,
		http.StatusOK,
		"Provider deleted successfully",
		nil,
	)

}
