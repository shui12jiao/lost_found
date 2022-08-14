// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: item.sql

package sqlc

import (
	"context"
	"time"
)

const addFound = `-- name: AddFound :one
/*
found
*/
INSERT INTO found (
    picker_openid,
    found_date,
    time_bucket,
    location_id,
    location_info,
    location_status,
    type_id,
    item_info,
    image,
    image_key,
    owner_info,
    addtional_info
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING id, create_at, picker_openid, found_date, time_bucket, location_id, location_info, location_status, type_id, item_info, image, image_key, owner_info, addtional_info
`

type AddFoundParams struct {
	PickerOpenid   string         `json:"pickerOpenid"`
	FoundDate      time.Time      `json:"foundDate"`
	TimeBucket     TimeBucket     `json:"timeBucket"`
	LocationID     int16          `json:"locationID"`
	LocationInfo   string         `json:"locationInfo"`
	LocationStatus LocationStatus `json:"locationStatus"`
	TypeID         int16          `json:"typeID"`
	ItemInfo       string         `json:"itemInfo"`
	Image          []byte         `json:"image"`
	ImageKey       string         `json:"imageKey"`
	OwnerInfo      string         `json:"ownerInfo"`
	AddtionalInfo  string         `json:"addtionalInfo"`
}

func (q *Queries) AddFound(ctx context.Context, arg AddFoundParams) (Found, error) {
	row := q.db.QueryRowContext(ctx, addFound,
		arg.PickerOpenid,
		arg.FoundDate,
		arg.TimeBucket,
		arg.LocationID,
		arg.LocationInfo,
		arg.LocationStatus,
		arg.TypeID,
		arg.ItemInfo,
		arg.Image,
		arg.ImageKey,
		arg.OwnerInfo,
		arg.AddtionalInfo,
	)
	var i Found
	err := row.Scan(
		&i.ID,
		&i.CreateAt,
		&i.PickerOpenid,
		&i.FoundDate,
		&i.TimeBucket,
		&i.LocationID,
		&i.LocationInfo,
		&i.LocationStatus,
		&i.TypeID,
		&i.ItemInfo,
		&i.Image,
		&i.ImageKey,
		&i.OwnerInfo,
		&i.AddtionalInfo,
	)
	return i, err
}

const addLost = `-- name: AddLost :one
/*
lost
*/
INSERT INTO lost (
    owner_openid,
    lost_date,
    time_bucket,
    type_id,
    item_info,
    image,
    image_key,
    location_id,
    location_id1,
    location_id2
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id, create_at, owner_openid, lost_date, time_bucket, type_id, item_info, image, image_key, location_id, location_id1, location_id2
`

type AddLostParams struct {
	OwnerOpenid string     `json:"ownerOpenid"`
	LostDate    time.Time  `json:"lostDate"`
	TimeBucket  TimeBucket `json:"timeBucket"`
	TypeID      int16      `json:"typeID"`
	ItemInfo    string     `json:"itemInfo"`
	Image       []byte     `json:"image"`
	ImageKey    string     `json:"imageKey"`
	LocationID  int16      `json:"locationID"`
	LocationId1 int16      `json:"locationId1"`
	LocationId2 int16      `json:"locationId2"`
}

func (q *Queries) AddLost(ctx context.Context, arg AddLostParams) (Lost, error) {
	row := q.db.QueryRowContext(ctx, addLost,
		arg.OwnerOpenid,
		arg.LostDate,
		arg.TimeBucket,
		arg.TypeID,
		arg.ItemInfo,
		arg.Image,
		arg.ImageKey,
		arg.LocationID,
		arg.LocationId1,
		arg.LocationId2,
	)
	var i Lost
	err := row.Scan(
		&i.ID,
		&i.CreateAt,
		&i.OwnerOpenid,
		&i.LostDate,
		&i.TimeBucket,
		&i.TypeID,
		&i.ItemInfo,
		&i.Image,
		&i.ImageKey,
		&i.LocationID,
		&i.LocationId1,
		&i.LocationId2,
	)
	return i, err
}

const addMatch = `-- name: AddMatch :one
/*
match
*/
INSERT INTO match (
    picker_openid,
    owner_openid,
    found_date,
    lost_date,
    type_id,
    item_info,
    image,
    image_key,
    comment
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, create_at, picker_openid, owner_openid, found_date, lost_date, type_id, item_info, image, image_key, comment
`

type AddMatchParams struct {
	PickerOpenid string    `json:"pickerOpenid"`
	OwnerOpenid  string    `json:"ownerOpenid"`
	FoundDate    time.Time `json:"foundDate"`
	LostDate     time.Time `json:"lostDate"`
	TypeID       int16     `json:"typeID"`
	ItemInfo     string    `json:"itemInfo"`
	Image        []byte    `json:"image"`
	ImageKey     string    `json:"imageKey"`
	Comment      string    `json:"comment"`
}

