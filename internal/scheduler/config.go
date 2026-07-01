package scheduler

type Frequency string

const (
	Daily   Frequency = "daily"
	Weekly  Frequency = "weekly"
	Monthly Frequency = "monthly"
)

type Schedule struct {
	Frequency Frequency

	Hour   int
	Minute int

	// For use in weekly: 0 := Sunday, 6 := Saturday
	Weekday int

	// For use for monthly
	Day int
}
