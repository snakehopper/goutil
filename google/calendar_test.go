package google

import (
	calendar "google.golang.org/api/calendar/v3"
	"testing"
	"time"
)

var (
	now, _ = time.Parse(time.RFC3339, "2014-12-28T00:00:00Z")
	edts   = []*calendar.EventDateTime{
		{DateTime: "2014-12-30T10:00:00Z"},
		{DateTime: "2014-12-30T18:00:00Z"},
		{Date: "2014-12-31"},
		{DateTime: "2014-12-31T02:00:00Z"},
		{DateTime: "2015-01-01T02:00:00Z"},
		{DateTime: "2015-01-02T02:00:00Z"},
		{DateTime: "2015-01-04T02:00:00Z"},
		{DateTime: "2015-01-08T02:00:00Z"},
		{DateTime: "2015-01-08T23:00:00Z"},
		{DateTime: "2015-01-09T00:00:00Z"},
	}
)

func setup_events() []*calendar.Event {

	events := make([]*calendar.Event, len(edts))
	for i, edt := range edts {
		events[i] = new(calendar.Event)
		events[i].Start = edt
		events[i].End = edt
	}
	return events
}
func TestGetDaysFirstEvent_todayEvent(t *testing.T) {
	evts := setup_events()

	today, _ := time.Parse(time.RFC3339, "2014-12-30T00:30:00Z")
	size := 3
	results := getDaysFirstEvent(evts, today, size)
	if len(results) != size {
		t.Errorf("getDaysFirstEvent length return %v, want %v", len(results), size)
	}

	want, _ := time.Parse(time.RFC3339, edts[0].DateTime)
	actual := results[0]
	if !actual.Start.Equal(want) {
		t.Errorf("today event return %v, want %v", actual, want)
	}
}
func TestGetDaysFirstEvent_noEventTomorrow(t *testing.T) {
	evts := setup_events()

	size := 3
	results := getDaysFirstEvent(evts, now, size)
	if len(results) != size {
		t.Errorf("getDaysFirstEvent length return %v, want %v", len(results), size)
	}
	if !results[0].Start.IsZero() || !results[1].Start.IsZero() {
		t.Errorf("the first 2 days event return %v %v, want zero",
			results[0].Start, results[1].Start)
	}

	r3th, _ := time.Parse(time.RFC3339, edts[0].DateTime)
	if !results[2].Start.Equal(r3th) {
		t.Errorf("3th days return %v, want %v", results[2].Start, r3th)
	}
}
func TestGetDaysFirstEvent_eventFullDay(t *testing.T) {
	evts := setup_events()

	size := 4
	results := getDaysFirstEvent(evts, now, size)
	if len(results) != size {
		t.Errorf("getDaysFirstEvent length return %v, want %v", len(results), size)
	}

	if !results[3].IsAllDay {
		t.Errorf("Event should be All Day event")
	}

	r, _ := time.Parse("2006-01-02", "2014-12-31")
	if !results[3].Start.Equal(r) {
		t.Errorf("day return %v, want %v", results[3].Start, r)
	}
}
func TestGetDaysFirstEvent_eventEmptyInMid(t *testing.T) {
	evts := setup_events()

	size := 20
	results := getDaysFirstEvent(evts, now, size)
	if len(results) != size {
		t.Errorf("getDaysFirstEvent length return %v, want %v", len(results), size)
	}

	want, _ := time.Parse(time.RFC3339, edts[4].DateTime)
	actual := results[4].Start
	if !actual.Equal(want) {
		t.Errorf("day return %v, want %v", actual, want)
	}

	want, _ = time.Parse(time.RFC3339, edts[5].DateTime)
	actual = results[5].Start
	if !actual.Equal(want) {
		t.Errorf("day return %v, want %v", actual, want)
	}

	want, _ = time.Parse(time.RFC3339, edts[6].DateTime)
	actual = results[7].Start
	if !actual.Equal(want) {
		t.Errorf("day return %v, want %v", actual, want)
	}

	want, _ = time.Parse(time.RFC3339, edts[7].DateTime)
	actual = results[11].Start
	if !actual.Equal(want) {
		t.Errorf("day return %v, want %v", actual, want)
	}

	want, _ = time.Parse(time.RFC3339, edts[9].DateTime)
	actual = results[12].Start
	if !actual.Equal(want) {
		t.Errorf("day return %v, want %v", actual, want)
	}
}

func TestDaysAgo(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	now, _ := time.Parse(time.RFC3339, "2015-01-09T01:30:26+08:00")
	current := now.In(loc)

	t1 := "2015-01-09T10:00:00+08:00"
	evs, _ := time.Parse(time.RFC3339, t1)

	actual, want := daysAgo(evs, current), 0
	if actual != want {
		t.Errorf("daysAgo return %v, want %v", actual, want)
	}
}
func TestDaysAgo_differentZone(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2015-01-08T16:30:26Z")

	t1 := "2015-01-09T10:00:00+08:00"
	evs, _ := time.Parse(time.RFC3339, t1)

	actual, want := daysAgo(evs, now), 0
	if actual != want {
		t.Errorf("daysAgo return %v, want %v", actual, want)
	}
}
