package gb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/launchdarkly/ld-find-code-refs/v2/internal/validation"
)

type ConfigurationError struct {
	error
}

func newConfigurationError(e string) ConfigurationError {
	return ConfigurationError{errors.New((e))}
}

var (
	NotFoundErr                       = errors.New("not found")
	ConflictErr                       = errors.New("conflict")
	RateLimitExceededErr              = errors.New("rate limit exceeded")
	InternalServiceErr                = errors.New("internal service error")
	ServiceUnavailableErr             = errors.New("service unavailable")
	BranchUpdateSequenceIdConflictErr = errors.New("updateSequenceId conflict")
	RepositoryDisabledErr             = newConfigurationError("repository is disabled")
	UnauthorizedErr                   = newConfigurationError("unauthorized, check your LaunchDarkly access token")
	EntityTooLargeErr                 = newConfigurationError("entity too large")
)

// IsTransient returns true if the error returned by the LaunchDarkly API is either unexpected, or unable to be resolved by the user.
func IsTransient(err error) bool {
	var e ConfigurationError
	return !errors.As(err, &e)
}

type BranchRep struct {
	Name             string              `json:"name"`
	Head             string              `json:"head"`
	UpdateSequenceId *int                `json:"updateSequenceId,omitempty"`
	SyncTime         int64               `json:"syncTime"`
	References       []ReferenceHunksRep `json:"references,omitempty"`
	CommitTime       int64               `json:"commitTime,omitempty"`
}

func (b BranchRep) TotalHunkCount() int {
	count := 0
	for _, r := range b.References {
		count += len(r.Hunks)
	}
	return count
}

func (b BranchRep) WriteToJSON(outDir, projKey, sha string) (path string, err error) {
	// Try to create a filename with a shortened sha, but if the sha is too short for some unexpected reason, use the branch name instead
	var tag string
	if len(sha) >= 7 {
		tag = sha[:7]
	} else {
		tag = b.Name
	}

	absPath, err := validation.NormalizeAndValidatePath(outDir)
	if err != nil {
		return "", fmt.Errorf("invalid outDir '%s': %w", outDir, err)
	}

	// replace any forward slashes in filename
	filename := strings.ReplaceAll(fmt.Sprintf("coderefs_%s_%s.json", projKey, tag), "/", "_")
	path = filepath.Join(absPath, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	records := make([][]string, 0, len(b.References)+1)
	for _, ref := range b.References {
		records = append(records, ref.toRecords()...)
	}

	// sort records by flag key
	sort.Slice(records, func(i, j int) bool {
		// sort by flagKey -> path -> startingLineNumber
		for k := 0; k < 3; k++ {
			if records[i][k] != records[j][k] {
				return records[i][k] < records[j][k]
			}
		}
		// above loop should always return since startingLineNumber is guaranteed to be unique
		return false
	})

	r, err := json.Marshal(records)
	if err != nil {
		return "", err
	}

	_, err = f.Write(r)

	return path, err
}

type ReferenceHunksRep struct {
	Path  string    `json:"path"`
	Hunks []HunkRep `json:"hunks"`
}

func (r ReferenceHunksRep) toRecords() [][]string {
	ret := make([][]string, 0, len(r.Hunks))
	for _, hunk := range r.Hunks {
		ret = append(ret, []string{hunk.FlagKey, hunk.ProjKey, r.Path, strconv.FormatInt(int64(hunk.StartingLineNumber), 10), hunk.Lines, strings.Join(hunk.Aliases, " "), hunk.ContentHash})
	}
	return ret
}

type HunkRep struct {
	StartingLineNumber int      `json:"startingLineNumber"`
	Lines              string   `json:"lines,omitempty"`
	ProjKey            string   `json:"projKey"`
	FlagKey            string   `json:"flagKey"`
	Aliases            []string `json:"aliases,omitempty"`
	ContentHash        string   `json:"contentHash,omitempty"`
}

// Returns the number of lines overlapping between the receiver (h) and the parameter (hr) hunkreps
// The return value will be negative if the hunks do not overlap
func (h HunkRep) Overlap(hr HunkRep) int {
	return h.StartingLineNumber + h.NumLines() - hr.StartingLineNumber
}

func (h HunkRep) NumLines() int {
	return strings.Count(h.Lines, "\n") + 1
}

type ExtinctionRep struct {
	Revision string `json:"revision"`
	Message  string `json:"message"`
	Time     int64  `json:"time"`
	ProjKey  string `json:"projKey"`
	FlagKey  string `json:"flagKey"`
}

type tableData [][]string

func (t tableData) Len() int {
	return len(t)
}

func (t tableData) Less(i, j int) bool {
	first, _ := strconv.ParseInt(t[i][1], 10, 32)
	second, _ := strconv.ParseInt(t[j][1], 10, 32)
	return first > second
}

func (t tableData) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

const maxFlagKeysDisplayed = 50

func (b BranchRep) CountAll() map[string]int64 {
	refCount := map[string]int64{}
	for _, ref := range b.References {
		for _, hunk := range ref.Hunks {
			refCount[hunk.FlagKey]++
		}
	}
	return refCount
}

func (b BranchRep) CountByProjectAndFlag(matcher [][]string, projects []string) map[string]map[string]int64 {
	refCountByFlag := map[string]map[string]int64{}
	for i, project := range projects {
		refCountByFlag[project] = map[string]int64{}
		for _, flag := range matcher[i] {
			refCountByFlag[project][flag] = 0
		}
		for _, ref := range b.References {
			for _, hunk := range ref.Hunks {
				if hunk.ProjKey == project {
					refCountByFlag[project][hunk.FlagKey]++
				}
			}
		}
	}
	return refCountByFlag
}

func (b BranchRep) PrintReferenceCountTable() {
	data := tableData{}

	for k, v := range b.CountAll() {
		data = append(data, []string{k, strconv.FormatInt(v, 10)})
	}
	sort.Sort(data)

	truncatedData := data
	var additionalRefCount int64 = 0
	if len(truncatedData) > maxFlagKeysDisplayed {
		truncatedData = data[0:maxFlagKeysDisplayed]

		for _, v := range data[maxFlagKeysDisplayed:] {
			i, _ := strconv.ParseInt(v[1], 10, 64)
			additionalRefCount += i
		}
	}
	truncatedData = append(truncatedData, []string{"Other flags", strconv.FormatInt(additionalRefCount, 10)})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Flag", "# References"})
	table.SetBorder(false)
	table.AppendBulk(truncatedData)
	table.Render()
}
