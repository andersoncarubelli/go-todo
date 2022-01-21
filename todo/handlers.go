package todo

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type TodoHandler struct {
	repository *TodoRepository
}

func (hanlder *TodoHandler) GetAll(c *fiber.Ctx) error {
	var todos []Todo = hanlder.repository.FindAll()
	return c.JSON(todos)
}

func (handler *TodoHandler) Get(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	todo, err := handler.repository.Find(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status": fiber.StatusNotFound,
			"error":  err,
		})
	}

	return c.JSON(todo)
}

func (handler *TodoHandler) Create(c *fiber.Ctx) error {
	data := new(Todo)

	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"error":   err,
		})
	}

	item, err := handler.repository.Create(*data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Failed creating item",
			"error":   err,
		})
	}

	return c.JSON(item)
}

func (handler *TodoHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Item not found",
			"error":   err,
		})
	}

	todo, err := handler.repository.Find(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Item not found",
		})
	}

	todoData := new(Todo)

	if err := c.BodyParser(todoData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"data":    err,
		})
	}

	todo.Name = todoData.Name
	todo.Description = todoData.Description
	todo.Status = todoData.Status

	item, err := handler.repository.Save(todo)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error updating todo",
			"error":   err,
		})
	}

	return c.JSON(item)
}

func (handler *TodoHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Failed to delete todo",
			"error":   err,
		})
	}

	rowsAffected := handler.repository.Delete(id)
	statusCode := fiber.StatusCreated

	if rowsAffected == 0 {
		statusCode = fiber.StatusBadRequest
	}

	return c.Status(statusCode).JSON(nil)
}

func NewTodoHandler(repository *TodoRepository) *TodoHandler {
	return &TodoHandler{
		repository: repository,
	}
}

func Register(router fiber.Router, database *gorm.DB) {
	//database.AutoMigrate(&Todo{})

	todoRepository := NewTodoRepository(database)
	todoHandler := NewTodoHandler(todoRepository)

	movieRouter := router.Group("/todo")
	movieRouter.Get("/", todoHandler.GetAll)
	movieRouter.Get("/:id", todoHandler.Get)
	movieRouter.Put("/:id", todoHandler.Update)
	movieRouter.Post("/", todoHandler.Create)
	movieRouter.Delete("/:id", todoHandler.Delete)
}
