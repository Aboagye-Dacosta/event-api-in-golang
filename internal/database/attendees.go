package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeesModel struct {
	DB *sql.DB
}

type Attendees struct {
	Id      int `json:"id"`
	UserId  int `json:"userId"`
	EventId int `json:"eventId"`
}

func (m *AttendeesModel) Insert(attendee *Attendees) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO attendees (user_id,event_id) VALUES ($1,$2) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, attendee.UserId, attendee.EventId).Scan(&attendee.Id)
}

func (m *AttendeesModel) GetAllAttendeesByEventID(eventID int) ([]*Attendees, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM attendees WHERE event_id = $1"

	rows, err := m.DB.QueryContext(ctx, query, eventID)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	defer rows.Close()

	attendees := []*Attendees{}

	for rows.Next() {
		var attendee Attendees

		err := rows.Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)
		if err != nil {
			return nil, err
		}

		attendees = append(attendees, &attendee)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attendees, nil
}

func (m *AttendeesModel) GetAttendeeByUserIDAndEventID(userID, eventID int) (*Attendees, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM attendees WHERE user_id = $1 AND event_id = $2"

	row := m.DB.QueryRowContext(ctx, query, userID, eventID)

	if err := row.Err(); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var attendee Attendees

	err := row.Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)

	if err != nil {
		return nil, err
	}

	return &attendee, nil
}

func (m *AttendeesModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM attendees WHERE id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}
