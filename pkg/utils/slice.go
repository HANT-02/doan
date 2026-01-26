package utils

func IsIntSliceContains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func IsInterfaceEmptySlice(value interface{}) bool {
	valueArray, ok := value.([]interface{})
	if !ok {
		return false
	}
	return len(valueArray) == 0
}

// FindDuplicates returns a slice of all values in a string slice that appear more than once.
// Each duplicate value appears at most once in the returned slice.
// The order of elements in the result is not guaranteed.
func FindDuplicates(ids []string) []string {
	seen := make(map[string]struct{}, len(ids))
	dups := make(map[string]struct{})
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			dups[id] = struct{}{}
			continue
		}
		seen[id] = struct{}{}
	}
	out := make([]string, 0, len(dups))
	for id := range dups {
		out = append(out, id)
	}
	return out
}

// DiffAB computes the asymmetric differences between two string slices A and B.
// X: elements that are present in B but not in A.
// Y: elements that are present in A but not in B.
// Each element appears at most once in X or Y regardless of how many
func DiffAB(A, B []string) (X, Y []string) {
	setA := make(map[string]struct{}, len(A))
	for _, a := range A {
		setA[a] = struct{}{}
	}

	// Tối đa X có thể bằng len(B)
	X = make([]string, 0, len(B))
	for _, b := range B {
		if _, ok := setA[b]; ok {
			// b có trong A → giao nhau; xoá để còn lại chỉ-a
			delete(setA, b)
		} else {
			// b không có trong A → thuộc X
			X = append(X, b)
		}
	}

	// Những gì còn lại trong setA là chỉ có ở A → Y
	Y = make([]string, 0, len(setA))
	for a := range setA {
		Y = append(Y, a)
	}
	return
}

// Intersect returns the common elements between two string slices A and B.
// Each element appears at most once in the result. If there is no intersection, it returns nil.
func Intersect(A, B []string) []string {
	if len(A) == 0 || len(B) == 0 {
		return nil
	}
	// Xây map từ mảng nhỏ hơn để tiết kiệm bộ nhớ
	if len(A) > len(B) {
		A, B = B, A
	}
	set := make(map[string]struct{}, len(A))
	for _, v := range A {
		set[v] = struct{}{}
	}
	out := make([]string, 0, len(A))
	for _, v := range B {
		if _, ok := set[v]; ok {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// UniqueStrings takes a slice of strings that may contain duplicates
// and returns a new slice containing only unique values,
// preserving the order of their first appearance.
func UniqueStrings(input []string) []string {
	if len(input) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))

	for _, v := range input {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}

	return out
}