func (q *Queries) AddMatch(ctx context.Context, arg AddMatchParams) (Match, error) {
	row := q.db.QueryRowContext(ctx, addMatch,
		arg.PickerOpenid,
		arg.OwnerOpenid,
		arg.FoundDate,
		arg.LostDate,
		arg.TypeID,
		arg.ItemInfo,
		arg.Image,
		arg.ImageKey,
		arg.Comment,
	)
	var i Match
	err := row.Scan(
		&i.ID,
		&i.CreateAt,
		&i.PickerOpenid,
		&i.OwnerOpenid,
		&i.FoundDate,
		&i.LostDate,
		&i.TypeID,
		&i.ItemInfo,
		&i.Image,
		&i.ImageKey,
		&i.Comment,
	)
	return i, err
}

const deleteFound = `-- name: DeleteFound :exec
DELETE FROM found
WHERE id = $1
`

func (q *Queries) DeleteFound(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteFound, id)
	return err
}

const deleteLost = `-- name: DeleteLost :exec
DELETE FROM lost
WHERE id = $1
`

func (q *Queries) DeleteLost(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteLost, id)
	return err
}

const deleteMatch = `-- name: DeleteMatch :exec
DELETE FROM match
WHERE id = $1
`

func (q *Queries) DeleteMatch(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteMatch, id)
	return err
}

const getFound = `-- name: GetFound :one
SELECT id, create_at, picker_openid, found_date, time_bucket, location_id, location_info, location_status, type_id, item_info, image, image_key, owner_info, addtional_info FROM found
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetFound(ctx context.Context, id int32) (Found, error) {
	row := q.db.QueryRowContext(ctx, getFound, id)
	var i Found
	err := row.Scan(
		&i.ID,
		&i.CreateAt,
		&i.PickerOpenid,
		&i.FoundDate,
		&i.TimeBucket,
		&i.LocationID,
		&i.LocationInfo,
		&i.LocationStatus,
		&i.TypeID,
		&i.ItemInfo,
		&i.Image,
		&i.ImageKey,
		&i.OwnerInfo,
		&i.AddtionalInfo,
	)
	return i, err
}

const getLost = `-- name: GetLost :one
SELECT id, create_at, owner_openid, lost_date, time_bucket, type_id, item_info, image, image_key, location_id, location_id1, location_id2 FROM lost
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetLost(ctx context.Context, id int32) (Lost, error) {
	row := q.db.QueryRowContext(ctx, getLost, id)
	var i Lost
	err := row.Scan(
		&i.ID,
		&i.CreateAt,
		&i.OwnerOpenid,
		&i.LostDate,
		&i.TimeBucket,
		&i.TypeID,
		&i.ItemInfo,
		&i.Image,
		&i.ImageKey,
		&i.LocationID,
		&i.LocationId1,
		&i.LocationId2,
	)
	return i, err
}

