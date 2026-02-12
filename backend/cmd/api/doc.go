package main

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

// Health godoc
// @Summary		Health check
// @Description	Returns service health status. Requires DB and Redis to be reachable.
// @Tags			health
// @Produce		json
// @Success	200	{object}	HealthResponse
// @Failure	503	{object}	object	"Unhealthy"
// @Router		/health [get]
func _healthDoc() {}
