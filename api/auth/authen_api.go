package auth

import (
	oAuth "LongTM/basic/o/auth"
	"LongTM/basic/x/rest"
	"github.com/gin-gonic/gin"
)

type AuthenServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewAuthenServer(parent *gin.RouterGroup, name string) {
	var s = AuthenServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/login", s.handleLogin)
	s.POST("/register", s.handleRegister)
}

func (s *AuthenServer) handleLogin(ctx *gin.Context) {
	var u *oAuth.LoginUser
	rest.AssertNil(ctx.BindJSON(&u))
	s.SendData(ctx, oAuth.LoginApp(u))
}

func (s *AuthenServer) handleRegister(ctx *gin.Context) {
	var u *oAuth.User
	rest.AssertNil(ctx.BindJSON(&u))
	s.SendData(ctx, oAuth.CreateUser(u))
}
