package Models

import "database/sql"

type User struct {
	Id                int64
	Name              string `sql:"size:255"`
	BillingAddress    Address       // Embedded struct
	BillingAddressID  sql.NullInt64 // Embedded struct's foreign key
}
type Address struct {
	ID        int
	Address1  string
	Address2  string
}

