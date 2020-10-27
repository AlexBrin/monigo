package monigo

type History struct {
	ID        uint    `gorm:"primary_key"`
	Type      string  `gorm:"column:typ;type:VARCHAR(7)"`
	Val       float64 `gorm:"type:'DOUBLE(5, 2)'"`
	CreatedAt int64   `sql:"DEFAULT:(strftime('%s', 'now'))"`
}
