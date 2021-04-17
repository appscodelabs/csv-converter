package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"kmodules.xyz/client-go/logs"

	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	base string
	keys []string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "filter-json",
		Short: "Filter keys from json",
		RunE: func(cmd *cobra.Command, args []string) error {
			return filter()
		},
	}
	flags := rootCmd.Flags()

	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&base, "base", base, "Path to base json file")
	flags.StringSliceVar(&keys, "keys", keys, "Keys to be kept")

	logs.ParseFlags()

	utilruntime.Must(rootCmd.Execute())
}

type Row map[string]interface{}

func filter() error {
	entries, err := LoadFile(base)
	if err != nil {
		return err
	}

	out := make([]Row, 0, len(entries))
	for _, row := range entries {
		filtereRow := Row{}
		for _, key := range keys {
			if v, ok := row[key]; ok {
				filtereRow[key] = v
			}
		}
		out = append(out, filtereRow)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i]["email"].(string) < out[j]["email"].(string)
	})
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(base)
	filename := filepath.Join(dir, fmt.Sprintf("%s_filtered.json", strings.TrimSuffix(filepath.Base(base), filepath.Ext(base))))
	return ioutil.WriteFile(filename, data, 0644)
}

// email -> Row
func LoadFile(filename string) ([]Row, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var entries []Row
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
