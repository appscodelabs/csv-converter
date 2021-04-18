package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"kmodules.xyz/client-go/logs"

	"github.com/gobuffalo/flect"
	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	in        string
	renames = map[string]string{}
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "csv-tats",
		Short: "Print csv stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return LoadFile(in)
		},
	}
	flags := rootCmd.Flags()

	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&in, "in", in, "Path to input csv file")
	flags.StringToStringVar(&renames, "renames", nil, "Provide a map of column renames")

	logs.ParseFlags()

	utilruntime.Must(rootCmd.Execute())
}

type Row map[string]interface{}

func LoadFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	r := csv.NewReader(bytes.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	total := len(records)-1

	emailIdx := -1
	for i, entry := range records[0] {
		key := KeyFunc(entry)
		if key == "email" {
			emailIdx = i
		}
	}

	n := 0
	for _, r := range records[1:] {
		if r[emailIdx] == "" {
			continue
		}
		n++
	}

	fmt.Printf("Total rows: %d\n" ,total)
	fmt.Printf("Total rows with email: %d\n" ,n)

	return nil
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
