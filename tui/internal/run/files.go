package run

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"high-tui/internal/model"
)

var refRE = regexp.MustCompile(`^(.+)::([A-Za-z_][A-Za-z0-9_]*)::([A-Za-z0-9_]+)$`)
var varRE = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)=(.*)$`)

// Resolve replaces chemin::NOM::label and chemin::NOM::value. Other text stays.
func Resolve(text, root, where string) (string, error) {
	if strings.Count(text, "::") < 2 {
		return text, nil
	}
	m := refRE.FindStringSubmatch(text)
	if m == nil {
		return text, nil
	}
	path, name, kind := m[1], m[2], m[3]
	if kind != "label" && kind != "value" {
		return "", model.Errf("le dernier mot n'est ni label ni value (%s)", where)
	}
	file := path
	if !filepath.IsAbs(file) {
		file = filepath.Join(root, path)
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		return "", model.Errf("fichier introuvable : %s", path)
	}
	label, value := "", ""
	haveValue := false
	labelRE := regexp.MustCompile(`^# label ` + regexp.QuoteMeta(name) + ` (.*)$`)
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimRight(line, "\r")
		if vm := varRE.FindStringSubmatch(line); vm != nil && vm[1] == name && !haveValue {
			value = vm[2]
			haveValue = true
		}
		if lm := labelRE.FindStringSubmatch(line); lm != nil && label == "" {
			label = lm[1]
		}
	}
	if kind == "value" {
		if !haveValue {
			return "", model.Errf("variable manquante (%s)", where)
		}
		return value, nil
	}
	if label == "" {
		return "", model.Errf("label manquant (%s)", where)
	}
	return label, nil
}

// ReadVar returns the right-hand side of NAME= in file.
func ReadVar(file, name string) (string, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return "", model.Errf("fichier introuvable : %s", file)
	}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimRight(line, "\r")
		m := varRE.FindStringSubmatch(line)
		if m != nil && m[1] == name {
			return m[2], nil
		}
	}
	return "", nil
}

// WriteVar replaces or appends NAME=value and keeps every other line.
func WriteVar(file, name, value string) error {
	raw, err := os.ReadFile(file)
	if err != nil {
		return model.Errf("fichier introuvable : %s", file)
	}
	lines := strings.Split(string(raw), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	found := false
	for i, line := range lines {
		body := strings.TrimRight(line, "\r")
		m := varRE.FindStringSubmatch(body)
		if m != nil && m[1] == name {
			lines[i] = name + "=" + value
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, name+"="+value)
	}
	return os.WriteFile(file, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

// Path joins root and a relative file. Absolute paths stay absolute.
func Path(root, file string) string {
	if filepath.IsAbs(file) {
		return file
	}
	return filepath.Join(root, file)
}
