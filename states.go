// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package states

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	Flags flag.FlagSet

	outFile            *string
	includeUnspecified *bool
	stateName          = regexp.MustCompile("^[A-Za-z]*State$")
)

func init() {
	outFile = Flags.String("out_file", "", "")
	includeUnspecified = Flags.Bool("include_unspecified", true, "")
}

func Analyze(plugin *protogen.Plugin) error {
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	var out io.Writer
	if f := *outFile; f != "" {
		file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		out = file
	} else {
		out = os.Stderr
	}
	w := func(f string, args ...any) {
		out.Write([]byte(fmt.Sprintf(f, args...) + "\n"))
	}

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		var states []*protogen.Enum

		nested := collectNestedEnums(file.Messages)
		enums := append(file.Enums, nested...)
		for _, e := range enums {
			if stateName.MatchString(string(e.Desc.Name())) {
				states = append(states, e)
			}
		}
		for _, state := range states {
			var names []string
			for _, v := range state.Values {
				n := string(v.Desc.Name())
				if !*includeUnspecified && strings.HasSuffix(n, "_UNSPECIFIED") {
					continue
				}
				names = append(names, n)
			}

			// proto pkg, enum name, values
			w("%s,%s,%s", file.Desc.Package(), state.Desc.Name(), strings.Join(names, ","))
		}
	}

	return nil
}

func collectNestedEnums(messages []*protogen.Message) (enums []*protogen.Enum) {
	for _, m := range messages {
		enums = append(enums, m.Enums...)
		enums = append(enums, collectNestedEnums(m.Messages)...)
	}
	return
}
