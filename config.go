package main

func getDBFilename() string {
	return "acbot.db"
}

func getDBTableName() string {
	return "highlight"
}

func getTurnipDBTableName() string {
	return "turnip"
}

// Price for turnip
type Price struct {
	buy  int
	sell [12]int
}

// MaxOfTurnip presents a reasonable price maximum
const MaxOfTurnip = 660

// Debug for detail api message
const Debug = false
