/*
Copyright AppsCode Inc. and Contributors

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/agpl-3.0.txt>.
*/

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

var in string

func main() {
	rootCmd := &cobra.Command{
		Use:   "listmonk-csv-to-json",
		Short: "Convert Listmonk ready csv file to json file",
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
	r := csv.NewReader(bytes.NewReader(data))
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	entries := make([]Row, 0, len(records))
	for _, record := range records {
		entry := map[string]interface{}{}
		if len(record) >= 3 {
			err := json.Unmarshal([]byte(record[2]), &entry)
			if err != nil {
				return err
			}
		}
		entry["name"] = record[1]
		entry["email"] = record[0]

		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i]["email"].(string) < entries[j]["email"].(string)
	})

	data, err = json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	outFile := fmt.Sprintf("%s_listmonk.csv", strings.TrimSuffix(filepath.Base(in), filepath.Ext(in)))
	return ioutil.WriteFile(filepath.Join(filepath.Dir(in), outFile), data, 0o644)
}
