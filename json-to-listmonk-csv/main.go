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
	"regexp"
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

		if email, ok := entry["email"]; ok {
			record[0] = email.(string)
			delete(entry, "email")
		} else {
			return fmt.Errorf("email missing is %+v", entry)
		}

		if name, ok := entry["name"]; ok && IsString(name) {
			record[1] = name.(string)
			delete(entry, "name")
			delete(entry, "first_name")
			delete(entry, "last_name")
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
			if len(nameparts) > 0 {
				record[1] = strings.Join(nameparts, " ")
			} else {
				record[1] = DetectNameFromEmail(record[0])
			}
		}

		rest, err := json.Marshal(entry)
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
	_ = w.Write([]string{"email", "name", "attributes"})
	_ = w.WriteAll(records) // calls Flush internally
	if err := w.Error(); err != nil {
		return err
	}

	outFile := fmt.Sprintf("%s_listmonk.csv", strings.TrimSuffix(filepath.Base(in), filepath.Ext(in)))
	return ioutil.WriteFile(filepath.Join(filepath.Dir(in), outFile), buf.Bytes(), 0644)
}

func DetectNameFromEmail(email string) string {
	idx := strings.IndexRune(email, '@')
	if idx > -1 {
		if idx <= 2 { // [2 letters]@domain.com, detect name from domain.com
			email = email[idx+1:]
			idx2 := strings.IndexRune(email, '.')
			if idx2 > -1 {
				email = email[0:idx2]
			}
		} else {
			name := email[0:idx]
			re := regexp.MustCompile("\\d+")
			name = re.ReplaceAllString(name, "")
			// in case numbers@qq.com type email
			if name != "" {
				email = name
			}
		}
	}

	email = strings.Replace(email, ".", " ", -1)
	email = strings.Replace(email, "_", " ", -1)
	email = strings.Replace(email, "-", " ", -1)
	return strings.Title(email)
}

func IsString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}
