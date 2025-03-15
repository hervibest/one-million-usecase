package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hervibest/one-million-usecase/internal/delivery/http/controller"
)

func SetupNewUploadRoute(app *fiber.App, uploadController controller.UploadController) {
	app.Post("/upload", uploadController.UploadFile)
}
