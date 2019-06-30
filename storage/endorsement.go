package storage

// Endorsement information for a specific level
type Endorsement struct {
	Delegate     string
	Level        int64
	Rights       []int64
	Endorsements []int64
	Misses       int64
	Block        string
	Cycle        int64
}

var endorsementStorage []Endorsement
var endorsementCycleMisses map[int64]int64

// RecordEndorsement in firestore
func RecordEndorsement(delegate string, level int64, rights []int64, endorsements []int64, block string) {
	e := Endorsement{
		Delegate:     delegate,
		Level:        level,
		Rights:       rights,
		Endorsements: endorsements,
		Misses:       int64(len(rights) - len(endorsements)),
		Block:        block,
		Cycle:        level / 4096,
	}
	endorsementStorage = append(endorsementStorage, e)
}

// GetLastRecordedEndorsementLevel so we can resume scanning from the returned level+1
func GetLastRecordedEndorsementLevel(delegate string) int64 {
	if len(endorsementStorage) == 0 {
		return -1
	}
	return endorsementStorage[len(endorsementStorage)-1].Level
}

// GetEndorsements returning the `count` most recent
func GetEndorsements(delegate string, count int) []Endorsement {
	if count > len(endorsementStorage) {
		count = len(endorsementStorage)
	}
	end := len(endorsementStorage)
	start := end - count
	return endorsementStorage[start:end]
}

// GetLatestEndorsementCycle that this delegate has reported from
func GetLatestEndorsementCycle(delegate string) int64 {
	if len(endorsementStorage) == 0 {
		return -1
	}
	return endorsementStorage[len(endorsementStorage)-1].Cycle
}

// GetCycleEndorsementMissCount ...
func GetCycleEndorsementMissCount(delegate string) int64 {
	latestCycle := GetLatestEndorsementCycle(delegate)
	return endorsementCycleMisses[latestCycle]
}
