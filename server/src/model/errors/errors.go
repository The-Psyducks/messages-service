package errors

import "github.com/gin-gonic/gin"

func NewErrorMessage(ctx *gin.Context, err error) {
	ctx.JSON(500, gin.H{"TODO: implement nice errors": err.Error()})
}
