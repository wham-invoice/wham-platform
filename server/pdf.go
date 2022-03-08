package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetPDFhandler returns a handler for the PDF download. Using Content type (?) we can force the browser to download the file.
func GetPDFhandler(c *gin.Context) {
	fmt.Printf("hello")
}
