package pgtype_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgtype/testutil"
)

func TestBoxCodec(t *testing.T) {
	testPgxCodec(t, "box", []PgxTranscodeTestCase{
		{
			pgtype.Box{
				P:     [2]pgtype.Vec2{{7.1, 5.2345678}, {3.14, 1.678}},
				Valid: true,
			},
			new(pgtype.Box),
			isExpectedEq(pgtype.Box{
				P:     [2]pgtype.Vec2{{7.1, 5.2345678}, {3.14, 1.678}},
				Valid: true,
			}),
		},
		{
			pgtype.Box{
				P:     [2]pgtype.Vec2{{7.1, 5.2345678}, {-13.14, -5.234}},
				Valid: true,
			},
			new(pgtype.Box),
			isExpectedEq(pgtype.Box{
				P:     [2]pgtype.Vec2{{7.1, 5.2345678}, {-13.14, -5.234}},
				Valid: true,
			}),
		},
		{nil, new(pgtype.Box), isExpectedEq(pgtype.Box{})},
	})
}

func TestBoxNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL: "select '3.14, 1.678, 7.1, 5.234'::box",
			Value: &pgtype.Box{
				P:     [2]pgtype.Vec2{{7.1, 5.234}, {3.14, 1.678}},
				Valid: true,
			},
		},
	})
}