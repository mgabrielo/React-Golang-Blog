package routes

import (
	"github/mgabrielo/React-Golang-Blog/controller"
	"github/mgabrielo/React-Golang-Blog/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register/", controller.Register)
	app.Post("/api/login/", controller.Login)
	app.Use(middleware.IsAuthenticate)
	app.Post("/api/create-post/", controller.CreatePost)
	app.Post("/api/image-upload", controller.UploadImage)
	app.Get("/api/all-post/", controller.AllPost)
	app.Get("/api/unique-post/:id", controller.UniquePost)
	app.Put("/api/update-post/:id", controller.UpdatePost)
	app.Delete("/api/delete-post/:id", controller.DeletePost)
	app.Static("/api/uploads", "./uploads")
}
