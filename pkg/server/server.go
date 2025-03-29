package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/influenzanet/survey-repository/pkg/backend"
	"github.com/influenzanet/survey-repository/pkg/config"
	"github.com/influenzanet/survey-repository/pkg/manager"
	"github.com/influenzanet/survey-repository/pkg/models"
	"github.com/influenzanet/survey-repository/pkg/surveys"
	"github.com/influenzanet/survey-repository/pkg/utils"
	"github.com/influenzanet/survey-repository/pkg/version"
	
	fiber "github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type HttpServer struct {
	app         *fiber.App
	config      *config.AppConfig
	manager     *manager.Manager
	start       time.Time
	counter     atomic.Uint64
	storeSurvey bool
	version  	version.VersionInfo
}

// ShortVersionMeta is a shorter structure to list survey versions
type ShortVersionMeta struct {
	ID		uint `json:"id"`
	Version string `json:"version"`
	PublishedAt int64 `json:"published"`
	PlatformID string `json:"platform"`
	Name	string `json:"name"`
	ModelType  string `json:"model_type"` // Model type 'definition','preview'
}

func NewHttpServer(config *config.AppConfig, manager *manager.Manager) *HttpServer {
	return &HttpServer{config: config, manager: manager, storeSurvey: true}
}

func (server *HttpServer) HomeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"Status":  "ok",
		"Version": server.version.Tag,
		"Started": server.start,
	})
}

func (server *HttpServer) NamespacesHandler(c *fiber.Ctx) error {
	namespaces := server.manager.GetNamespaces()
	return c.JSON(namespaces)
}

func (server *HttpServer) StatsHandler(c *fiber.Ctx) error {
	namespace := c.Params("namespace")

	ns := server.manager.GetNamespaceID(namespace)
	if ns == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Unknown namespace")
	}

	stats, err := server.manager.GetSurveysStats(ns)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}
	return c.JSON(stats)
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

	platform := c.FormValue("platform")
	if platform == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Platform code must be provided",
		})
	}

	version := c.FormValue("version")
	name := c.FormValue("name")
	
	count := server.counter.Add(1)

	var fn string
	if server.storeSurvey {
		// Store survey data in temporary file. It will be renamed in cas of success with file id
		fn = fmt.Sprintf("%s/%s-%s-%x-%x.json", server.config.SurveyPath, namespace, platform, time.Now().Unix(), count)
		err = os.WriteFile(fn, survey, 0666)
		if err != nil {
			log.Printf("Error writing survey in %s", fn)
		}
	}

	descriptor, err := surveys.ExtractSurveyMetadata([]byte(survey))
	
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}

	if(descriptor.VersionID == "") {
		if(version == "") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "The survey doesnt contains version, please provide it with `version` field in the POST request",
			})
		}
	} else {
		version = descriptor.VersionID
	}

	if(descriptor.Name == "") {
		if(name == "") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "The survey doesnt contains name, please provide it with `name` field in the POST request",
			})
		}
		descriptor.Name = ""
	}

	username := string(c.Locals("_user").(string))

	modelType := ""

	if(descriptor.ModelVersion == "preview") {
		modelType = models.SurveyModelPreview
	} else {
		modelType = models.SurveyModelDefinition
	}

	meta := models.SurveyMetadata{
		Namespace:  ns,
		PlatformID: platform,
		Version: version,
		ModelType: modelType,
		ImportedAt: time.Now().Unix(),
		ImportedBy: username,
		Descriptor: *descriptor,
	}

	var id uint

	id, err = server.manager.FindSurvey(meta)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}
	if id != 0 {
		return c.Status(fiber.StatusAlreadyReported).JSON(fiber.Map{
			"id": id,
		})
	}

	id, err = server.manager.ImportSurvey(meta, fn, []byte(survey))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}

	m := ShortVersionMeta{
		ID:  id,
		Version: meta.Version,
		PublishedAt: meta.Descriptor.Published,
		PlatformID: meta.PlatformID, 
		ModelType: meta.ModelType, 
		Name: meta.Descriptor.Name,
	}

	return c.Status(fiber.StatusCreated).JSON(m)
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

