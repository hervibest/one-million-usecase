package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hervibest/one-million-usecase/internal/usecase"
)

type UploadController interface {
	UploadFile(ctx *fiber.Ctx) error
}

type uploadController struct {
	uploadUseCase usecase.UploadUseCase
}

func NewUploadController(uploadUseCase usecase.UploadUseCase) UploadController {
	return &uploadController{uploadUseCase: uploadUseCase}
}

func (c *uploadController) UploadFile(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "incorrect file name")
	}

	// Simpan file ke disk sebelum diproses di goroutine
	tempFilePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := ctx.SaveFile(file, tempFilePath); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// Proses file di goroutine
	go func(filePath string) {
		if err := c.uploadUseCase.UploadFile(filePath); err != nil {
			log.Println("Error saat mengupload file:", err)
		}
	}(tempFilePath)

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"status": "success"})
}
