package categoryhandler

import (
	"log"
	"net/http"
	"todoapp/internal/param"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) createCategory(c echo.Context) error {
	var req param.CreateCategoryRequest
	log.Printf("🔵 Create category request received")

	if err:=c.Bind(&req);err!=nil{
			log.Printf("❌ Bind error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	claims:=claim.GetClaimsFromEchoContext(c)
	req.UserID=claims.UserID

	log.Printf("🔵 Creating category for user: %d, name: %s", req.UserID, req.Name)
	
	if fieldErrors,err:=h.categoryValidator.ValidateCreateCategory(c.Request().Context(),req);err!=nil{
		log.Printf("❌ Validationerror: %v, fields: %v", err, fieldErrors)
		msg,code:=httpmsgerrorhandler.Error(err)
		return c.JSON(code,echo.Map{
			"message":msg,
			"errors":fieldErrors,
		})
	}

	log.Printf("🔵 Validation passed, calling service... %v:\n",req)

	resp,err:=h.categorySvc.CreateCategory(c.Request().Context(),req)
	if err != nil {
		log.Printf("❌ Service error: %v", err)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}
	log.Printf("✅ Category created successfully: ID=%d", resp.Category.ID)

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "category created successfully",
		"data":    resp,
	})
}