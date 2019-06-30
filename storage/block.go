package storage

// Block information to indicate that it's been analyzed
type Block struct {
	Level     int64
	BlockHash string
}

var blockStorage []Block

// RecordBlock in firestore
func RecordBlock(level int64, blockHash string) {
	b := Block{
		Level:     level,
		BlockHash: blockHash,
	}
	blockStorage = append(blockStorage, b)
}

// GetLastRecordedBlockLevel so we can resume scanning from the returned level+1
func GetLastRecordedBlockLevel() int64 {
	if len(blockStorage) == 0 {
		return -1
	}
	return blockStorage[len(blockStorage)-1].Level
}
