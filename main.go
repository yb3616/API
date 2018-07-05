package main

import (
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	// 数据库初始化
	var db *gorm.DB = initDB()
	defer db.Close()

	r := gin.Default()

	store := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://127.0.0.1", "http://localhost", "http://127.0.0.1:80", "http://localhost:80"}
	r.Use(cors.New(config))

	a := gormadapter.NewAdapter("sqlite3", "./data/data.db")
	e := casbin.NewEnforcer("./auth/rbac_with_deny_model.conf", a)
	e.LoadPolicy()
	policyInit(e)
	e.SavePolicy()
	r.Use(newAuthorizer(e))

	// 每次运行程序检查Roles表
	var count int
	db.Table("roles").Count(&count)
	if count == 0 {
		tx := db.Begin()
		if err := addRole(tx, &Role{RoleName: "administrator"}); err != nil {
			tx.Rollback()
			return
		}
		if err := addRole(tx, &Role{RoleName: "anonymous"}); err != nil {
			tx.Rollback()
			return
		}
		if err := addRole(tx, &Role{RoleName: "user"}); err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}

	addAdmin(db, e)

	// 路由
	route(r, db, e)

	r.Run("127.0.0.1:8080")
}
