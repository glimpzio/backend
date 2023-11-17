package misc

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

type contextKey string

const ginContextKey contextKey = "ginContextKey"

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginContextKey)
	if ginContext == nil {
		return nil, errors.New("could not retrieve gin context")
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		return nil, errors.New("gin context has wrong type")
	}

	return gc, nil
}
