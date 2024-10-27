package application

import (
	"net/http"
	"github.com/gin-gonic/gin"
    er "github.com/betterreads/internal/pkg/errors"
    "log"
) 

func ErrorMiddleware(c *gin.Context){
    logger := log.New(log.Writer(), "[Error] ", log.LstdFlags)
    c.Next() // This is a placeholder for the next middleware

    err := c.Errors.Last()
    if  err != nil {
        switch e := err.Err.(type) {
        case *er.ErrorDetails:
            if e.Status == http.StatusInternalServerError {
                c.JSON(e.Status, gin.H{"error": "Internal Server Error"}) 
                logger.Print(e.Error())
            } else {
                er.SendError(c, e)
            }
        case *er.ErrorDetailsWithParams:
            if e.Status == http.StatusInternalServerError {
                c.JSON(e.Status, gin.H{"error": "Internal Server Error"}) 
                logger.Print(e.Error())
            } else {
                er.SendErrorWithParams(c, e)
            }
        default: 
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
            logger.Print(err.Error())
        }
    }
}


