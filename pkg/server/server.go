package server

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/influenzanet/survey-repository/pkg/config"
	"github.com/influenzanet/survey-repository/pkg/manager"
	"github.com/influenzanet/survey-repository/pkg/models"
	"github.com/influenzanet/survey-repository/pkg/surveys"
	"github.com/influenzanet/survey-repository/pkg/utils"

	fiber "github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

type HttpServer struct {
	app     *fiber.App
	config  *config.AppConfig
	manager *manager.Manager
	start   time.Time
}

func NewHttpServer(config *config.AppConfig, manager *manager.Manager) *HttpServer {
	return &HttpServer{config: config, manager: manager}
}

func (server *HttpServer) HomeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"Status":  "ok",
		"Started": server.start,
	})
}

func (server *HttpServer) NamespacesHandler(c *fiber.Ctx) error {
	namespaces := server.manager.GetNamespaces()
	return c.JSON(namespaces)
}

func (server *HttpServer) ImportHandler(c *fiber.Ctx) error {
	namespace := c.Params("namespace")

	ns := server.manager.GetNamespaceID(namespace)
	if ns == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Unknown namespace")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}
	files := form.File["survey"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file",
		})
	}
	file, err := files[0].Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read survey data",
		})
	}
	survey, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to read survey data",
		})
	}

	fmt.Println("survey")
	fmt.Println(survey)

	platform := c.FormValue("platform")

	descriptor, err := surveys.ExtractSurveyMetadata([]byte(survey))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}

	username := string(c.Locals("_user").(string))

	meta := models.SurveyMetadata{
		Namespace:  ns,
		PlatformID: platform,
		ImportedAt: time.Now().Unix(),
		ImportedBy: username,
		Descriptor: *descriptor,
	}

	var id uint

	id, err = server.manager.ImportSurvey(meta, []byte(survey))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}
	return c.JSON(fiber.Map{
		"id": id,
	})
}

func (server *HttpServer) SurveyDataHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}
	data, err := server.manager.GetSurveyData(uint(id), true)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}
	c.SendString(string(data))
	return nil
}

func (server *HttpServer) SurveyMetaHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}
	data, err := server.manager.GetSurveyMeta(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}
	return c.JSON(data)
}

func (server *HttpServer) BasicAuthorizer(user, password string) bool {
	hash, ok := server.config.Users[user]
	if !ok {
		return false
	}
	check, err := utils.CheckPassword(hash, password)
	if err != nil {
		log.Printf("Error checking password hash : %s", err)
	}
	return check
}

func (server *HttpServer) Start() error {

	app := fiber.New()

	fiberlog.SetLevel(fiberlog.LevelInfo)

	server.app = app
	//server.instance = uuid.NewString()
	server.start = time.Now()

	authMiddleware := basicauth.New(basicauth.Config{
		Users:      nil,
		Realm:      "Forbidden",
		Authorizer: server.BasicAuthorizer,
		Unauthorized: func(c *fiber.Ctx) error {
			c.JSON(fiber.Map{
				"Status": "Unauthorized",
			})
			return nil
		},
		ContextUsername: "_user",
		ContextPassword: "_pass",
	})

	app.Get("/", server.HomeHandler)
	app.Get("/_/ns", server.NamespacesHandler)
	app.Post("/import/:namespace", authMiddleware, server.ImportHandler)
	app.Get("/survey/:id/data", server.SurveyDataHandler)
	app.Get("/survey/:id", server.SurveyMetaHandler)

	return app.Listen(server.config.Host)
}
