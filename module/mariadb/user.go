package mariadb

//Field const
const (
	FieldID      = "id"
	FieldName    = "name"
	FieldAccount = "account"
	FieldAge     = "age"
)

type User struct {
	id      string
	Account string
	Name    string
	Age     uint
}

func (u *User) Create(m *Maria) error {
	db := m.db.Create(u)
	return db.Error
}
