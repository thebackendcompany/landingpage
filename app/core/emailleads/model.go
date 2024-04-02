package emailleads

import "time"

type SheetsEmailLeads struct {
	Email       string
	Body        string
	ContactedAt time.Time
}
