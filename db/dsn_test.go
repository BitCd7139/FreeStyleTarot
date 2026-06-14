package db

import (
	"strings"
	"testing"
)

func TestNormalizeDatabaseURL_SessionPooler(t *testing.T) {
	raw := "database:/postgres.myref:secret%2Fword@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres"
	got, err := NormalizeDatabaseURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantPrefix := "postgresql://postgres.myref:secret%2Fword@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres"
	if got != wantPrefix+"?sslmode=require" {
		t.Fatalf("got %q, want %q", got, wantPrefix+"?sslmode=require")
	}
}

func TestNormalizeDatabaseURL_RepairsConcatenatedPoolerURL(t *testing.T) {
	raw := "postgresql://hidden:secret%2Fword@hidden:5432//postgres.myref:secret%2Fword@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres"
	got, err := NormalizeDatabaseURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "postgres.myref:secret%2Fword@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres") {
		t.Fatalf("got %q", got)
	}
}

func TestNormalizeDatabaseURL_RejectsTransactionPoolerPort(t *testing.T) {
	raw := "postgresql://postgres.myref:secret@aws-1-ap-northeast-1.pooler.supabase.com:6543/postgres?sslmode=require"
	_, err := NormalizeDatabaseURL(raw)
	if err == nil {
		t.Fatal("expected error for transaction pooler port")
	}
}

func TestNormalizeDatabaseURL_LocalPostgres(t *testing.T) {
	raw := "postgres://postgres:postgres@localhost:5432/freestyletarot?sslmode=disable"
	got, err := NormalizeDatabaseURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "postgresql://postgres:postgres@localhost:5432/freestyletarot?sslmode=disable" {
		t.Fatalf("got %q", got)
	}
}

func TestToMigrateURL(t *testing.T) {
	in := "postgresql://postgres.myref:secret@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres?sslmode=require"
	want := "pgx5://postgres.myref:secret@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres?sslmode=require"
	if got := ToMigrateURL(in); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
