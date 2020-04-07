package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * user: ZY
 * Date: 2020/4/7 17:02
 */
type ResponseNorm struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}

var (
	Param = ResponseNorm{
		Status: 10010,
		Info:   "param error",
	}
	Success = ResponseNorm{
		Status: 10000,
		Info:   "success",
	}
	Verify = ResponseNorm{
		Status: 10011,
		Info:   "verify failed",
	}
	Internal = ResponseNorm{
		Status: 10020,
		Info:   "internal error",
	}
	Permission = ResponseNorm{
		Status: 10012,
		Info:   "permission failed",
	}
)



func OK(c *gin.Context) {
	c.JSON(http.StatusOK, Success)
}

func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Internal)
}

func VerifyError(c *gin.Context) {
	c.JSON(http.StatusForbidden, Verify)
}

func ParamError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Param)
}

func PermissionError(c *gin.Context){
	c.JSON(http.StatusForbidden,Permission)
}

func HandleError(c *gin.Context, err ResponseNorm) {
	c.JSON(http.StatusOK, err)
}
