package repository

import (
	"github.com/magicnana999/im/entity"
)

var DefaultSequenceRepository = &SequenceRepository{}

type SequenceRepository struct {
}

func InitSequenceRepository() *SequenceRepository {
	initGorm()
	return DefaultSequenceRepository
}

func (s *SequenceRepository) SelectById(appId, cId string) (*entity.Sequence, error) {
	var entry entity.Sequence
	err := db.Where("app_id = ? AND c_id = ?", appId, cId).First(&entry).Error
	return &entry, err
}

func (s *SequenceRepository) Insert(entry *entity.Sequence) error {
	return db.Create(entry).Error
}

func (s *SequenceRepository) Update(entry *entity.Sequence) error {
	return db.Save(entry).Error
}

func (s *SequenceRepository) Delete(entry *entity.Sequence) error {
	return db.Delete(entry).Error
}
