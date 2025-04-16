package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func svgToPNG(svgData []byte, outputPath string) error {
	// Save the SVG data to a temporary file
	tmpSVGFile, err := os.CreateTemp("", "temp-*.svg")
	if err != nil {
		return fmt.Errorf("error creating temporary SVG file: %w", err)
	}
	defer os.Remove(tmpSVGFile.Name())
	defer tmpSVGFile.Close()

	if _, err := tmpSVGFile.Write(svgData); err != nil {
		return fmt.Errorf("error writing SVG data to file: %w", err)
	}

	fmt.Println("[INFO] Converting SVG to PNG...")
	cmd := exec.Command("rsvg-convert", "-f", "png", "-o", outputPath, tmpSVGFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error converting SVG to PNG: %w, output: %s", err, string(output))
	}
	fmt.Println("[SUCCESS] SVG converted to PNG:", outputPath)
	return nil
}

func convertSVGToPDF(c *gin.Context) {
	svgData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		msg := fmt.Sprintf("Error reading request body: %v", err)
		fmt.Println("[ERROR]", msg)
		c.String(http.StatusBadRequest, msg)
		return
	}
	defer c.Request.Body.Close()

	if len(svgData) == 0 {
		msg := "SVG data is empty"
		fmt.Println("[ERROR]", msg)
		c.String(http.StatusBadRequest, msg)
		return
	}

	tmpDir, err := os.MkdirTemp("", "pdfgen-*")
	if err != nil {
		msg := fmt.Sprintf("Error creating temp dir: %v", err)
		fmt.Println("[ERROR]", msg)
		c.String(http.StatusInternalServerError, msg)
		return
	}
	defer os.RemoveAll(tmpDir)

	pngPath := filepath.Join(tmpDir, "image.png")

	err = svgToPNG(svgData, pngPath)
	if err != nil {
		msg := fmt.Sprintf("Error in svgToPNG: %v", err)
		fmt.Println("[ERROR]", msg)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	pngFile, err := os.Open(pngPath)
	if err != nil {
		msg := fmt.Sprintf("Error opening PNG: %v", err)
		fmt.Println("[ERROR]", msg)
		c.String(http.StatusInternalServerError, msg)
		return
	}
	defer pngFile.Close()

	var pdfBuf bytes.Buffer

	fmt.Println("[INFO] Importing PNG into PDF...")
	err = api.ImportImages(nil, &pdfBuf, []io.Reader{pngFile}, nil, nil)
	if err != nil {
		msg := fmt.Sprintf("Error importing PNG to PDF: %v", err)
		fmt.Println("[ERROR]", msg)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	fmt.Println("[SUCCESS] PDF created. Sending to client.")
	c.DataFromReader(http.StatusOK, int64(pdfBuf.Len()), "application/pdf", bytes.NewReader(pdfBuf.Bytes()), map[string]string{
		"Content-Disposition": `attachment; filename="converted.pdf"`,
	})
}

func main() {
	router := gin.Default()
	// Serve the HTML file
	router.StaticFile("/", "./index.html")

	router.POST("/svg-to-pdf", convertSVGToPDF)

	port := "8080"
	fmt.Println("[INFO] Server running on port", port)
	err := router.Run(":" + port)
	if err != nil {
		fmt.Println("[ERROR] Failed to start server:", err)
	}
}
