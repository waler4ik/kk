package pathchecker

import "strings"

type PathChecker struct {
	pathbase, templateFolder string
	exclusivePaths           []string
}

func New(pathBase, templateFolder string, exclusivePaths ...string) PathChecker {
	return PathChecker{
		pathbase:       pathBase,
		templateFolder: templateFolder,
		exclusivePaths: exclusivePaths,
	}
}

func (pc PathChecker) Skip(path string) bool {
	if path == pc.templateFolder {
		return true
	}
	if len(pc.exclusivePaths) > 0 && !contains(path, pc.exclusivePaths) {
		return true
	}
	return false
}

func (pc PathChecker) Rename(path string) (renamed string) {
	return strings.Replace(path, "{{PathBase}}", pc.pathbase, -1)
}

func contains(t string, s []string) bool {
	for _, p := range s {
		if p == t {
			return true
		}
	}
	return false
}
