package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// StringToPGText converts a regular string to a pgtype.Text, which can handle NULL values in PostgreSQL.
func StringToPGText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}

// StringToPGUUID converts a regular string to a pgtype.UUID, which can handle NULL values in PostgreSQL.
func StringToPGUUID(s string) pgtype.UUID {
	return pgtype.UUID{
		Bytes: [16]byte([]byte(s)),
		Valid: len(s) == 36, // UUIDs are typically 36 characters long (including hyphens)
	}
}

func TimeToPGTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: !t.IsZero(),
	}
}
