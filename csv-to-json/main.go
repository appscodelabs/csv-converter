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
	"strconv"
	"strings"

	"kmodules.xyz/client-go/logs"

	"github.com/gobuffalo/flect"
	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	inFile  string
	outDir  string
	renames = map[string]string{}
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "csv-to-json",
		Short: "Convert CSV files to json",
		RunE: func(cmd *cobra.Command, args []string) error {
			return convert()
		},
	}
	flags := rootCmd.Flags()

	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&inFile, "in", inFile, "Path to csv file")
	flags.StringVar(&outDir, "out", outDir, "Path to outDir directory")
	flags.StringToStringVar(&renames, "renames", nil, "Provide a map of column renames")

	logs.ParseFlags()

	utilruntime.Must(rootCmd.Execute())
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

func ValueFunc(v string) interface{} {
	v = strings.TrimSpace(v)
	smallV := strings.ToLower(v)
	if smallV == "true" || smallV == "t" || smallV == "y" || smallV == "yes" {
		return true
	}
	if smallV == "false" || smallV == "f" || smallV == "n" || smallV == "no" {
		return false
	}
	if i, err := strconv.Atoi(smallV); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f
	}
	return v
}

func convert() error {
	input, err := ioutil.ReadFile(inFile)
	if err != nil {
		return err
	}

	r := csv.NewReader(bytes.NewReader(input))

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	rows := make([]interface{}, 0, len(records))

	emailIdx := -1
	headers := make([]string, 0, len(records[0]))
	for i, entry := range records[0] {
		key := KeyFunc(entry)
		if key == "email" {
			emailIdx = i
		}
		headers = append(headers, key)
	}

	for _, r := range records[1:] {
		if r[emailIdx] == "" {
			continue
		}

		x := map[string]interface{}{}
		for i, v := range r {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			x[headers[i]] = ValueFunc(v)
		}
		rows = append(rows, x)
	}

	base := filepath.Base(inFile)
	ext := filepath.Ext(inFile)
	outFilename := filepath.Join(outDir, fmt.Sprintf("%s.json", strings.TrimSuffix(base, ext)))

	data, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outFilename, data, 0o644)
}
