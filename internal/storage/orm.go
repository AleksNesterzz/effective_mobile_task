package storage

import (
	"future_today/models"

	"gorm.io/gorm"
)

type OrmRequestManager struct {
	db *gorm.DB
}

func NewOrmRequestManager(db *gorm.DB) *OrmRequestManager {
	return &OrmRequestManager{db: db}
}

func (orm *OrmRequestManager) Create(person *models.Person) error {
	return orm.db.Create(person).Error
}

func (orm *OrmRequestManager) GetByID(id uint) (*models.Person, error) {
	var person models.Person
	err := orm.db.First(&person).Error
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (orm *OrmRequestManager) GetAll(
	name, surname, patronymics, gender, nationality *string,
	minAge, maxAge *int,
	limit, offset int) ([]models.Person, error) {
	var persons []models.Person
	query := orm.db.Model(&models.Person{}).Where("is_active= ?", true)
	if name != nil {
		query = query.Where("name ILIKE ?", "%"+*name+"%")
	}
	if surname != nil {
		query = query.Where("surname ILIKE ?", "%"+*surname+"%")
	}
	if minAge != nil {
		query = query.Where("age >= ?", *minAge)
	}
	if maxAge != nil {
		query = query.Where("age <= ?", *maxAge)
	}
	if gender != nil {
		query = query.Where("gender ILIKE ?", "%"+*gender+"%")
	}
	if nationality != nil {
		query = query.Where("nationality ILIKE ?", "%"+*nationality+"%")
	}

	err := query.Limit(limit).Offset(offset).Find(&persons).Error
	return persons, err
}

func (orm *OrmRequestManager) Update(person *models.Person) error {
	return orm.db.Save(person).Error
}

func (orm *OrmRequestManager) Delete(id uint) error {
	return orm.db.Model(&models.Person{}).Where("id", id).Update("is_active", false).Error
}
