package storage

// Baking information about a delegate @ level
type Baking struct {
	Delegate       string
	Level          int64
	BakerPriority  int64
	DelegateRights int64
	DelegateMissed bool
	DelegateStole  bool
	DelegateBaked  bool
	BlockHash      string
	Cycle          int64
}

var bakingStorage []Baking
var bakingCycleMisses map[int64]int64

// RecordBaking in local storage
func RecordBaking(delegate string, level int64, delegateRights int64, bakerPriority int64, blockHash string) {
	b := Baking{
		Delegate:       delegate,
		Level:          level,
		BakerPriority:  bakerPriority,
		DelegateRights: delegateRights,
		DelegateMissed: delegateRights >= 0 && bakerPriority > delegateRights,
		DelegateStole:  bakerPriority == delegateRights && delegateRights > 0,
		DelegateBaked:  bakerPriority == delegateRights,
		BlockHash:      blockHash,
		Cycle:          level / 4096,
	}
	if b.DelegateMissed {
		bakingCycleMisses[b.Cycle]++
	}
	bakingStorage = append(bakingStorage, b)
}

// GetLastRecordedBakeLevel so we can resume scanning from the returned level+1
func GetLastRecordedBakeLevel(delegate string) int64 {
	if len(bakingStorage) == 0 {
		return -1
	}
	return bakingStorage[len(bakingStorage)-1].Level
}

// GetLatestBakingCycle that this delegate has reported from
func GetLatestBakingCycle(delegate string) int64 {
	if len(bakingStorage) == 0 {
		return -1
	}
	return bakingStorage[len(bakingStorage)-1].Cycle
}

// GetCycleBakeMissCount returns the number of misses in the current cycle
func GetCycleBakeMissCount(delegate string) (int64, int64) {
	latestCycle := GetLatestBakingCycle(delegate)
	return bakingCycleMisses[latestCycle], latestCycle
}
