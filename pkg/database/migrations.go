package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// DropColumnsManual elimina campos específicos de tablas
// IMPORTANTE: Solo ejecutar después de haber eliminado los campos de los structs
func DropColumnsManual(db *gorm.DB) error {
	log.Println("Ejecutando migraciones manuales para eliminar campos...")

	// Lista de campos a eliminar
	// Formato: "tabla.campo"
	fieldsToRemove := []struct {
		Table  string
		Column string
	}{
		// EJEMPLO: {"users", "old_field"},
		// EJEMPLO: {"courses", "deprecated_column"},

		// AGREGAR AQUÍ LOS CAMPOS QUE ELIMINES:
		// {"tabla", "campo_eliminado"},
		// {"students", "grade"}, // COMENTADO: ya fue eliminado
		// {"questionnaire_responses", "user_id"}, // COMENTADO: ya fue eliminado
	}

	// Ejecutar drops
	for _, field := range fieldsToRemove {
		if err := dropColumnIfExists(db, field.Table, field.Column); err != nil {
			return fmt.Errorf("error eliminando %s.%s: %w", field.Table, field.Column, err)
		}
	}

	log.Println("Migraciones manuales completadas")
	return nil
}

// dropColumnIfExists elimina una columna solo si existe
func dropColumnIfExists(db *gorm.DB, tableName, columnName string) error {
	// Verificar si la columna existe
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = ? AND column_name = ?
		)
	`

	if err := db.Raw(query, tableName, columnName).Scan(&exists).Error; err != nil {
		return err
	}

	if exists {
		log.Printf("Eliminando columna %s.%s...", tableName, columnName)

		// Usar Migrator para eliminar la columna
		if err := db.Migrator().DropColumn(tableName, columnName); err != nil {
			return err
		}

		log.Printf("Columna %s.%s eliminada", tableName, columnName)
	} else {
		log.Printf("Columna %s.%s no existe, saltando...", tableName, columnName)
	}

	return nil
}
