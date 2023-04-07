package moels 

Import "gorm.io/gorm"

// modelo o esquema donde se manejara la informacion a guardar o extraer 
type Books struct {

	ID  uint `gorm:"primary key; autoIncrement"  json:"id"` 

	Author *string `json:"author"`

	Title *string `json:"title"`

	Publisher *string `json:"publisher"`// manejo de envio para formato json

}
