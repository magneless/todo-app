package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Handler) createList(c *gin.Context) {
	id, _ := c.Get(userCtx)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (s *Handler) getAllLists(c *gin.Context) {

}

func (s *Handler) getListById(c *gin.Context) {

}

func (s *Handler) updateList(c *gin.Context) {

}

func (s *Handler) deleteList(c *gin.Context) {

}
