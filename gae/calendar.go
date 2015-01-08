package gae

import (
	"appengine"
	calendar "github.com/snakehopper/goutil/google"
)

func NewCalendar(c appengine.Context, id, tz string) (*calendar.Calendar, error) {
	oauth := OAuthClient(c, []string{calendar.SCOPE_CALENDAR})

	return calendar.NewCalendar(oauth, id, tz)
}

func CreateCalendar(c appengine.Context, name, tz string) (string, error) {
	oauth := OAuthClient(c, []string{calendar.SCOPE_CALENDAR})

	return calendar.CreateCalendar(oauth, name, tz)
}
