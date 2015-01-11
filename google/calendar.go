package google

import (
	"fmt"
	calendar "google.golang.org/api/calendar/v3"
	"net/http"
	"time"
)

const (
	TZ_TAIWAN = "Asia/Taipei"

	SCOPE_CALENDAR = calendar.CalendarScope
	YYYYMMDD       = "2006-01-02"
)

type Calendar struct {
	Id       string
	Timezone *time.Location
	svc      *calendar.Service
}
type Event struct {
	IsAllDay bool
	Start    time.Time
	End      time.Time
}

func NewCalendar(oauthClient *http.Client, id, tz string) (*Calendar, error) {
	svc, err := calendar.New(oauthClient)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	return &Calendar{id, loc, svc}, nil
}

//NextFewDaysFristEvent return datetime array with the filtered event summary.
//It use the first match event, and ignore others same day event
//It always return same length array with requested `days`, even if error occurs
func (c *Calendar) NextFewDaysFristEvent(filterName string, days int) ([]Event, error) {
	var events = make([]Event, days)

	now := time.Now().In(c.Timezone)
	today := now.Format(time.RFC3339)
	lastDay := now.AddDate(0, 0, days).Format(time.RFC3339)

	res, err := c.svc.Events.List(c.Id).Q(filterName).Fields("items(start,end)").
		TimeMin(today).TimeMax(lastDay).SingleEvents(true).OrderBy("startTime").
		TimeZone(c.Timezone.String()).Do()
	if err != nil {
		return events, err
	}

	events = getDaysFirstEvent(res.Items, now, days)

	return events, nil
}
func getDaysFirstEvent(evts []*calendar.Event, current time.Time, withinDays int) []Event {
	events := make([]Event, withinDays)
	thDay := 0

	for _, ev := range evts {
		var evs, eve time.Time
		var isAllDay bool

		if isAllDay = ev.Start.Date != ""; isAllDay {
			evs, _ = time.Parse(YYYYMMDD, ev.Start.Date)
			eve, _ = time.Parse(YYYYMMDD, ev.End.Date)
		} else {
			evs, _ = time.Parse(time.RFC3339, ev.Start.DateTime)
			eve, _ = time.Parse(time.RFC3339, ev.End.DateTime)
		}

		var diff int
		diff = daysAgo(evs, current)
		if diff >= withinDays {
			break
		}

		if diff == thDay {
			events[thDay] = MakeEvent(isAllDay, evs, eve)
			thDay++
		} else if diff > thDay {
			events[diff] = MakeEvent(isAllDay, evs, eve)
			thDay = diff + 1
		} else {
			continue
		}
	}

	return events
}

// daysAgo return elapsed days between times provided. Less than 24 hours will return zero
// Both time should be in same timezone, to avoid confused or unexpected result
func daysAgo(src, dest time.Time) int {
	dt1 := time.Date(src.Year(), src.Month(), src.Day(), 0, 0, 0, 0, src.Location())
	dt2 := time.Date(dest.Year(), dest.Month(), dest.Day(), 0, 0, 0, 0, dest.Location())

	diff := dt1.Sub(dt2).Hours() / 24
	return int(diff)
}
func MakeEvent(allDay bool, start, end time.Time) Event {
	return Event{allDay, start, end}
}

func CreateCalendar(oauthClient *http.Client, name, tz string) (string, error) {
	svc, err := calendar.New(oauthClient)
	if err != nil {
		return "", fmt.Errorf("Unable to create Calendar service: %v", err)
	}

	cal := &calendar.Calendar{Summary: name, TimeZone: tz}
	calRes, err := svc.Calendars.Insert(cal).Do()
	if err != nil {
		return "", fmt.Errorf("Unable to insert calendar: %v", err)
	}

	return calRes.Id, nil
}
func (c *Calendar) AddCalendarUserACL(role, scopeType, scopeValue string) error {
	acl := &calendar.AclRule{Role: role,
		Scope: &calendar.AclRuleScope{scopeType, scopeValue}}

	if _, err := c.svc.Acl.Insert(c.Id, acl).Do(); err != nil {
		return err
	}
	return nil
}
func (c *Calendar) DeleteCalendarUserACL(writer string) error {
	rule := fmt.Sprintf("user:%s", writer)
	return c.svc.Acl.Delete(c.Id, rule).Do()
}
func (c *Calendar) PutACL(owner, googleGroup string, writer []string) error {
	if owner != "" {
		if err := c.AddCalendarUserACL("owner", "user", owner); err != nil {
			return err
		}
	}
	if googleGroup != "" {
		if err := c.AddCalendarUserACL("writer", "group", googleGroup); err != nil {
			return err
		}
	}
	for _, wt := range writer {
		if wt == "" {
			continue
		}
		if err := c.AddCalendarUserACL("writer", "user", wt); err != nil {
			return err
		}
	}

	return nil
}

func (c *Calendar) UpdateInfo(summary, description, location string) (*calendar.Calendar, error) {
	cal := &calendar.Calendar{}
	if summary != "" {
		cal.Summary = summary
	}
	if description != "" {
		cal.Description = description
	}
	if location != "" {
		cal.Location = location
	}

	return c.svc.Calendars.Patch(c.Id, cal).Do()
}
