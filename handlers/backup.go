package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func BackupConfig(c *gin.Context) {
	userID := c.GetInt("userID")
	path := filepath.Join("backups", fmt.Sprintf("%d_config.json", userID))
	data, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "config not found"})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func RestoreConfig(c *gin.Context) {
	userID := c.GetInt("userID")

	path := filepath.Join("backups", fmt.Sprintf("%d_config.json", userID))
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	os.WriteFile(path, data, 0644)
	c.JSON(http.StatusOK, gin.H{"status": "config saved"})
}

func BackupState(c *gin.Context) {
	userID := c.GetInt("userID")

	path := filepath.Join("backups", fmt.Sprintf("%d_state.json", userID))
	data, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "state not found"})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func RestoreState(c *gin.Context) {
	userID := c.GetInt("userID")

	path := filepath.Join("backups", fmt.Sprintf("%d_state.json", userID))
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	os.WriteFile(path, data, 0644)
	c.JSON(http.StatusOK, gin.H{"status": "state saved"})
}
