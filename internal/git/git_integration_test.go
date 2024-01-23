package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"

	"github.com/growthbook/gb-find-code-refs/internal/gb"
	"github.com/growthbook/gb-find-code-refs/internal/log"
	"github.com/growthbook/gb-find-code-refs/search"
)

const (
	repoDir = "testdata/repo"
	flag1   = "flag1"
	flag2   = "flag2"
	flag3   = "flag3"
)

func TestMain(m *testing.M) {
	log.Init(true)
	os.Exit(m.Run())
}

func setupRepo(t *testing.T) *git.Repository {
	os.RemoveAll(repoDir)
	require.NoError(t, os.MkdirAll(repoDir, 0700))
	repo, err := git.PlainInit(repoDir, false)
	require.NoError(t, err)
	return repo
}

// TestFindExtinctions is an integration test against a real Git repository stored under the testdata directory.
func TestFindExtinctions(t *testing.T) {
	repo := setupRepo(t)

	// Create commit history
	flagFile, err := os.Create(filepath.Join(repoDir, "flag1.txt"))
	require.NoError(t, err)
	_, err = flagFile.WriteString(flag1)
	require.NoError(t, err)
	require.NoError(t, flagFile.Close())
	flagFile, err = os.Create(filepath.Join(repoDir, "flag2.txt"))
	require.NoError(t, err)
	_, err = flagFile.WriteString(flag2)
	require.NoError(t, err)
	require.NoError(t, flagFile.Close())

	wt, err := repo.Worktree()
	require.NoError(t, err)

	who := object.Signature{Name: "GrowthBook", Email: "dev@growthbook.com", When: time.Unix(100000000, 0)}

	wt.Add("flag1.txt")
	wt.Add("flag2.txt")
	_, err = wt.Commit("add flags", &git.CommitOptions{All: true, Committer: &who, Author: &who})
	require.NoError(t, err)

	// Test with a removed file
	err = os.Remove(filepath.Join(repoDir, "flag1.txt"))
	require.NoError(t, err)

	who.When = who.When.Add(time.Minute)
	message2 := "remove flag1"
	commit2, err := wt.Commit(message2, &git.CommitOptions{All: true, Committer: &who, Author: &who})
	require.NoError(t, err)

	// Test with an updated (truncated) file
	flagFile, err = os.Create(filepath.Join(repoDir, "flag2.txt"))
	require.NoError(t, err)
	require.NoError(t, flagFile.Close())

	who.When = who.When.Add(time.Minute)
	message3 := "remove flag2"
	commit3, err := wt.Commit("remove flag2", &git.CommitOptions{All: true, Committer: &who, Author: &who})
	require.NoError(t, err)

	c := Client{workspace: repoDir}
	missingFlags := []string{flag1, flag2}
	matcher := search.Matcher{
		Elements: []search.ElementMatcher{
			search.NewElementMatcher(``, ``, []string{flag1, flag2}, nil),
		},
	}

	extinctions := make([]gb.ExtinctionRep, 0)
	extinctionsByProject, err := c.FindExtinctions(missingFlags, matcher, 10)
	require.NoError(t, err)
	extinctions = append(extinctions, extinctionsByProject...)

	expected := []gb.ExtinctionRep{
		{
			Revision: commit3.String(),
			Message:  message3,
			Time:     who.When.Unix() * 1000,
			FlagKey:  flag2,
		},
		{
			Revision: commit2.String(),
			Message:  message2,
			Time:     who.When.Add(-time.Minute).Unix() * 1000,
			FlagKey:  flag1,
		},
	}

	require.Equal(t, expected, extinctions)
}
