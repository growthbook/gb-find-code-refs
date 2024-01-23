package gb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/growthbook/gb-find-code-refs/internal/log"
)

func TestMain(m *testing.M) {
	log.Init(true)
	os.Exit(m.Run())
}

func TestCountAll(t *testing.T) {
	flagKey := "testFlag"

	h := HunkRep{
		StartingLineNumber: 1,
		Lines:              "testtest",
		FlagKey:            flagKey,
		Aliases:            []string{},
	}
	b := BranchRep{
		Name:     "",
		Head:     "",
		SyncTime: 0,
		References: []ReferenceHunksRep{{
			Hunks: []HunkRep{h},
		}},
	}
	count := b.CountAll()
	want := make(map[string]int64)
	want[flagKey] = 1
	require.Equal(t, count, want)

}

func TestCountByProjectAndFlag(t *testing.T) {
	flagKey := "testFlag"
	notFoundKey := "notFoundFlag"
	notFoundKey2 := "notFoundFlag2"

	h := HunkRep{
		StartingLineNumber: 1,
		Lines:              "testtest",
		FlagKey:            flagKey,
		Aliases:            []string{},
	}
	notFound := HunkRep{
		StartingLineNumber: 1,
		Lines:              "testtest",
		FlagKey:            flagKey,
		Aliases:            []string{},
	}
	b := BranchRep{
		Name:     "",
		Head:     "",
		SyncTime: 0,
		References: []ReferenceHunksRep{{
			Hunks: []HunkRep{h, notFound},
		}},
	}
	elements := [][]string{{flagKey, notFoundKey, notFoundKey2}}
	count := b.CountByFlag(elements)
	want := make(map[string]int64)
	want[flagKey] = 2
	want[notFoundKey] = 0
	want[notFoundKey2] = 0
	require.Equal(t, count, want)

}
