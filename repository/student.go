package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type StudentRepository interface {
	FetchAll() ([]model.Student, error)
	FetchByID(id int) (*model.Student, error)
	Store(s *model.Student) error
	Update(id int, s *model.Student) error
	Delete(id int) error
	FetchWithClass() (*[]model.StudentClass, error)
}

type studentRepoImpl struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) *studentRepoImpl {
	return &studentRepoImpl{db}
}

func (s *studentRepoImpl) FetchAll() ([]model.Student, error) {
	student := []model.Student{}
	err := s.db.Find(&student).Error
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (s *studentRepoImpl) Store(student *model.Student) error {
	if result := s.db.Create(&student); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *studentRepoImpl) Update(id int, student *model.Student) error {
	err := s.db.Where("id = ?", id).First(&student).Error
	if err != nil {
		return err
	}

	err = s.db.Model(&model.Student{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":     student.Name,
		"address":  student.Address,
		"class_id": student.ClassId,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func (s *studentRepoImpl) Delete(id int) error {

	err := s.db.Where("id = ?", id).First(&model.Student{}).Error
	if err != nil {
		return err
	}

	err = s.db.Delete(&model.Student{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *studentRepoImpl) FetchByID(id int) (*model.Student, error) {
	var student model.Student
	err := s.db.Where("id = ?", id).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (s *studentRepoImpl) FetchWithClass() (*[]model.StudentClass, error) {
	result := []model.StudentClass{}
	s.db.Table("students").
		Select("students.name AS name, students.address AS address, classes.name AS class_name, classes.professor AS professor, classes.room_number AS room_number").
		Joins("INNER JOIN classes ON classes.id = students.class_id").
		Scan(&result)
	return &result, nil
}