const getMatch = `-- name: GetMatch :one
SELECT id, create_at, picker_openid, owner_openid, found_date, lost_date, type_id, item_info, image, image_key, comment FROM match
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetMatch(ctx context.Context, id int32) (Match, error) {
	row := q.db.QueryRowContext(ctx, getMatch, id)
	var i Match
	err := row.Scan(
		&i.ID,
		&i.CreateAt,
		&i.PickerOpenid,
		&i.OwnerOpenid,
		&i.FoundDate,
		&i.LostDate,
		&i.TypeID,
		&i.ItemInfo,
		&i.Image,
		&i.ImageKey,
		&i.Comment,
	)
	return i, err
}

const listFound = `-- name: ListFound :many
SELECT id, create_at, picker_openid, found_date, time_bucket, location_id, location_info, location_status, type_id, item_info, image, image_key, owner_info, addtional_info FROM found
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListFoundParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListFound(ctx context.Context, arg ListFoundParams) ([]Found, error) {
	rows, err := q.db.QueryContext(ctx, listFound, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Found{}
	for rows.Next() {
		var i Found
		if err := rows.Scan(
			&i.ID,
			&i.CreateAt,
			&i.PickerOpenid,
			&i.FoundDate,
			&i.TimeBucket,
			&i.LocationID,
			&i.LocationInfo,
			&i.LocationStatus,
			&i.TypeID,
			&i.ItemInfo,
			&i.Image,
			&i.ImageKey,
			&i.OwnerInfo,
			&i.AddtionalInfo,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listFoundByPicker = `-- name: ListFoundByPicker :many
SELECT id, create_at, picker_openid, found_date, time_bucket, location_id, location_info, location_status, type_id, item_info, image, image_key, owner_info, addtional_info FROM found
WHERE picker_openid = $1
ORDER BY id
`

func (q *Queries) ListFoundByPicker(ctx context.Context, pickerOpenid string) ([]Found, error) {
	rows, err := q.db.QueryContext(ctx, listFoundByPicker, pickerOpenid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Found{}
	for rows.Next() {
		var i Found
		if err := rows.Scan(
			&i.ID,
			&i.CreateAt,
			&i.PickerOpenid,
			&i.FoundDate,
			&i.TimeBucket,
			&i.LocationID,
			&i.LocationInfo,
			&i.LocationStatus,
			&i.TypeID,
			&i.ItemInfo,
			&i.Image,
			&i.ImageKey,
			&i.OwnerInfo,
			&i.AddtionalInfo,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLost = `-- name: ListLost :many
SELECT id, create_at, owner_openid, lost_date, time_bucket, type_id, item_info, image, image_key, location_id, location_id1, location_id2 FROM lost
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListLostParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListLost(ctx context.Context, arg ListLostParams) ([]Lost, error) {
	rows, err := q.db.QueryContext(ctx, listLost, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Lost{}
	for rows.Next() {
		var i Lost
		if err := rows.Scan(
			&i.ID,
			&i.CreateAt,
			&i.OwnerOpenid,
			&i.LostDate,
			&i.TimeBucket,
			&i.TypeID,
			&i.ItemInfo,
			&i.Image,
			&i.ImageKey,
			&i.LocationID,
			&i.LocationId1,
			&i.LocationId2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLostByOwner = `-- name: ListLostByOwner :many
SELECT id, create_at, owner_openid, lost_date, time_bucket, type_id, item_info, image, image_key, location_id, location_id1, location_id2 FROM lost
WHERE owner_openid = $1
ORDER BY id
`

func (q *Queries) ListLostByOwner(ctx context.Context, ownerOpenid string) ([]Lost, error) {
	rows, err := q.db.QueryContext(ctx, listLostByOwner, ownerOpenid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Lost{}
	for rows.Next() {
		var i Lost
		if err := rows.Scan(
			&i.ID,
			&i.CreateAt,
			&i.OwnerOpenid,
			&i.LostDate,
			&i.TimeBucket,
			&i.TypeID,
			&i.ItemInfo,
			&i.Image,
			&i.ImageKey,
			&i.LocationID,
			&i.LocationId1,
			&i.LocationId2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMatch = `-- name: ListMatch :many
SELECT id, create_at, picker_openid, owner_openid, found_date, lost_date, type_id, item_info, image, image_key, comment FROM match
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListMatchParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListMatch(ctx context.Context, arg ListMatchParams) ([]Match, error) {
	rows, err := q.db.QueryContext(ctx, listMatch, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Match{}
	for rows.Next() {
		var i Match
		if err := rows.Scan(
			&i.ID,
			&i.CreateAt,
			&i.PickerOpenid,
			&i.OwnerOpenid,
			&i.FoundDate,
			&i.LostDate,
			&i.TypeID,
			&i.ItemInfo,
			&i.Image,
			&i.ImageKey,
			&i.Comment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMatchByOwner = `-- name: ListMatchByOwner :many
SELECT id, create_at, picker_openid, owner_openid, found_date, lost_date, type_id, item_info, image, image_key, comment FROM match
WHERE owner_openid = $1
ORDER BY id
`

func (q *Queries) ListMatchByOwner(ctx context.Context, ownerOpenid string) ([]Match, error) {
	rows, err := q.db.QueryContext(ctx, listMatchByOwner, ownerOpenid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Match{}
	for rows.Next() {
		var i Match
		if err := rows.Scan(
			&i.ID,
			&i.CreateAt,
			&i.PickerOpenid,
			&i.OwnerOpenid,
			&i.FoundDate,
			&i.LostDate,
			&i.TypeID,
			&i.ItemInfo,
			&i.Image,
			&i.ImageKey,
			&i.Comment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMatchByPicker = `-- name: ListMatchByPicker :many
SELECT id, create_at, picker_openid, owner_openid, found_date, lost_date, type_id, item_info, image, image_key, comment FROM match
WHERE picker_openid = $1
ORDER BY id
`

func (q *Queries) ListMatchByPicker(ctx context.Context, pickerOpenid string) ([]Match, error) {
	rows, err := q.db.QueryContext(ctx, listMatchByPicker, pickerOpenid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Match{}
	for rows.Next() {
		var i Match
		if err := rows.Scan(
			&i.ID,
			&i.CreateAt,
			&i.PickerOpenid,
			&i.OwnerOpenid,
			&i.FoundDate,
			&i.LostDate,
			&i.TypeID,
			&i.ItemInfo,
			&i.Image,
			&i.ImageKey,
			&i.Comment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
