package main


func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func getAvgFloat64(floatNumbers []float64) float64 {
	floatSum := float64(0)
	for _, floatNumber := range floatNumbers {
		floatSum += floatNumber
	}

	return floatSum / float64(len(floatNumbers))
}

func getAvgInt64(int64Numbers []int64) int64 {
	intSum  := int64(0)
	for _, intNumber := range int64Numbers {
		intSum += intNumber
	}

	return intSum / int64(len(int64Numbers))
}

func getGreaterThanStone(floatNumbers []float64, stone float64) int {
	count := 0
	for _, floatNumber := range floatNumbers {
		if floatNumber > stone {
			count++
		}
	}

	return count
}
