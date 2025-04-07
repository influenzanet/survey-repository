package manager

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
	"context"

	"github.com/influenzanet/survey-repository/pkg/backend"
	gormBackend "github.com/influenzanet/survey-repository/pkg/backend/gorm"
	"github.com/influenzanet/survey-repository/pkg/config"
	"github.com/influenzanet/survey-repository/pkg/models"
)

var (
	ErrUnknownNamespace = errors.New("unknown namespace")
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenNotFound = errors.New("token not found")
	ErrInBackend = errors.New("unable to get record")
)

type AuthResponse struct {
	Key string `json:"key"`
	TTL int64 `json:"ttl"`
}

type Manager struct {
	db         backend.Backend
	SurveyPath string
	keyTTL	   int64
	namespaces NsRegistry
	cleanupTicker *time.Ticker
}

func NewManager(config *config.AppConfig) *Manager {

	db := gormBackend.NewGormBackend(gormBackend.GormBackendConfig{DSN: config.DB.DSN, Debug: config.DB.Debug})
	cleanupTicker := time.NewTicker(config.Auth.CleanupDuration)

	return &Manager{
		db:         db,
		SurveyPath: config.SurveyPath,
		keyTTL: config.Auth.AuthKeyTTL,
		cleanupTicker: cleanupTicker,
	}
}

func (manager *Manager) StartRoutines(ctx context.Context) error {
	
	// First run on startup
	manager.cleanupKeys()
	
	// Run cleanup routine in sub routine
	go manager.cleanupHander(ctx)
	return nil
}

func (manager *Manager) Start() error {

	if manager.SurveyPath != "" {
		s, err := os.Stat(manager.SurveyPath)
		if os.IsNotExist(err) {
			return err
		}
		if !s.IsDir() {
			return errors.New("SurveyPath must be a directory")
		}
	}

	err := manager.db.Start()
	if err != nil {
		return err
	}

	err = manager.Migrate()
	if err != nil {
		return err
	}

	err = manager.loadNamespaces()
	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) loadNamespaces() error {
	var nn []models.Namespace
	nn, err := manager.db.GetNamespaces()
	if err != nil {
		return err
	}
	manager.namespaces = createRegistry(nn)
	return nil
}

func (manager *Manager) Migrate() error {
	manager.db.Migrate()
	return nil
}

func (manager *Manager) GetNamespaces() map[uint]string {
	return manager.namespaces.toName
}

func (manager *Manager) GetNamespaceID(name string) uint {
	id, ok := manager.namespaces.toID[name]
	if !ok {
		return 0
	}
	return id
}

var nameRegexp = regexp.MustCompile("^[a-z]+$")

func (manager *Manager) CreateNamespace(name string) (uint, error) {
	if !nameRegexp.MatchString(name) {
		return 0, fmt.Errorf("namespace name must be only with lowercase alpha chars, given '%s'", name)
	}
	return manager.db.CreateNamespace(name)
}

func (manager *Manager) ImportSurvey(meta models.SurveyMetadata, filePath string, survey []byte) (uint, error) {
	id, err := manager.db.ImportSurvey(meta, survey)
	if err != nil {
		return id, err
	}
	if manager.SurveyPath != "" {
		fn := fmt.Sprintf("%s/%d.json", manager.SurveyPath, id)

		if filePath == "" {
			err = os.WriteFile(fn, survey, 0666)
		} else {
			err = os.Rename(filePath, fn)
		}
		if err != nil {
			log.Printf("Error writing survey in %s", fn)
		} else {
			log.Printf("Survey imported and saved in %s", fn)
		}
	}

	return id, err
}

func (manager *Manager) FindSurvey(meta models.SurveyMetadata) (uint, error) {
	return manager.db.FindSurvey(meta)
}

func (manager *Manager) GetSurveysStats(namespace uint) ([]models.SurveyStats, error) {
	return manager.db.GetSurveysStats(namespace)
}

func (manager *Manager) GetSurveyData(id uint, decompress bool) ([]byte, error) {
	return manager.db.GetSurveyData(id, decompress)
}

func (manager *Manager) GetSurveyMeta(id uint) (models.SurveyMetadata, error) {
	return manager.db.GetSurveyMeta(id)
}

func (manager *Manager) GetSurveys(namespace uint, filters backend.SurveyFilter) (backend.PaginatedResult[models.SurveyMetadata], error) {
	return manager.db.GetSurveys(namespace, filters)
}

func (manager *Manager) CreateAuthKey(user string) (AuthResponse, error) {
	auth, err := manager.db.CreateAuthKey(user)
	if err != nil {
		return AuthResponse{}, err
	}
	return AuthResponse{Key: auth.Key, TTL: manager.keyTTL}, nil
}

func (manager *Manager) keyExpirationTime() time.Time {
	return time.Now().Add(-time.Duration(manager.keyTTL) * time.Second)
}

func (manager *Manager) FindUserFromAuthKey(key string) (string, error) {
	auth, err := manager.db.FindUserFromAuthKey(key)
	if err != nil {
		log.Printf("Error searching for key %s", err)
		return "", ErrInBackend
	}

	expires := manager.keyExpirationTime().Unix()
	
	if(auth.Created < expires) {
		return "", ErrTokenExpired
	}
	
	if(auth.User == "") {
		return "", ErrTokenNotFound
	}
	return auth.User, nil
}

func (manager *Manager) cleanupKeys() error {
	expires := manager.keyExpirationTime().Unix()
	removed, err := manager.db.CleanupKeys(expires)
	if err != nil {
		log.Printf("Error in cleanup routine : %s", err)
	} else {
		log.Printf("Auth key removed : %d", removed)
	}
	return err
}


func (manager *Manager) cleanupHander(ctx context.Context) {
	log.Printf("Starting Cleanup Hander",)
	for {
		select {
		case <-ctx.Done():
			log.Println("cleanup routine", ctx.Err())
			return

		case <-manager.cleanupTicker.C:
			log.Println("Running cleanup routine")
			manager.cleanupKeys()
		}
	}
}


type NsRegistry struct {
	toName map[uint]string
	toID   map[string]uint
}

func createRegistry(namespaces []models.Namespace) NsRegistry {
	toName := make(map[uint]string, len(namespaces))
	toID := make(map[string]uint, len(namespaces))
	for _, n := range namespaces {
		toName[n.ID] = n.Name
		toID[n.Name] = n.ID
	}
	return NsRegistry{
		toName: toName,
		toID:   toID,
	}
}
