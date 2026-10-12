package compiler

func defaultCheckersForWorkload(parallelism, checkedFiles int) int {
	byMachine := max(4, min(8, parallelism))
	return max(4, min(byMachine, checkedFiles/32))
}
