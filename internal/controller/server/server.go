package server

import "github.com/gin-gonic/gin"

func Run(router *gin.Engine, port string) error {
	if err := router.Run(port); err != nil {
		return err
	}

	return nil
}
