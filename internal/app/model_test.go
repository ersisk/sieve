package app

import (
	"sort"
	"testing"
	"time"

	"github.com/ersanisk/sieve/pkg/logentry"
)

func TestToggleSort(t *testing.T) {
	model := Model{
		sortOrder: SortAsc,
	}

	if model.sortOrder != SortAsc {
		t.Errorf("Expected initial sort order to be ascending, got %v", model.sortOrder)
	}

	model.toggleSort()
	if model.sortOrder != SortDesc {
		t.Errorf("Expected sort order to be descending after toggle, got %v", model.sortOrder)
	}

	model.toggleSort()
	if model.sortOrder != SortAsc {
		t.Errorf("Expected sort order to be ascending after second toggle, got %v", model.sortOrder)
	}
}

func TestSortLogic(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Hour)

	entries := []logentry.Entry{
		{Timestamp: later, Message: "second"},
		{Timestamp: now, Message: "first"},
		{Timestamp: now.Add(time.Minute * 30), Message: "middle"},
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	if entries[0].Message != "first" {
		t.Errorf("Expected first message 'first' when sorted ascending, got '%s'", entries[0].Message)
	}
	if entries[1].Message != "middle" {
		t.Errorf("Expected second message 'middle' when sorted ascending, got '%s'", entries[1].Message)
	}
	if entries[2].Message != "second" {
		t.Errorf("Expected third message 'second' when sorted ascending, got '%s'", entries[2].Message)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})

	if entries[0].Message != "second" {
		t.Errorf("Expected first message 'second' when sorted descending, got '%s'", entries[0].Message)
	}
	if entries[1].Message != "middle" {
		t.Errorf("Expected second message 'middle' when sorted descending, got '%s'", entries[1].Message)
	}
	if entries[2].Message != "first" {
		t.Errorf("Expected third message 'first' when sorted descending, got '%s'", entries[2].Message)
	}
}
