package helper

import "cache/internal/model"

//	func BinarySearch(a []model.KeyTTL, search int) (result int, searchCount int) {
//		mid := len(a) / 2
//		switch {
//		case len(a) == 0:
//			result = -1 // not found
//		case int(a[mid].TTL.UnixMilli()) > search:
//			result, searchCount = BinarySearch(a[:mid], search)
//		case int(a[mid].TTL.UnixMilli()) < search:
//			result, searchCount = BinarySearch(a[mid+1:], search)
//			if result >= 0 { // if anything but the -1 "not found" result
//				result += mid + 1
//			}
//		default: // a[mid] == search
//			result = mid // found
//		}
//		searchCount++
//		return
//	}
func BinarySearch(a []model.KeyTTL, x int64) int {
	r := -1 // not found
	start := 0
	end := len(a) - 1
	for start <= end {
		mid := (start + end) / 2
		if a[mid].TTL.UnixMilli() == x {
			r = mid // found
			break
		} else if a[mid].TTL.UnixMilli() < x {
			start = mid + 1
		} else if a[mid].TTL.UnixMilli() > x {
			end = mid - 1
		}
	}
	return r
}
