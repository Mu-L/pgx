package stdlib_test

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func getSelectRowsCounts(b *testing.B) []int64 {
	var rowCounts []int64
	{
		s := os.Getenv("PGX_BENCH_SELECT_ROWS_COUNTS")
		if s != "" {
			for _, p := range strings.Split(s, " ") {
				n, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					b.Fatalf("Bad PGX_BENCH_SELECT_ROWS_COUNTS value: %v", err)
				}
				rowCounts = append(rowCounts, n)
			}
		}
	}

	if len(rowCounts) == 0 {
		rowCounts = []int64{1, 10, 100, 1000}
	}

	return rowCounts
}

type BenchRowSimple struct {
	ID         int32
	FirstName  string
	LastName   string
	Sex        string
	BirthDate  time.Time
	Weight     int32
	Height     int32
	UpdateTime time.Time
}

func BenchmarkSelectRowsScanSimple(b *testing.B) {
	db := openDB(b)
	defer closeDB(b, db)

	rowCounts := getSelectRowsCounts(b)

	for _, rowCount := range rowCounts {
		b.Run(fmt.Sprintf("%d rows", rowCount), func(b *testing.B) {
			br := &BenchRowSimple{}
			for i := 0; i < b.N; i++ {
				rows, err := db.Query("select n, 'Adam', 'Smith ' || n, 'male', '1952-06-16'::date, 258, 72, '2001-01-28 01:02:03-05'::timestamptz from generate_series(1, $1) n", rowCount)
				if err != nil {
					b.Fatal(err)
				}

				for rows.Next() {
					rows.Scan(&br.ID, &br.FirstName, &br.LastName, &br.Sex, &br.BirthDate, &br.Weight, &br.Height, &br.UpdateTime)
				}

				if rows.Err() != nil {
					b.Fatal(rows.Err())
				}
			}
		})
	}
}

type BenchRowNull struct {
	ID         sql.NullInt32
	FirstName  sql.NullString
	LastName   sql.NullString
	Sex        sql.NullString
	BirthDate  sql.NullTime
	Weight     sql.NullInt32
	Height     sql.NullInt32
	UpdateTime sql.NullTime
}

func BenchmarkSelectRowsScanNull(b *testing.B) {
	db := openDB(b)
	defer closeDB(b, db)

	rowCounts := getSelectRowsCounts(b)

	for _, rowCount := range rowCounts {
		b.Run(fmt.Sprintf("%d rows", rowCount), func(b *testing.B) {
			br := &BenchRowSimple{}
			for i := 0; i < b.N; i++ {
				rows, err := db.Query("select n, 'Adam', 'Smith ' || n, 'male', '1952-06-16'::date, 258, 72, '2001-01-28 01:02:03-05'::timestamptz from generate_series(100000, 100000 +  $1) n", rowCount)
				if err != nil {
					b.Fatal(err)
				}

				for rows.Next() {
					rows.Scan(&br.ID, &br.FirstName, &br.LastName, &br.Sex, &br.BirthDate, &br.Weight, &br.Height, &br.UpdateTime)
				}

				if rows.Err() != nil {
					b.Fatal(rows.Err())
				}
			}
		})
	}
}

func BenchmarkFlatArrayEncodeArgument(b *testing.B) {
	db := openDB(b)
	defer closeDB(b, db)

	input := make(pgtype.FlatArray[string], 10)
	for i := range input {
		input[i] = fmt.Sprintf("String %d", i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var n int64
		err := db.QueryRow("select cardinality($1::text[])", input).Scan(&n)
		if err != nil {
			b.Fatal(err)
		}
		if n != int64(len(input)) {
			b.Fatalf("Expected %d, got %d", len(input), n)
		}
	}
}

func BenchmarkFlatArrayScanResult(b *testing.B) {
	db := openDB(b)
	defer closeDB(b, db)

	var input string
	for i := 0; i < 10; i++ {
		if i > 0 {
			input += ","
		}
		input += fmt.Sprintf(`'String %d'`, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result pgtype.FlatArray[string]
		err := db.QueryRow(fmt.Sprintf("select array[%s]::text[]", input)).Scan(&result)
		if err != nil {
			b.Fatal(err)
		}
		if len(result) != 10 {
			b.Fatalf("Expected %d, got %d", len(result), 10)
		}
	}
}
