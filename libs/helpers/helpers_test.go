package helpers

import (
	"errors"
	"reflect"
	"sort"
	"testing"
)

var testsLatestVersion = []struct {
	list    []string
	wantStr string
	wantErr error
}{
	{
		list:    []string{"0.1", "0.10", "1.11", "0.5", "0.11", "1.0", "1.1"},
		wantStr: "1.11",
		wantErr: nil,
	},
	{
		list:    []string{"0.1.0", "0.10.0", "0.10.1", "0.5.0", "0.1.1", "0.10.0", "0.10.0"},
		wantStr: "0.10.1",
		wantErr: nil,
	},
	{
		list:    []string{"1.1.0.0", "1.1.0.1", "0.1.10.1", "1.0.5.0", "1.0.5.11", "1.0.10.0", "1.0.10.10"},
		wantStr: "1.1.0.1",
		wantErr: nil,
	},
	{
		list:    []string{"0.1", "0.10", "1.0", "1.2.3"},
		wantStr: "",
		wantErr: errors.New("Unable to determine latest version, format mismatch"),
	},
	{
		list:    []string{"1", "2", "3", "4", "5"},
		wantStr: "",
		wantErr: errors.New("Unable to determine latest version, format mismatch"),
	},
	{
		list:    []string{},
		wantStr: "",
		wantErr: errors.New("Unable to determine latest version, empty slice"),
	},
}

func TestLatestVersion(t *testing.T) {
	for _, tc := range testsLatestVersion {
		gotStr, gotErr := LatestVersion(tc.list)
		if !reflect.DeepEqual(gotErr, tc.wantErr) {
			t.Fatalf("LatestVersion(%v) returned error: '%v'; want: '%v'", tc.list, gotErr, tc.wantErr)
		}
		if gotStr != tc.wantStr {
			t.Fatalf("LatestVersion(%v) = %v, want %v.", tc.list, gotStr, tc.wantStr)
		}
	}
}

func BenchmarkLatestVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range testsLatestVersion {
			if tc.wantErr == nil {
				LatestVersion(tc.list)
			}
		}
	}
}

var testsStringInSlice = []struct {
	str  string
	list []string
	want bool
}{
	{
		str:  "a",
		list: []string{"a", "b", "a", "c", "c"},
		want: true,
	},
	{
		str:  "x",
		list: []string{"a", "b", "a", "c", "c"},
		want: false,
	},
	{
		str:  "Sunday",
		list: []string{"Monday", "Tuesday", "Thursday", "Friday", "Sunday"},
		want: true,
	},
	{
		str:  "sunday",
		list: []string{"Monday", "Tuesday", "Thursday", "Friday", "Sunday"},
		want: false,
	},
}

func TestStringInSlice(t *testing.T) {
	for _, tc := range testsStringInSlice {
		got := StringInSlice(tc.str, tc.list)
		if got != tc.want {
			t.Fatalf("StringInSlice(%v, %v) = %v, want %v.", tc.str, tc.list, got, tc.want)
		}
	}
}

func BenchmarkStringInSlice(b *testing.B) {
	// bench combined time to run through all test cases
	for i := 0; i < b.N; i++ {
		for _, tc := range testsStringInSlice {
			StringInSlice(tc.str, tc.list)
		}
	}
}

var testsUniqueSlice = []struct {
	test, want []string
}{
	{
		test: []string{"a", "b", "a", "c", "c"},
		want: []string{"a", "b", "c"},
	},
	{
		test: []string{"one", "one", "one", "one", "one"},
		want: []string{"one"},
	},
	{
		test: []string{"one", "two", "three", "one", "two", "three"},
		want: []string{"one", "two", "three"},
	},
	{
		test: []string{"one fish", "one fish", "two fish", "two fish", "three fish", "three fish", "four"},
		want: []string{"one fish", "two fish", "three fish", "four"},
	},
	{
		test: []string{" ", " ", " "},
		want: []string{" "},
	},
}

func TestUniqueSlice(t *testing.T) {
	for _, tc := range testsUniqueSlice {
		got := UniqueSlice(tc.test)
		sort.Strings(tc.want)
		sort.Strings(got)
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("UniqueSlice(%v) = %v, want %v.", tc.test, got, tc.want)
		}
	}
}

func BenchmarkUniqueSlice(b *testing.B) {
	// bench combined time to run through all test cases
	for i := 0; i < b.N; i++ {
		for _, tc := range testsUniqueSlice {
			UniqueSlice(tc.test)
		}
	}
}
