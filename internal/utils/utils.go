package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

// CheckVersion checks whether the version satisfies the provided version
// constraints.
func CheckVersion(versionString string, versionConstraint string) (bool, error) {
	ver, err := version.NewVersion(versionString)
	if err != nil {
		return false, fmt.Errorf("Unable to parse version %q: %v", versionString, err)
	}

	constraint, err := version.NewConstraint(versionConstraint)
	if err != nil {
		return false, fmt.Errorf("Unable to parse version constraint %q: %v", versionConstraint, err)
	}

	return constraint.Check(ver), nil
}

// ContextTimeout converts context deadline into timeout (time until it runs
// out). If the context has no deadline, the default timeout is returned.
func ContextTimeout(ctx context.Context, def time.Duration) int {
	deadline, ok := ctx.Deadline()
	if ok {
		return int(time.Until(deadline).Seconds())
	}

	return int(def.Seconds())
}

// HasAnyPrefix checks whether a value has any of the prefixes.
func HasAnyPrefix(value string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(value, p) {
			return true
		}
	}

	return false
}

// ValueInSlice checks whether a value is present in the given slice.
func ValueInSlice[T comparable](value T, slice []T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}

// SortMapKeys returns map keys sorted alphabetically.
func SortMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

// DiffSlice compares two slices and returns removed and added elements.
// Note: Does not find differences for duplicate elements.
func DiffSlices[T comparable](sliceA []T, sliceB []T) ([]T, []T) {
	mapA := make(map[T]bool, len(sliceA))
	mapB := make(map[T]bool, len(sliceB))

	for _, k := range sliceA {
		mapA[k] = true
	}

	for _, k := range sliceB {
		mapB[k] = true
	}

	var removed []T
	var added []T

	// Find elements in listA but not in listB (removed elements)
	for k := range mapA {
		if mapB[k] {
			removed = append(removed, k)
		}
	}

	// Find elements in listB but not in listA (added elements)
	for k := range mapB {
		if mapA[k] {
			added = append(added, k)
		}
	}

	return removed, added
}

// ToPrettyJSON converts the given value into JSON string. If value cannot
// be marshaled into JSON, an empty string is returned.
func ToPrettyJSON(v any) string {
	bytes, _ := json.MarshalIndent(v, "", "    ")
	return string(bytes)
}

// IsDir returns true if the given path is a directory.
func IsDir(name string) bool {
	stat, err := os.Stat(name)
	if err != nil {
		return false
	}

	return stat.IsDir()
}
