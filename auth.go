package main

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func newAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//	var roles []string
		//	s := sessions.Default(c)

		//	if temp := s.Get("roles"); temp == nil {
		//		roles = []string{"2"}
		//	} else {
		//		roles = temp.([]string)
		//	}

		var name string
		if err := sessions.Default(c).Get("uid"); err != nil {
			name = fmt.Sprintf("u%d", err)
		} else {
			name = "u0"
		}
		path := c.Request.URL.Path
		method := c.Request.Method

		// DEBUG
		// fmt.Println("Current Role for user_id: ", name, e.GetRolesForUser(name))

		if !e.Enforce(name, path, method) {
			c.JSON(403, gin.H{
				"msg": "Unauthorized",
			})
			c.Abort()
		}

		c.Next()
	}
}
