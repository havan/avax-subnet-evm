// (c) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rawdb

import (
	"bytes"

	"github.com/ava-labs/libevm/common"
	ethrawdb "github.com/ava-labs/libevm/core/rawdb"
	"github.com/ava-labs/libevm/ethdb"
)

// InspectDatabase traverses the entire database and checks the size
// of all different categories of data.
func InspectDatabase(db ethdb.Database, keyPrefix, keyStart []byte) error {
	type stat = ethrawdb.DatabaseStat
	stats := []struct {
		name      string
		keyLen    int
		keyPrefix []byte
		stat      *stat
	}{
		{"Trie segments", syncSegmentsKeyLength, syncSegmentsPrefix, &stat{}},
		{"Storage tries to fetch", syncStorageTriesKeyLength, syncStorageTriesPrefix, &stat{}},
		{"Code to fetch", codeToFetchKeyLength, CodeToFetchPrefix, &stat{}},
		{"Block numbers synced to", syncPerformedKeyLength, syncPerformedPrefix, &stat{}},
	}

	options := []ethrawdb.InspectDatabaseOption{
		ethrawdb.WithDatabaseMetadataKeys(func(key []byte) bool {
			return bytes.Equal(key, snapshotBlockHashKey) ||
				bytes.Equal(key, syncRootKey) ||
				(bytes.HasPrefix(key, upgradeConfigPrefix) && len(key) == len(upgradeConfigPrefix)+common.HashLength)
		}),
		ethrawdb.WithDatabaseStatRecorder(func(key []byte, size common.StorageSize) bool {
			for _, s := range stats {
				if len(key) == s.keyLen && bytes.HasPrefix(key, s.keyPrefix) {
					s.stat.Add(size)
					return true
				}
			}
			return false
		}),
		ethrawdb.WithDatabaseStatsTransformer(func(rows [][]string) [][]string {
			newRows := make([][]string, 0, len(rows))
			for _, row := range rows {
				switch db, cat := row[0], row[1]; {
				// Discard rows specific to libevm (geth) but irrelevant to coreth.
				case db == "Key-Value store" && (cat == "Difficulties" || cat == "Beacon sync headers"):
				case db == "Ancient store (Chain)":
				default:
					newRows = append(newRows, row)
				}
			}
			for _, s := range stats {
				newRows = append(newRows, []string{"State sync", s.name, s.stat.Size(), s.stat.Count()})
			}
			return newRows
		}),
	}

	return ethrawdb.InspectDatabase(db, keyPrefix, keyStart, options...)
}
