package user

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/guilhermemena/agenda-zap-server/cmd/api/middleware"
	"github.com/guilhermemena/agenda-zap-server/cmd/configs"
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

func (h *UserHandler) HandleRegister(c *fiber.Ctx) error {
	user := new(types.RegisterUserPayload)

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

	_, err := h.userStore.GetByEmail(c.Context(), user.Email)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Já existe um usuário com esse email"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Tivemos um problema ao criar o usuário"})
	}
	user.Password = string(hashedPassword)

	createdUser, err := h.userStore.Create(c.Context(), &types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]interface{}{
		"message": "Usuário criado com sucesso",
		"user":    createdUser,
	})
}

func (h *UserHandler) HandleLogin(c *fiber.Ctx) error {
	user := new(types.LoginUserPayload)

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
			default:
				message = e.Error()
			}
			errorMessages[e.Field()] = message
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errorMessages})
	}

	existingUser, err := h.userStore.GetByEmail(c.Context(), user.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email ou senha incorretos"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email ou senha incorretos"})
	}

	secret := []byte(configs.Envs.JWTSecret)
	token, err := middleware.CreateJWT(secret, existingUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login efetuado com sucesso", "token": token})
}

func (h *UserHandler) HandleMe(c *fiber.Ctx) error {
	userID := c.Locals(string(middleware.UserKey)).(string)
	user, err := h.userStore.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Usuário não encontrado"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": user})
}
