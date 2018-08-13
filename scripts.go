//go:generate go run scripts.go

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	goPath = os.Getenv("GOPATH")
)

const (
	currentDir = "."      // Current directory
	vendorDir  = "vendor" // Vendor directory
	protoExt   = ".proto" // Proto file extension
	protoc     = "protoc" // Proto compiler
)

//func main() {
//	if err := generateProtos(currentDir); err != nil {
//		fmt.Printf("%+v", err)
//	}
//}

func generateProtos(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Skip vendor directory if any
		if info.IsDir() && info.Name() == vendorDir {
			return filepath.SkipDir
		}

		// Find all the protobuf files
		if filepath.Ext(path) == protoExt {
			// Prepare the args
			args := []string{
				"-I=.",
				fmt.Sprintf("-I=%s", filepath.Join(goPath, "src")),
				fmt.Sprintf("-I=%s", filepath.Join(goPath, "src", "github.com", "gogo", "protobuf", "protobuf")),
				fmt.Sprintf("--proto_path=%s", filepath.Join(goPath, "src", "github.com")),
				"--gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:.",
				path,
			}

			// fmt.Println(args)
			// Execute protoc one by one
			cmd := exec.Command(protoc, args...)
			if err := cmd.Run(); err != nil {
				return err
			}
		}

		return nil
	})
}
