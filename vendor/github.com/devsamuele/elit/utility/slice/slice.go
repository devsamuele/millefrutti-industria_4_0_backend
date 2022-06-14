package slice

func Find[T any](slice []T, compareFn func(element T) bool) (T, bool) {
	var t T
	for i := 0; i < len(slice); i++ {
		if compareFn(slice[i]) {
			return slice[i], true
		}
	}

	return t, false
}

func Filter[T any](sliceA []T, compareFn func(elementA T) bool) []T {
	result := make([]T, 0)
	for i := 0; i < len(sliceA); i++ {
		if compareFn(sliceA[i]) {
			result = append(result, sliceA[i])
		}
	}
	return result
}

// func Intersection[T any](sliceA []T, sliceB []T, compareFn func(elementA T, elementB T) bool) []T {
// 	intersection := make([]T, 0)
// 	var found bool
// 	for i := 0; i < len(sliceA); i++ {
// 		for j := 0; j < len(sliceB) && !found; j++ {
// 			if compareFn(sliceA[i], sliceB[j]) {
// 				intersection = append(intersection, sliceA[i])
// 				found = true
// 			}
// 		}
// 		found = false
// 	}
// 	return intersection
// }

func Map[T any, K any](slice []T, mapFn func(element T) K) []K {
	aux := make([]K, 0)
	for _, item := range slice {
		aux = append(aux, mapFn(item))
	}
	return aux
}
