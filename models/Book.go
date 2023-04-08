package models

import "gorm.io/gorm"

// modelo o esquema donde se manejara la informacion a guardar
type Books struct {
	ID uint `gorm:"primary key; autoIncrement"  json:"id"`

	Author *string `json:"author"`

	Title *string `json:"title"`

	Publisher *string `json:"publisher"` // manejo de envio para formato json

}

// la funcion AutoMigrate migrara los datos a la base de datos y creara una table con el esquema que se expecifica en la struct
func MigrateBook(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
