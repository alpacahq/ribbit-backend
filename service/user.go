package service

import (
	"net/http"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository/user"
	"github.com/alpacahq/ribbit-backend/request"

	"github.com/gin-gonic/gin"
)

// User represents the user http service
type User struct {
	svc *user.Service
}

// UserRouter declares the orutes for users router group
func UserRouter(svc *user.Service, r *gin.RouterGroup) {
	u := User{
		svc: svc,
	}
	ur := r.Group("/users")
	ur.GET("", u.list)
	ur.GET("/:id", u.view)
	ur.PATCH("/:id", u.update)
	ur.DELETE("/:id", u.delete)
}

type listResponse struct {
	Users []model.User `json:"users"`
	Page  int          `json:"page"`
}

func (u *User) list(c *gin.Context) {
	p, err := request.Paginate(c)
	if err != nil {
		return
	}
	result, err := u.svc.List(c, &model.Pagination{
		Limit: p.Limit, Offset: p.Offset,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, listResponse{
		Users: result,
		Page:  p.Page,
	})
}

func (u *User) view(c *gin.Context) {
	id, err := request.ID(c)
	if err != nil {
		return
	}
	result, err := u.svc.View(c, id)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (u *User) update(c *gin.Context) {
	updateUser, err := request.UserUpdate(c)
	if err != nil {
		return
	}
	userUpdate, err := u.svc.Update(c, &user.Update{
		ID:        updateUser.ID,
		FirstName: updateUser.FirstName,
		LastName:  updateUser.LastName,
		Mobile:    updateUser.Mobile,
		Phone:     updateUser.Phone,
		Address:   updateUser.Address,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, userUpdate)
}

func (u *User) delete(c *gin.Context) {
	id, err := request.ID(c)
	if err != nil {
		return
	}
	if err := u.svc.Delete(c, id); err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
