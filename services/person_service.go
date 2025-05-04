package person_service

import (
	"fmt"
	"future_today/internal/addition"
	"future_today/internal/storage"
	"future_today/models"
)

type PersonService struct {
	add *addition.Addition
	orm *storage.OrmRequestManager
}

func NewPersonService(add *addition.Addition, orm *storage.OrmRequestManager) *PersonService {
	return &PersonService{add: add, orm: orm}
}

func (s *PersonService) CreatePerson(req *models.CreatePersonRequest) (*models.Person, error) {

	age, gender, nation, err := s.add.GetAdditionAsync(req.Name)
	if err != nil {
		return nil, fmt.Errorf("error adding person data: %v", err)
	}

	person := &models.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nation,
		IsActive:    true,
	}

	err = s.orm.Create(person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (s *PersonService) GetPerson(id uint) (*models.Person, error) {
	person, err := s.orm.GetByID(id)
	if err != nil {
		return nil, err
	}
	return person, nil
}

func (s *PersonService) GetAllPersons(
	limit, offset int,
	name, surname, patronymic *string,
	minAge, maxAge *int,
	gender, nationality *string,
) ([]models.Person, error) {
	persons, err := s.orm.GetAll(name, surname, patronymic, gender, nationality, minAge, maxAge, limit, offset)
	if err != nil {
		return nil, err
	}
	return persons, nil
}

func (s *PersonService) UpdatePerson(id uint, upd *models.UpdatePersonRequest) (*models.Person, error) {
	person, err := s.orm.GetByID(id)
	if err != nil {
		return nil, err
	}

	if upd.Name != nil {
		person.Name = *upd.Name
	}
	if upd.Surname != nil {
		person.Surname = *upd.Surname
	}
	if upd.Patronymic != nil {
		person.Patronymic = *upd.Patronymic
	}
	if upd.Age != nil {
		person.Age = *upd.Age
	}
	if upd.Gender != nil {
		person.Gender = *upd.Gender
	}
	if upd.Nationality != nil {
		person.Nationality = *upd.Nationality
	}

	err = s.orm.Update(person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (s *PersonService) DeletePerson(id uint) error {
	err := s.orm.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
