package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
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
		Use:   "json-tats",
		Short: "Print json stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return LoadFile(in)
		},
	}
	flags := rootCmd.Flags()

	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&in, "in", in, "Path to input json file")
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
	var entries []Row
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return err
	}

	total := len(entries)
	n := 0
	for _, row := range entries {
		for k, v := range row {
			if KeyFunc(k) == "email" && !reflect.ValueOf(v).IsZero() {
				n++
				break
			}
		}
	}

	fmt.Printf("Total objects: %d\n" ,total)
	fmt.Printf("Total objects with email: %d\n" ,n)

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
