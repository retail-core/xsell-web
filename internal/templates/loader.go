package templates

import (
	"html/template"
	"io/fs"
	"path/filepath"
)

func Load() (*template.Template, error) {
	t := template.New("").Funcs(template.FuncMap{
		"dict": func(values ...any) map[string]any {
			d := map[string]any{}
			for i := 0; i < len(values); i += 2 {
				d[values[i].(string)] = values[i+1]
			}
			return d
		},

		"js": func(s string) template.JS {
			return template.JS(s)
		},
	})

	var files []string
	err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// log.Println("Parsed template files:", files)

	return t.ParseFiles(files...)
}
