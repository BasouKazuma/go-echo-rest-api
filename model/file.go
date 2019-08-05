package model

type File struct {
	Id			int64		`json:"id" db:"id"`
	UserId		int64		`json:"user_id" db:"user_id"`
	Name		string		`json:"name" db:"name"`
	Hash		string		`json:"hash" db:"hash"`
}
