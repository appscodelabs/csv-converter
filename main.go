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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"kmodules.xyz/client-go/logs"

	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

type Location struct {
	App     string
	Version string
}

func main() {
	var (
		input  string
		output string
	)
	var rootCmd = &cobra.Command{
		Use:   "csv-converter",
		Short: "Convert CSV files to Listmonk format",
		RunE: func(cmd *cobra.Command, args []string) error {
			return convert(output, input)
		},
	}
	flags := rootCmd.Flags()

	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&input, "in", input, "Path to csv file")
	flags.StringVar(&output, "out", output, "Path to output directory")

	logs.ParseFlags()

	utilruntime.Must(rootCmd.Execute())
}

func convert(outDir, in string) error {
	data, err := ioutil.ReadFile(in)
	if err != nil {
		return err
	}

	r := csv.NewReader(bytes.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(records)

	base := filepath.Base(in)
	ext := filepath.Ext(in)
	outFilename := filepath.Join(outDir, fmt.Sprintf("%s_listmonk.%s", strings.TrimSuffix(base, ext), ext))

	return ioutil.WriteFile(outFilename, data, 0644)
}
