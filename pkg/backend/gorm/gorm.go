package gorm

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/influenzanet/survey-repository/pkg/backend"
	"github.com/influenzanet/survey-repository/pkg/models"
	"github.com/klauspost/compress/zstd"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Create a writer that caches compressors.
// For this operation type we supply a nil Reader.
var encoder, _ = zstd.NewWriter(nil)

func Compress(src []byte) string {
	cz := encoder.EncodeAll(src, make([]byte, 0, len(src)))
	return base64.StdEncoding.EncodeToString(cz)
}

func DecompressStd(in io.Reader, out io.Writer) error {
	d, err := zstd.NewReader(in)
	if err != nil {
		return err
	}
	defer d.Close()

	// Copy content...
	_, err = io.Copy(out, d)
	return err
}

func Decompress(src string) ([]byte, error) {
	cz, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer // A Buffer needs no initialization.
	err = DecompressStd(bytes.NewReader(cz), bufio.NewWriter(&b))
	return b.Bytes(), err
}

type Backend interface {
}

type GormBackedConfig struct {
	DSN string
}

type GormBackend struct {
	config GormBackedConfig
	db     *gorm.DB
}

func NewGormBackend(config GormBackedConfig) *GormBackend {

	return &GormBackend{
		config: config,
	}

}

func (gb *GormBackend) Start() error {

	cfg, err := ParseDSN(gb.config.DSN)
	if err != nil {
		return err
	}
	if cfg.Driver != "sqlite" {
		return fmt.Errorf("database driver '%s' is not available", cfg.Driver)
	}

	db, err := gorm.Open(sqlite.Open(cfg.Connexion), &gorm.Config{})
	if err != nil {
		return err
	}
	gb.db = db
	return nil
}

func (gb *GormBackend) ImportSurvey(meta models.SurveyMetadata, data []byte) (uint, error) {
	sz := Compress(data)

	meta.SurveyData = models.SurveyData{Survey: string(sz)}
	result := gb.db.Create(&meta)
	if result.Error != nil {
		return 0, result.Error
	}
	return meta.ID, nil
}

func (gb *GormBackend) FindSurvey(meta models.SurveyMetadata) (uint, error) {
	sd := models.SurveyMetadata{
		Namespace:  meta.Namespace,
		PlatformID: meta.PlatformID,
		Descriptor: models.SurveyDescriptor{
			Name:      meta.Descriptor.Name,
			VersionID: meta.Descriptor.VersionID,
		},
	}
	r := models.DBId{}
	result := gb.db.Model(&sd).Where(&sd).Select("id").First(&r)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return r.ID, nil
}

func rangeFilter(db *gorm.DB, field string, filter backend.RangeFilter) {
	if filter.From > 0 {
		db.Where(fmt.Sprintf("%s > ?", field), filter.From)
	}
	if filter.To > 0 {
		db.Where(fmt.Sprintf("%s < ?", field), filter.From)
	}
}

func (gb *GormBackend) GetSurveys(namespace uint, filters backend.SurveyFilter) (backend.PaginatedResult[models.SurveyMetadata], error) {

	db := gb.db

	if len(filters.Platforms) > 0 {
		db.Where("platforms IN ?", filters.Platforms)
	}

	rangeFilter(db, "imported_at", filters.ImporterAt)
	rangeFilter(db, "descriptor_published", filters.Published)

	page := backend.PaginatedResult[models.SurveyMetadata]{}

	db.Count(&page.Total)

	if filters.Limit > 0 {
		page.Limit = int64(filters.Limit)
		db.Limit(filters.Limit)
	}

	if filters.Offset > 0 {
		page.Offset = int64(filters.Offset)
		db.Offset(filters.Offset)
	}

	result := db.Find(&page.Data)
	if result.Error != nil {
		return page, result.Error
	}
	return page, nil
}

func (gb *GormBackend) GetSurveyMeta(id uint) (models.SurveyMetadata, error) {
	sd := models.SurveyMetadata{ID: id}
	result := gb.db.Model(sd).First(&sd)
	if result.Error != nil {
		return sd, result.Error
	}
	return sd, nil
}

func (gb *GormBackend) GetSurveyData(id uint, decompress bool) ([]byte, error) {
	sd := models.SurveyData{SurveyID: id}
	result := gb.db.Model(sd).First(&sd)
	if result.Error != nil {
		return nil, result.Error
	}
	data, err := Decompress(sd.Survey)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (gb *GormBackend) GetNamespaces() ([]models.Namespace, error) {
	var nn []models.Namespace
	result := gb.db.Find(&nn)
	if result.Error != nil {
		return nil, result.Error
	}
	return nn, nil
}

func (gb *GormBackend) CreateNamespace(name string) (uint, error) {
	ns := models.Namespace{Name: name}
	result := gb.db.Create(&ns)
	if result.Error != nil {
		return 0, result.Error
	}
	return ns.ID, nil
}

func (gb *GormBackend) Migrate() error {
	gb.db.AutoMigrate(&models.Namespace{})
	gb.db.AutoMigrate(&models.SurveyMetadata{})
	gb.db.AutoMigrate(&models.SurveyData{})
	return nil
}
