// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package sqlc

import (
	"context"
)

type Querier interface {
	AddFound(ctx context.Context, arg AddFoundParams) (Found, error)
	AddLocationNarrow(ctx context.Context, arg AddLocationNarrowParams) (LocationNarrow, error)
	AddLocationWide(ctx context.Context, arg AddLocationWideParams) (LocationWide, error)
	AddLost(ctx context.Context, arg AddLostParams) (Lost, error)
	AddManager(ctx context.Context, arg AddManagerParams) (Manager, error)
	AddMatch(ctx context.Context, arg AddMatchParams) (Match, error)
	AddTypeNarrow(ctx context.Context, arg AddTypeNarrowParams) (TypeNarrow, error)
	AddTypeWide(ctx context.Context, name string) (TypeWide, error)
	AddUsr(ctx context.Context, arg AddUsrParams) (Usr, error)
	DeleteFound(ctx context.Context, id int32) error
	DeleteLocationNarrow(ctx context.Context, id int16) error
	DeleteLocationWide(ctx context.Context, id int16) error
	DeleteLost(ctx context.Context, id int32) error
	DeleteManager(ctx context.Context, id int16) error
	DeleteMatch(ctx context.Context, id int32) error
	DeleteTypeNarrow(ctx context.Context, id int16) error
	DeleteTypeWide(ctx context.Context, id int16) error
	GetFound(ctx context.Context, id int32) (Found, error)
	GetLocationNarrow(ctx context.Context, id int16) (LocationNarrow, error)
	GetLocationWide(ctx context.Context, id int16) (LocationWide, error)
	GetLost(ctx context.Context, id int32) (Lost, error)
	GetManager(ctx context.Context, id int16) (Manager, error)
	GetMatch(ctx context.Context, id int32) (Match, error)
	GetTypeNarrow(ctx context.Context, id int16) (TypeNarrow, error)
	GetTypeWide(ctx context.Context, id int16) (TypeWide, error)
	GetUsr(ctx context.Context, openid string) (Usr, error)
	ListFound(ctx context.Context, arg ListFoundParams) ([]Found, error)
	ListLocationNarrow(ctx context.Context) ([]LocationNarrow, error)
	ListLocationWide(ctx context.Context) ([]LocationWide, error)
	ListLost(ctx context.Context, arg ListLostParams) ([]Lost, error)
	ListManager(ctx context.Context) ([]Manager, error)
	ListMatch(ctx context.Context, arg ListMatchParams) ([]Match, error)
	ListTypeNarrow(ctx context.Context) ([]TypeNarrow, error)
	ListTypeWide(ctx context.Context) ([]TypeWide, error)
}

var _ Querier = (*Queries)(nil)
