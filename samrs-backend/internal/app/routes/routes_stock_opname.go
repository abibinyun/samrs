package routes

import (
	"github.com/gin-gonic/gin"
)

func registerStockOpnameRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/stock-opnames", requirePerm(deps, "stock_opname:create"), handlers.StockOpname.CreateSession)
	v1.GET("/stock-opnames", requirePerm(deps, "stock_opname:read"), handlers.StockOpname.GetAllSessions)
	v1.GET("/stock-opnames/:id", requirePerm(deps, "stock_opname:read"), handlers.StockOpname.GetSessionByID)
	v1.PATCH("/stock-opnames/:id", requirePerm(deps, "stock_opname:update"), handlers.StockOpname.UpdateSession)
	v1.PATCH("/stock-opnames/:id/close", requirePerm(deps, "stock_opname:close"), handlers.StockOpname.CloseSession)
	v1.POST("/stock-opnames/:id/items", requirePerm(deps, "stock_opname:update"), handlers.StockOpname.AddItem)
	v1.GET("/stock-opnames/:id/items", requirePerm(deps, "stock_opname:read"), handlers.StockOpname.ListItems)
}
