package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"

	"kmodules.xyz/client-go/logs"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	base      string
	others    []string
	overwrite bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "merge-json",
		Short: "Merge other json files into a base json file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return merge()
		},
	}
	flags := rootCmd.Flags()

	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&base, "base", base, "Path to base json file")
	flags.StringSliceVar(&others, "others", others, "Path to other json file")
	flags.BoolVar(&overwrite, "overwrite", overwrite, "If true, merge override non-empty base attributes with non-empty other attributes values.")

	logs.ParseFlags()

	utilruntime.Must(rootCmd.Execute())
}

type Row map[string]interface{}

func merge() error {
	baseEntries, err := LoadFile(base)
	if err != nil {
		return err
	}
	for _, otherfile := range others {
		otherEntries, err := LoadFile(otherfile)
		if err != nil {
			return err
		}

		for email, other := range otherEntries {
			existing, found := baseEntries[email]
			if !found {
				baseEntries[email] = other
			} else {
				err = mergo.Map(existing, other, func(config *mergo.Config) {
					config.Overwrite = overwrite
				})
				if err != nil {
					return err
				}
				baseEntries[email] = existing
			}
		}
	}

	out := make([]Row, len(baseEntries))
	for _, entry := range baseEntries {
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i]["email"].(string) < out[j]["email"].(string)
	})
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(base, data, 0644)
}

// email -> Row
func LoadFile(filename string) (map[string]Row, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var entries []Row
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}
	out := make(map[string]Row, len(entries))
	for _, x := range entries {
		email, ok := x["email"]
		if !ok {
			return nil, fmt.Errorf("email missing is %+v", x)
		}
		out[email.(string)] = x
	}
	return out, nil
}
