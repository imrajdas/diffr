package diffr

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type jsonResult struct {
	Left  string     `json:"left"`
	Right string     `json:"right"`
	Stats jsonStats  `json:"stats"`
	Files []jsonFile `json:"files"`
}

type jsonStats struct {
	Changed   int    `json:"changed"`
	Added     int    `json:"added"`
	Removed   int    `json:"removed"`
	Identical int    `json:"identical"`
	Elapsed   string `json:"elapsed"`
}

type jsonFile struct {
	Path      string `json:"path"`
	Status    string `json:"status"`
	Kind      string `json:"kind"`
	Summary   string `json:"summary,omitempty"`
	LeftHash  string `json:"left_hash,omitempty"`
	RightHash string `json:"right_hash,omitempty"`
}

func (r *Result) JSON() jsonResult {
	out := jsonResult{
		Left:  r.LeftAbs,
		Right: r.RightAbs,
		Stats: jsonStats{
			Changed:   r.Stats.Changed,
			Added:     r.Stats.Added,
			Removed:   r.Stats.Removed,
			Identical: r.Stats.Identical,
			Elapsed:   r.Stats.Elapsed,
		},
		Files: make([]jsonFile, 0, len(r.Files)),
	}
	for _, f := range r.Files {
		out.Files = append(out.Files, jsonFile{
			Path:      f.RelPath,
			Status:    f.Status.String(),
			Kind:      f.Kind.String(),
			Summary:   f.Summary,
			LeftHash:  f.LeftHash,
			RightHash: f.RightHash,
		})
	}
	return out
}

func WriteCLI(res *Result, stdout, jsonOut bool, patchFile string) error {
	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(res.JSON()); err != nil {
			return err
		}
	} else if stdout {
		patch := res.Patch()
		if stdoutIsTTY() {
			patch = colorizePatch(patch)
		}
		if _, err := os.Stdout.WriteString(patch); err != nil {
			return err
		}
		if patch != "" && !strings.HasSuffix(patch, "\n") {
			if _, err := os.Stdout.WriteString("\n"); err != nil {
				return err
			}
		}
	}

	if patchFile != "" {
		patch := res.Patch()
		if patchFile == "-" {
			if !stdout && !jsonOut {
				if _, err := os.Stdout.WriteString(patch); err != nil {
					return err
				}
			}
		} else if err := os.WriteFile(patchFile, []byte(patch), 0o644); err != nil {
			return err
		} else {
			fmt.Fprintf(os.Stderr, "wrote patch to %s\n", patchFile)
		}
	}

	fmt.Fprintf(os.Stderr, "diffr: %d changed, %d added, %d removed (%d identical) in %s\n",
		res.Stats.Changed, res.Stats.Added, res.Stats.Removed, res.Stats.Identical, res.Stats.Duration.Round(time.Millisecond))
	return nil
}
