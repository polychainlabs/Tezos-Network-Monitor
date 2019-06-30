package alert

// Alias the address with a friendly alias if one exists, or truncate
func Alias(aliases map[string]string, address string) string {
	if _, ok := aliases[address]; ok {
		return aliases[address]
	}
	return shortenPkh(address)
}

// Helper to truncate long Public key hashes
func shortenPkh(pkh string) string {
	end := len(pkh) - 1
	if end > 6 {
		end = 6
	}
	return pkh[0:end] + "..."
}
