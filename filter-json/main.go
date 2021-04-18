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

	"github.com/gobuffalo/flect"
	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	in        string
	renames   = map[string]string{}
	keys      []string
	blocklist string
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
	flags.StringVar(&in, "in", in, "Path to input json file")
	flags.StringToStringVar(&renames, "renames", nil, "Provide a map of column renames")
	flags.StringSliceVar(&keys, "keys", keys, "Keys to be kept")
	flags.StringVar(&blocklist, "blocklist", blocklist, "Path to block list json file. Matching emails from this file will be removed")

	logs.ParseFlags()

	utilruntime.Must(rootCmd.Execute())
}

type Row map[string]interface{}

func filter() error {
	blocked := map[string]bool{}
	if blocklist != "" {
		entries, err := LoadFile(blocklist)
		if err != nil {
			return err
		}
		for _, row := range entries {
			email, ok := row["email"]
			if !ok {
				continue
			}
			blocked[email.(string)] = true
		}
	}

	entries, err := LoadFile(in)
	if err != nil {
		return err
	}

	out := make([]Row, 0, len(entries))
	for _, row := range entries {

		email, ok := row["email"]
		if !ok {
			continue
		}
		if blocked[email.(string)] {
			continue
		}

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

	dir := filepath.Dir(in)
	filename := filepath.Join(dir, fmt.Sprintf("%s_filtered.json", strings.TrimSuffix(filepath.Base(in), filepath.Ext(in))))
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

	out := make([]Row, 0, len(entries))
	for _, row := range entries {
		x := map[string]interface{}{}
		for k, v := range row {
			x[KeyFunc(k)] = v
		}
		out = append(out, x)
	}

	return out, nil
}

func KeyFunc(key string) string {
	if replace, ok := renames[key]; ok {
		return replace
	}
	key = flect.Underscore(key)
	if strings.HasPrefix(key, "email") {
		return "email"
	}
	if replace, ok := renames[key]; ok {
		return replace
	}
	return key
}
