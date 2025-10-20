/**
 * MIT License
 *
 * Copyright (c) 2025 Viki (VN - initials of my first and last name)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 *
 *
 * Contact: contact@viki.design  
 * Website: https://www.viki.design
 * 
 */

package bashgo

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Append(filename string, data interface{}) error {
	if info, err := os.Stat(filename); err == nil && info.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", filename)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	switch v := data.(type) {
	case string:
		_, err = file.WriteString(strings.TrimRight(v, "\n") + "\n")
	case []string:
		for _, line := range v {
			_, err = file.WriteString(strings.TrimRight(line, "\n") + "\n")
			if err != nil {
				break
			}
		}
	default:
		return fmt.Errorf("data must be string or []string")
	}

	return err
}

func Write(filename string, data interface{}, linenr int) error {
	if info, err := os.Stat(filename); err == nil && info.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", filename)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.WriteFile(filename, []byte(""), 0644)
	}

	content := []string{}
	switch v := data.(type) {
	case string:
		content = []string{v}
	case []string:
		content = v
	default:
		return fmt.Errorf("data must be string or []string")
	}

	existing := []string{}
	if b, err := os.ReadFile(filename); err == nil {
		existing = strings.Split(string(b), "\n")
		if len(existing) > 0 && existing[len(existing)-1] == "" {
			existing = existing[:len(existing)-1]
		}
	}

	if linenr < 0 {
		return os.WriteFile(filename, []byte(strings.Join(content, "\n")+"\n"), 0644)
	}

	for len(existing) < linenr {
		existing = append(existing, "")
	}
	for i, line := range content {
		if linenr+i < len(existing) {
			existing[linenr+i] = line
		} else {
			existing = append(existing, line)
		}
	}

	return os.WriteFile(filename, []byte(strings.Join(existing, "\n")+"\n"), 0644)
}

func Prepend(filename string, data interface{}) error {
	if info, err := os.Stat(filename); err == nil && info.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", filename)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.WriteFile(filename, []byte(""), 0644)
	}

	content := []string{}
	switch v := data.(type) {
	case string:
		content = []string{v}
	case []string:
		content = v
	default:
		return fmt.Errorf("data must be string or []string")
	}

	existing := []string{}
	if b, err := os.ReadFile(filename); err == nil {
		existing = strings.Split(string(b), "\n")
		if len(existing) > 0 && existing[len(existing)-1] == "" {
			existing = existing[:len(existing)-1]
		}
	}

	all := append(content, existing...)
	return os.WriteFile(filename, []byte(strings.Join(all, "\n")+"\n"), 0644)
}

func Read(filename string, linenr int) ([]string, error) {
	if info, err := os.Stat(filename); err == nil && info.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory, not a file", filename)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if linenr < 0 {
		return lines, nil
	} else if linenr < len(lines) {
		return []string{lines[linenr]}, nil
	} else {
		return nil, nil
	}
}

func Remove(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(filename)
}

func List(path string) ([]fs.FileInfo, error) {
	if path == "" {
		cwd, _ := os.Getwd()
		path = cwd
	}
	return ioutil.ReadDir(path)
}

func ListDirs(path string) ([]string, error) {
	items, err := List(path)
	if err != nil {
		return nil, err
	}
	dirs := []string{}
	for _, item := range items {
		if item.IsDir() {
			dirs = append(dirs, item.Name())
		}
	}
	return dirs, nil
}

func ListFiles(path string) ([]string, error) {
	items, err := List(path)
	if err != nil {
		return nil, err
	}
	files := []string{}
	for _, item := range items {
		if !item.IsDir() {
			files = append(files, item.Name())
		}
	}
	return files, nil
}

func Mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

func Size(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		var total int64
		err := filepath.Walk(path, func(_ string, f os.FileInfo, err error) error {
			if err == nil && !f.IsDir() {
				total += f.Size()
			}
			return nil
		})
		return total, err
	}
	return info.Size(), nil
}

func SizeKB(path string) (float64, error) {
	size, err := Size(path)
	return float64(size) / 1024.0, err
}

func SizeMB(path string) (float64, error) {
	size, err := Size(path)
	return float64(size) / (1024.0 * 1024.0), err
}

func SizeGB(path string) (float64, error) {
	size, err := Size(path)
	return float64(size) / (1024.0 * 1024.0 * 1024.0), err
}

func BasePath(path string) string {
	return filepath.Dir(path)
}

func LeafNode(path string) string {
	return filepath.Base(path)
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsEmpty(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	if info.IsDir() {
		files, _ := ioutil.ReadDir(path)
		return len(files) == 0
	} else if info.Mode().IsRegular() {
		return info.Size() == 0
	}
	return false
}

