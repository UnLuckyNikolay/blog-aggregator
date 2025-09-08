package fetch

import (
	"database/sql"
	"time"
)

// toNullString parses RSS text into sql.NullString.
//
// Returns empty string and Valid=false if parsing fails.
func toNullString(text string) sql.NullString {
	if text == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: text, Valid: true}
}

// toNullTime parses RSS timestamp into sql.NullTime.
//
// Returns empty time.Time and Valid=false if parsing fails.
func toNullTime(timestamp string) sql.NullTime {
	parsedTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", timestamp)
	if err != nil {
		//fmt.Printf("Time Parse error: %s\n", err)
		return sql.NullTime{Time: time.Time{}, Valid: false}
	}
	//fmt.Printf("Time Parse: parsed %s as %s\n", timestamp, parsedTime)
	return sql.NullTime{Time: parsedTime, Valid: true}
}
