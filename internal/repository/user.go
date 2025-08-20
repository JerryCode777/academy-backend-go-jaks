package repository

import (
	"errors"
	"backend-academi/internal/models"
	"gorm.io/gorm"
)

// UserRepository maneja todas las operaciones de base de datos para usuarios
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserta un nuevo usuario en la base de datos
// Retorna error si el email ya existe (unique constraint)
func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByEmail busca un usuario por su email
// Usado principalmente para login
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	// GORM usa prepared statements - protegido contra SQL injection
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByID busca un usuario por su ID
// Usado por middleware para validar tokens JWT
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Exists verifica si un email ya está registrado en el sistema
// Usado antes del registro para dar mejor feedback al usuario
func (r *UserRepository) Exists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update actualiza los datos de un usuario existente
func (r *UserRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	return result.Error
}

// Delete elimina un usuario por ID (soft delete si está configurado en el modelo)
func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}