func (server *HttpServer) PlatformsHandler(c *fiber.Ctx) error {
	return c.JSON(models.WellKnownPlatforms)
}

func parseCommaList(s string) []string {
	ss := strings.Split(s, ",")
	o := make([]string, 0, len(ss))
	for _, v := range ss {
		o = append(o, strings.TrimSpace(v))
	}
	return o
}

func (server *HttpServer) NamespaceSurveysFullHandler(c *fiber.Ctx) error {
	return server.loadNamespaceSurveys(c, false)
}

func (server *HttpServer) NamespaceSurveysVersionsHandler(c *fiber.Ctx) error {
	return server.loadNamespaceSurveys(c, true)
}


func (server *HttpServer) loadNamespaceSurveys(c *fiber.Ctx, onlyVersion bool) error {
	namespace := c.Params("namespace")
	filters := backend.SurveyFilter{}

	qPlatform := c.Query("platforms")
	if qPlatform != "" {
		filters.Platforms = parseCommaList(qPlatform)
	}

	qName := c.Query("names")
	if qName != "" {
		filters.Names = parseCommaList(qName)
	}

	qTypes := c.Query("types")
	if qTypes != "" {
		filters.ModelTypes = parseCommaList(qTypes)
	}

	limit := c.QueryInt("limit", 0)
	offset := c.QueryInt("offset", 0)
	if limit > 0 {
		filters.Limit = limit
		filters.Offset = offset
	} else {
		if(offset > 0) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "`offset` param can only be used whith `limit`",
			})
		}
	}

	publishedFrom := c.QueryInt("published_from", 0)
	if publishedFrom > 0 {
		filters.Published.From = int64(publishedFrom)
	}
	publishedTo := c.QueryInt("published_to", 0)
	if publishedTo > 0 {
		filters.Published.To = int64(publishedTo)
	}

	id := server.manager.GetNamespaceID(namespace)
	if id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Unknown namspace '%s'", namespace),
		})
	}

	data, err := server.manager.GetSurveys(id, filters)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("%s", err),
		})
	}
	if(onlyVersion) {

		versions := make([]ShortVersionMeta, 0, len(data.Data))
		for _, sv := range data.Data {

			m := ShortVersionMeta{
				ID:  sv.ID,
				Version: sv.Version,
				PublishedAt: sv.Descriptor.Published,
				PlatformID: sv.PlatformID, 
				ModelType: sv.ModelType, 
				Name: sv.Descriptor.Name,
			}
			versions = append(versions, m)
		}
		p := backend.PaginatedResult[ShortVersionMeta]{
			PaginateInfo: backend.PaginateInfo{
				Total: data.Total,
				Offset: data.Offset,
				Limit: data.Limit,
			},
			Data: versions,
		}
		return c.JSON(p)
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

	server.version = version.Version()

	cfg := server.config.Server

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

	ratelimiter := limiter.New(limiter.Config{
		Max:          cfg.LimiterMax,
		Expiration:     time.Duration(int64(cfg.LimiterWindow)) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
	})

	app.Get("/", server.HomeHandler)
	app.Get("/refs/platforms", server.PlatformsHandler)
	app.Get("/namespaces", server.NamespacesHandler)
	app.Get("/namespace/:namespace/surveys", server.NamespaceSurveysFullHandler)
	app.Get("/namespace/:namespace/surveys/versions", server.NamespaceSurveysVersionsHandler)
	app.Get("/namespace/:namespace/surveys/stats", server.StatsHandler)
	app.Post("/import/:namespace", ratelimiter, authMiddleware, server.ImportHandler)
	app.Get("/survey/:id/data", server.SurveyDataHandler)
	app.Get("/survey/:id", server.SurveyMetaHandler)

	return app.Listen(cfg.Host)
}
