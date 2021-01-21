package util

import (
	"bufio"
	"bytes"
	"strings"
)

const yamlSeparator = "---"

func IsYamlFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml")
}

func SplitYAMLDocument(document []byte) (documents [][]byte) {
	scanner := bufio.NewScanner(bytes.NewReader(document))
	scanner.Split(splitYAMLDocument)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		documents = append(documents, scanner.Bytes())
	}
	return documents
}

func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// We have a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// we can't read any more characters
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
