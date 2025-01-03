package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"GoMin/miniohelper"
	"github.com/labstack/echo/v4"
)

type APIHandler struct {
	Minio *miniohelper.MinioHelper
}

func NewAPIHandler(minioHelper *miniohelper.MinioHelper) *APIHandler {
	return &APIHandler{Minio: minioHelper}
}

func (h *APIHandler) UploadFile(c echo.Context) error {
	err := CheckHeader(c)
	if err != nil {
		return err
	}

	bucketName := c.Param("bucketName")

	fileUrl := c.FormValue("file_url")
	if fileUrl == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "FileUrl is required",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to read file from request",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open file",
		})
	}
	defer src.Close()

	// create temp file for persisting uploaded file
	tempFile, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create temp file",
		})
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// copy file from temp file
	if _, err := io.Copy(tempFile, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to write to temp file",
		})
	}

	// upload to minio
	err = h.Minio.UploadFile(bucketName, fileUrl, tempFile.Name())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to upload file: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"messages": "File uploaded success!",
	})
}

func (h *APIHandler) RemoveFile(c echo.Context) error {
	if err := CheckHeader(c); err != nil {
		return err
	}

	bucketName := c.Param("bucketName")

	fileUrl := c.FormValue("file_url")
	if fileUrl == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File url is required",
		})
	}

	if err := h.Minio.RemoveFile(bucketName, fileUrl); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to remove object with url: %v from minio error: %v", fileUrl, err),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Remove file success!",
	})
}

func (h *APIHandler) GetFile(c echo.Context) error {
	if err := CheckHeader(c); err != nil {
		return err
	}

	bucketName := c.Param("bucketName")

	fileUrl := c.FormValue("file_url")
	if fileUrl == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File url is required",
		})
	}

	object, err := h.Minio.GetFileStream(bucketName, fileUrl)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error:": fmt.Sprintf("Can not fetch file %v", err),
		})
	}

	return c.Stream(http.StatusOK, "application/octet-stream", object)
}
