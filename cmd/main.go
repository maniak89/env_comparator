package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"

	"github.com/maniak89/env_comparator/internal/comparator"
)

func main() {
	logger := log.With().Logger()
	ctx := logger.WithContext(context.Background())

	dir1 := flag.String("dir1", "", "directory with configs")
	dir2 := flag.String("dir2", "", "another directory with configs")
	notPresentSide := flag.String("np", "dir1/dir2", "filter only not present")
	flag.Parse()
	if dir1 == nil || *dir1 == "" {
		log.Fatal().Msg("dir1 is not set")
		return
	}
	if dir2 == nil || *dir2 == "" {
		log.Fatal().Msg("dir2 is not set")
		return
	}

	filesDir1, err := ioutil.ReadDir(*dir1)
	if err != nil {
		logger.Fatal().Err(err).Str("dir", *dir1).Msg("Failed list files in directory")
		return
	}
	var tableData [][]string
	for _, f := range filesDir1 {
		f1 := path.Join(*dir1, f.Name())
		if strings.HasSuffix(f1, ".yml") || strings.HasSuffix(f1, ".yaml") {
			f2 := path.Join(*dir2, f.Name())
			result, err := comparator.CompareYaml(ctx, f1, f2)
			if err != nil {
				logger.Error().Err(err).Str("file1", f1).Str("file2", f2).Msg("failed compare files")
				continue
			}
			if len(result) > 0 {
				for _, r := range result {
					for _, e := range r.Envs {
						if notPresentSide != nil && *notPresentSide != "" {
							if *notPresentSide == "dir1" && e.Val1 == "" {
								tableData = append(tableData, []string{r.Name, e.Name, e.Val1, e.Val2})
							}
							if *notPresentSide == "dir2" && e.Val2 == "" {
								tableData = append(tableData, []string{r.Name, e.Name, e.Val1, e.Val2})
							}
						} else {
							tableData = append(tableData, []string{r.Name, e.Name, e.Val1, e.Val2})
						}
					}
				}
			}
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Service", "Name", "val1", "val2"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, v := range tableData {
		table.Append(v)
	}
	table.Render() // Send output
}
