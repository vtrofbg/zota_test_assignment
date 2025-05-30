package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	error_brake_chance = 0 // disable occasion faulties

	Log_init()
	DB_init()
	// clean up DB before running tests
	_, _ = db.Exec(`DELETE FROM key_value`)
	os.Exit(m.Run())
}

func TestPutAndGet(t *testing.T) {
	key, val := "test_key", "test_val"
	Put(key, val)

	got := Get(key)
	if got != val {
		t.Errorf("GET(%q) = %q expected %q", key, got, val)
	}
}

func TestGetNonExistentKey(t *testing.T) {
	got := Get("non_existent_key")
	if got != "" {
		t.Errorf("GET(non_existent_key) = %q, expected empty string", got)
	}
}

func TestDelete(t *testing.T) {
	key, val := "delete_key", "delete_val"
	Put(key, val)

	Delete(key)

	got := Get(key)
	if got != "" {
		t.Errorf("after DELETE, GET(%q) = %q expected empty string", key, got)
	}
}

func TestFaultyPut(t *testing.T) {
	// set error_brake_chance at 100% for test
	oldChance := error_brake_chance
	error_brake_chance = 1.0

	key, val := "faulty_put_key", "faulty_val"
	Put(key, val)

	// as faulty PUT not insert, GET should return empty or different value
	got := Get(key)
	log.Println(got)
	if got == val {
		t.Errorf("'faulty PUT issue, expected no value or different, got %q", got)
	}

	error_brake_chance = oldChance
}

func TestFaultyGet(t *testing.T) {
	// set error_brake_chance at 100% for test
	oldChance := error_brake_chance
	error_brake_chance = 1.0

	Put("faulty_key1", "val1")

	got := Get("faulty_key1")
	if got == "val1" {
		t.Errorf("Faulty Get: expected incorrect value, got correct %q", got)
	}

	error_brake_chance = oldChance
}

func TestFaultyDelete(t *testing.T) {
	// set error_brake_chance at 100% for test
	oldChance := error_brake_chance
	error_brake_chance = 1.0

	Put("del1", "v1")
	Put("del2", "v2")
	Delete("del1")

	//reset for normal GET results
	error_brake_chance = oldChance
	val1 := Get("del1")

	// key should be deleted, but not which desired initially
	if val1 == "v1" {
		t.Errorf("faulty DELETE, expected del1 to remain, but it was deleted")
	}
}

func TestSwapTwoRndValues(t *testing.T) {

	Put("mut_key1", "val1")
	Put("mut_key2", "val2")

	// call manual mutation
	swap_two_rnd_values()

	got1 := Get("mut_key1")
	got2 := Get("mut_key2")

	if got1 != "val2" || got2 != "val1" {
		t.Errorf("Error, expected mut_key1 = 'val2', mut_key2 = 'val1', BUT got mut_key1 = %q, mut_key2 = %q", got1, got2)
	}
}

func TestDump(t *testing.T) {
	// clean up once again DB
	_, _ = db.Exec(`DELETE FROM key_value`)

	Put("dump1", "val1")
	Put("dump2", "val2")
	Put("dump3", "val3")

	// Capture stdout
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Dump()

	// Restore stdout
	w.Close()
	os.Stdout = stdout
	buf.ReadFrom(r)

	output := buf.String()

	// split into lines
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// expect 1 header + 3 records(output rows) = 4
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines (header + records), got %d\nOutput:\n%s", len(lines), output)
	}
}
