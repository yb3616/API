package main

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func route(r *gin.Engine, db *gorm.DB, e *casbin.Enforcer) {
	// 路由
	r.POST("/admin/role", createRole(db))
	r.GET("/admin/role", getRoles(db))
    r.PUT("/admin/role/:rid", editRoles(db))
    r.DELETE("/admin/role/:rid", deleteRole(db))

	r.POST("/admin/logon", createAdmin(db, e))
	r.GET("/user/info", getUserInfo(db, e))
	r.POST("/logon", createUser(db, e))
	r.POST("/login", login(db))
	r.GET("/logout", logout)
}
