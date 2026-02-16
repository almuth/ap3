package repository

import (
	"ahadpos-go/pkg/database"
	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

func NewBaseRepository() *BaseRepository {
	return &BaseRepository{DB: database.DB}
}

func (r *BaseRepository) GetDB() *gorm.DB {
	return r.DB
}

func (r *BaseRepository) Create(model interface{}) error {
	return r.DB.Create(model).Error
}

func (r *BaseRepository) Update(model interface{}) error {
	return r.DB.Save(model).Error
}

func (r *BaseRepository) Delete(model interface{}) error {
	return r.DB.Delete(model).Error
}

func (r *BaseRepository) First(model interface{}, conditions ...interface{}) error {
	return r.DB.First(model, conditions...).Error
}

func (r *BaseRepository) Find(model interface{}, conditions ...interface{}) error {
	return r.DB.Find(model, conditions...).Error
}

func (r *BaseRepository) WithPreloads(preloads []string) *gorm.DB {
	db := r.DB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

func (r *BaseRepository) WithPagination(page, pageSize int) *gorm.DB {
	offset := (page - 1) * pageSize
	return r.DB.Offset(offset).Limit(pageSize)
}

func (r *BaseRepository) Count(model interface{}) (int64, error) {
	var count int64
	err := r.DB.Model(model).Count(&count).Error
	return count, err
}
