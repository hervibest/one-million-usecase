package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/hervibest/one-million-usecase/internal/config"
	"github.com/hervibest/one-million-usecase/internal/delivery/http/controller"
	"github.com/hervibest/one-million-usecase/internal/delivery/http/route"
	"github.com/hervibest/one-million-usecase/internal/repository"
	"github.com/hervibest/one-million-usecase/internal/usecase"
)

func main() {
	db := config.NewPostgresDatabase()
	domainRepository := repository.NewDomainRepository(db)
	uploadUseCase := usecase.NewUploadUseCase(domainRepository)
	uploadController := controller.NewUploadController(uploadUseCase)

	app := config.NewFiber()
	route.SetupNewUploadRoute(app, uploadController)

	go func() {
		interuptSignal := make(chan os.Signal, 1)
		signal.Notify(interuptSignal, os.Interrupt)
		<-interuptSignal
		if err := app.Shutdown(); err != nil {
			log.Fatalf("failed to shutdown fiber app %s", err.Error())
		}
	}()

	if err := app.Listen(fmt.Sprintf(":%d", 5000)); err != nil {
		log.Printf("Cannot run server %v", err)
	}

}
