package main

import "github.com/gofiber/fiber/v2"

func init() {

}

type App struct {
	fb *fiber.App
}

func NewApp() *App {
	return &App{}
}
