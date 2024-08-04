package user

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/guilhermemena/agenda-zap-server/storage"
	"github.com/guilhermemena/agenda-zap-server/types"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userStore storage.UserStorage
	validate  *validator.Validate
}

func NewUserHandler(userStore storage.UserStorage) *UserHandler {
	validate := validator.New()
	return &UserHandler{
		userStore: userStore,
		validate:  validate,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	user := new(types.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if err := h.validate.Struct(user); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)
		for _, e := range errs {
			var message string
			switch e.Tag() {
			case "required":
				message = fmt.Sprintf("%s é obrigatório", e.Field())
			case "email":
				message = "O email fornecido é inválido"
			case "min":
				message = fmt.Sprintf("%s deve ter no mínimo %s caracteres", e.Field(), e.Param())
			case "max":
				message = fmt.Sprintf("%s deve ter no máximo %s caracteres", e.Field(), e.Param())
			default:
				message = e.Error()
			}
			errorMessages[e.Field()] = message
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errorMessages})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Tivemos um problema ao criar o usuário"})
	}
	user.Password = string(hashedPassword)

	createdUser, err := h.userStore.Create(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]interface{}{
		"message": "Usuário criado com sucesso",
		"user":    createdUser,
	})
}
