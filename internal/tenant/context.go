package tenant

import "github.com/gin-gonic/gin"

const tenantKey = "tenant_id"

func SetTenant(c *gin.Context, tenantID int64) {
	c.Set(tenantKey, tenantID)
}

func GetTenant(c *gin.Context) (int64, bool) {
	val, ok := c.Get(tenantKey)
	if !ok {
		return 0, false
	}
	id, ok := val.(int64)
	return id, ok
}

func CurrentTenantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tenantID, ok := GetTenant(c); ok {
			c.JSON(200, gin.H{"tenant_id": tenantID})
			return
		}
		c.JSON(404, gin.H{"error": "tenant not found"})
	}
}
