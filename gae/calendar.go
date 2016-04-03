package gae

import (
	calendar "github.com/snakehopper/goutil/google"
	"golang.org/x/net/context"
)

func NewCalendar(c context.Context, id, tz string) (*calendar.Calendar, error) {
	oauth := OAuthClient(c, []string{calendar.SCOPE_CALENDAR})

	return calendar.NewCalendar(oauth, id, tz)
}

func CreateCalendar(c context.Context, name, tz string) (string, error) {
	oauth := OAuthClient(c, []string{calendar.SCOPE_CALENDAR})

	return calendar.CreateCalendar(oauth, name, tz)
}
