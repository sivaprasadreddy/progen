package crud

import (
	"fmt"
	"os"
	"path/filepath"
)

func Generate(entity string) {
	dir := "generated"

	_ = os.MkdirAll(dir, 0755)

	filename := filepath.Join(dir, entity+".java")

	content := fmt.Sprintf(
		"public class %s {\n}\n",
		entity,
	)

	_ = os.WriteFile(filename, []byte(content), 0644)
}
