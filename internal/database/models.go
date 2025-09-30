package database

import "database/sql"

type Models struct {
	Users     UserModel
	Attendees AttendeesModel
	Events    EventModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Users:     UserModel{DB: db},
		Attendees: AttendeesModel{DB: db},
		Events:    EventModel{DB: db},
	}
}
