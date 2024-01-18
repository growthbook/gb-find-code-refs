package gb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/launchdarkly/ld-find-code-refs/v2/internal/log"
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
		ProjKey:            "example",
		FlagKey:            flagKey,
		Aliases:            []string{},
	}
	b := BranchRep{
		Name:             "",
		Head:             "",
		UpdateSequenceId: nil,
		SyncTime:         0,
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

	projectKey := "exampleProject"
	h := HunkRep{
		StartingLineNumber: 1,
		Lines:              "testtest",
		ProjKey:            projectKey,
		FlagKey:            flagKey,
		Aliases:            []string{},
	}
	notFound := HunkRep{
		StartingLineNumber: 1,
		Lines:              "testtest",
		ProjKey:            "notfound",
		FlagKey:            flagKey,
		Aliases:            []string{},
	}
	b := BranchRep{
		Name:             "",
		Head:             "",
		UpdateSequenceId: nil,
		SyncTime:         0,
		References: []ReferenceHunksRep{{
			Hunks: []HunkRep{h, notFound},
		}},
	}
	projects := []string{"exampleProject"}
	elements := [][]string{{flagKey, notFoundKey, notFoundKey2}}
	count := b.CountByProjectAndFlag(elements, projects)
	want := make(map[string]map[string]int64)
	want[projectKey] = make(map[string]int64)
	want[projectKey][flagKey] = 1
	want[projectKey][notFoundKey] = 0
	want[projectKey][notFoundKey2] = 0
	require.Equal(t, count, want)

}
