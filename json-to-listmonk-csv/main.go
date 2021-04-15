package main

import (
	"bytes"
	"encoding/csv"
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
	in string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "json-to-listmonk-csv",
		Short: "Convert json file to Listmonk ready csv file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return convert()
		},
	}
	flags := rootCmd.Flags()

	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&in, "in", in, "Path to input json file")

	logs.ParseFlags()

	utilruntime.Must(rootCmd.Execute())
}

type Row map[string]interface{}

func convert() error {
	data, err := ioutil.ReadFile(in)
	if err != nil {
		return err
	}
	var entries []Row
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return err
	}

	records := make([][]string, 0, len(entries))
	for _, entry := range entries {
		record := make([]string, 3)

		if name, ok := entry["email"]; ok {
			record[0] = name.(string)
			delete(entry, "email")
		} else {
			return fmt.Errorf("email missing is %+v", entry)
		}

		if name, ok := entry["name"]; ok {
			record[1] = name.(string)
			delete(entry, "name")
		} else {
			nameparts := make([]string, 0, 2)
			if first, found := entry["first_name"]; found {
				nameparts = append(nameparts, first.(string))
				delete(entry, "first_name")
			}
			if last, found := entry["last_name"]; found {
				nameparts = append(nameparts, last.(string))
				delete(entry, "last_name")
			}
			record[1] = strings.Join(nameparts, " ")
		}

		rest, err := json.Marshal(record)
		if err != nil {
			return err
		}
		record[2] = string(rest)

		records = append(records, record)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i][0] < records[j][0]
	})

	var buf bytes.Buffer

	w := csv.NewWriter(&buf)
	_ = w.WriteAll(records) // calls Flush internally
	if err := w.Error(); err != nil {
		return err
	}

	outFile := fmt.Sprintf("%s_listmonk.csv", strings.TrimSuffix(filepath.Base(in), filepath.Ext(in)))
	return ioutil.WriteFile(filepath.Join(filepath.Dir(in), outFile), buf.Bytes(), 0644)
}
