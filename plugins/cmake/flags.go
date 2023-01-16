package cmake

import (
	"fmt"
	"strings"
)

var (
	flagHandlers map[string]flagHandler = map[string]flagHandler{}
)

type flagHandler func(string, []string) ([]string, error)

func handleSingleFlag(cmakeFlag string, key string, values []string) ([]string, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("too many values for key '%v'", key)
	}
	return []string{fmt.Sprintf("-D%v=%v", cmakeFlag, values[0])}, nil
}

func handleJoinFlag(cmakeFlag string, joiner string, key string, values []string) ([]string, error) {
	return []string{fmt.Sprintf("-D%v=%v", cmakeFlag, strings.Join(values, joiner))}, nil
}

func handleCompilerFlags(cmakeFlag string, key string, values []string) ([]string, error) {
	rawFlags := strings.Join(values, " ")
	return []string{fmt.Sprintf("-D%v=%v", cmakeFlag, rawFlags)}, nil
}

func init() {
	singleFlags := [][2]string{
		{"cmake.build_type", "CMAKE_BUILD_TYPE"},
		{"cmake.cc", "CMAKE_C_COMPILER"},
		{"cmake.cxx", "CMAKE_CXX_COMPILER"},
		{"cmake.toolchain_file", "CMAKE_TOOLCHAIN_FILE"},
	}
	for i := range singleFlags {
		key, cmakeArg := singleFlags[i][0], singleFlags[i][1]
		flagHandlers[key] = func(key string, values []string) ([]string, error) {
			return handleSingleFlag(cmakeArg, key, values)
		}
	}

	flagHandlers["cmake.generator"] = func(key string, values []string) ([]string, error) {
		if len(values) != 1 {
			return nil, fmt.Errorf("too many values for key '%v'", key)
		}
		return []string{"-G", values[0]}, nil
	}

	joinFlags := [][3]string{
		{"cmake.module_path", "CMAKE_MODULE_PATH", ";"},
	}
	for i := range joinFlags {
		key, cmakeArg, joiner := joinFlags[i][0], joinFlags[i][1], joinFlags[i][2]
		flagHandlers[key] = func(key string, values []string) ([]string, error) {
			return handleJoinFlag(cmakeArg, joiner, key, values)
		}
	}

	variantFlagsBase := [][2]string{
		{"cmake.cflags", "CMAKE_C_FLAGS"},
		{"cmake.cxxflags", "CMAKE_CXX_FLAGS"},
	}

	linkTypes := [][2]string{
		{"exe", "EXE"},
		{"module", "MODULE"},
		{"shared", "SHARED"},
		{"static", "STATIC"},
	}

	for i := range linkTypes {
		key := fmt.Sprintf("cmake.ldflags.%v", linkTypes[i][0])
		flag := fmt.Sprintf("CMAKE_%v_LINKER_FLAGS", linkTypes[i][1])
		variantFlagsBase = append(variantFlagsBase, [2]string{key, flag})
	}

	variantFlagsSuffix := [][2]string{
		{"debug", "DEBUG"},
		{"minsizerel", "MINSIZEREL"},
		{"release", "RELEASE"},
		{"relwithdebinfo", "RELWITHDEBINFO"},
	}

	for i := range variantFlagsBase {
		key, cmakeArg := variantFlagsBase[i][0], variantFlagsBase[i][1]
		flagHandlers[key] = func(key string, values []string) ([]string, error) {
			return handleCompilerFlags(cmakeArg, key, values)
		}
		for j := range variantFlagsSuffix {
			fullKey := fmt.Sprintf("%v.%v", key, variantFlagsSuffix[j][0])
			fullArgument := fmt.Sprintf("%v_%v", cmakeArg, variantFlagsSuffix[j][1])
			flagHandlers[fullKey] = func(key string, values []string) ([]string, error) {
				return handleCompilerFlags(fullArgument, key, values)
			}
		}
	}
}
