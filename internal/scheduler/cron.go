package scheduler

import "fmt"

func CronExpression(s Schedule) string {
	switch s.Frequency {
	case Daily:
		return fmt.Sprintf(
			"%d %d * * *",
			s.Minute,
			s.Hour,
		)

	case Weekly:
		return fmt.Sprintf(
			"%d %d * * %d",
			s.Minute,
			s.Hour,
			s.Weekday,
		)

	case Monthly:
		return fmt.Sprintf(
			"%d %d %d * * *",
			s.Minute,
			s.Hour,
			s.Day,
		)

	default:
		return fmt.Sprintf(
			"%d %d * * *",
			s.Minute,
			s.Hour,
		)

	}
}
