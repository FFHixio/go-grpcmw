package descriptor

import (
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Parse parses the given protobuf request into a map of packages (key) and of
// files information (value).
func Parse(pb *plugin.CodeGeneratorRequest) (pkgs map[string][]*File, err error) {
	// TODO: Do this in multiple goroutines
	pkgs = make(map[string][]*File)
	for _, file := range pb.GetProtoFile() {
		if parsed, err := GetFile(file); err != nil {
			return nil, err
		} else if parsed != nil {
			pkgs[file.GetPackage()] = append(pkgs[file.GetPackage()], parsed)
		}
	}
	return
}
