package handler

import "github.com/gin-gonic/gin"

func GinHandler() *gin.Engine {
	h := NewHandler()
	r := gin.Default()
	r.GET("/ping", h.Ping)
	r.POST("/insert_test_data", h.InsertTestData)
	r.POST("/campaign", h.CreateCampaign)
	r.POST("/campaign/:campaign_id/start", h.StartCampaign)
	r.POST("/campaign/:campaign_id/claim", h.ClaimRedEnvelope)

	return r
}
