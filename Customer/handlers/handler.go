package handlers

import (
	"CustomerOrderApi/Customer/entities"
	"CustomerOrderApi/Customer/repositories"
	"CustomerOrderApi/Customer/services"
	"fmt"
	"github.com/erenkaratas99/COApiCore/pkg"
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"github.com/erenkaratas99/COApiCore/shared/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	service    *services.Service
	echo       *echo.Echo
	repository *repositories.Repository
}

func NewHandler(service *services.Service, echo *echo.Echo, repository *repositories.Repository) *Handler {

	return &Handler{service, echo, repository}
}

func (h *Handler) InitEndpoints() {
	e := h.echo
	g := e.Group("/customer")
	g.POST("/", h.CreateCustomer)
	g.GET("/", h.GetAllCustomer)
	g.GET("/:customerId", h.GetById)
	//g.DELETE("/:customerId", h.DeleteById)

}

func (h *Handler) CreateCustomer(c echo.Context) error {
	customerReq := entities.CustomerRequestModel{}
	err := c.Bind(&customerReq)
	if err != nil {
		fmt.Println(err)
		return customErrors.BindErr
	}

	//err = c.Validate(customerReq)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	insertedID, err := h.service.CreateCustomerService(&customerReq)
	if err != nil {
		fmt.Println(err)
		return err
	}

	srm := types.GetSRM(*insertedID)
	return c.JSON(http.StatusCreated, srm)
}

func (h *Handler) GetAllCustomer(c echo.Context) error {
	l := c.QueryParam("limit")
	o := c.QueryParam("offset")
	limit, offset := pkg.LimitOffsetValidation(l, o)
	customers, err := h.service.GetAllCustomersService(limit, offset)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.JSON(http.StatusOK, customers)
}

func (h *Handler) GetById(c echo.Context) error {
	id := c.Param("customerId")
	_, err := uuid.Parse(id)
	if err != nil {
		fmt.Println(err)
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr",
			"Id has not been validated.")
	}
	customerResp, err := h.service.CustomerGetById(id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.JSON(http.StatusOK, *customerResp)
}

//
//func (h *Handler) DeleteCustomer(c echo.Context) error {
//	id := c.Param("customerid")
//	_, err := uuid.Parse(id)
//	if err != nil {
//		fmt.Println(err)
//		return customErrors.NewHTTPError(http.StatusBadRequest,
//			"IdErr",
//			"Id has not been validated.")
//	}
//	corID := helper.GetCorrelationID(c)
//	err = h.service.DeleteCustomerService(id, corID)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	return nil
//}
