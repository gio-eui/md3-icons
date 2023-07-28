// Copyright (c) 2023 https://github.com/gio-eui
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// SPDX-License-Identifier: MIT

//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gio-eui/ivgconv"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Convert all svg files from the assets folder to iconVG files

// srcFolder is the folder where the SVG files are located
const srcFolder = "./assets/"

// destFolder is the folder where the icon.go files will be generated
const destFolder = "./icons/"

type IconType string

const (
	Base     IconType = "materialicons"
	Outlined IconType = "materialiconsoutlined"
	Round    IconType = "materialiconsround"
	Sharp    IconType = "materialiconssharp"
	TwoTone  IconType = "materialiconstwotone"
)

// toIconType convert from string to IconType
func toIconType(s string) IconType {
	switch s {
	case "materialicons":
		return Base
	case "materialiconsoutlined":
		return Outlined
	case "materialiconsround":
		return Round
	case "materialiconssharp":
		return Sharp
	case "materialiconstwotone":
		return TwoTone
	default:
		return Base
	}
}

// ext convert from IconType to ext value
func (it IconType) ext() string {
	switch it {
	case Base:
		return ""
	case Outlined:
		return "outlined"
	case Round:
		return "round"
	case Sharp:
		return "sharp"
	case TwoTone:
		return "twotone"
	default:
		return ""
	}
}

// desc convert from IconType to desc value
func (it IconType) desc() string {
	switch it {
	case Base:
		return "Base"
	case Outlined:
		return "Outlined"
	case Round:
		return "Round"
	case Sharp:
		return "Sharp"
	case TwoTone:
		return "TwoTone"
	default:
		return ""
	}
}

// Icon is a struct that contains the icon info
type Icon struct {
	Category string
	Name     string
	Type     IconType
	Size     string
	Path     string

	IconName    string
	PackageName string
	PackagePath string
}

// IconError is a struct that contains the icon error info
type IconError struct {
	Icon  *Icon
	Error error
}

func main() {
	log.Println("Collecting icons...")
	icons, err := collectIcons()
	if err != nil {
		log.Fatal(err)
	}

	// Keep track of all errors
	var errs []IconError

	log.Println("Generating icons...")
	for i, icon := range icons {
		log.Printf("Generating icon %d/%d: %s %s %s\n", i+1, len(icons), icon.Category, icon.Name, icon.Type.ext())
		// Create the categories folder
		if err := os.Mkdir(filepath.Join(destFolder, icon.Category), os.ModePerm); err != nil && !os.IsExist(err) {
			errs = append(errs, IconError{Icon: icon, Error: err})
		}
		// Create the name folder
		if err := os.Mkdir(filepath.Join(destFolder, icon.Category, icon.Name), os.ModePerm); err != nil && !os.IsExist(err) {
			errs = append(errs, IconError{Icon: icon, Error: err})
		}
		// Create the type folder
		if err := os.Mkdir(filepath.Join(destFolder, icon.Category, icon.Name, icon.Type.ext()), os.ModePerm); err != nil && !os.IsExist(err) {
			errs = append(errs, IconError{Icon: icon, Error: err})
		}
		// Update the icon PackagePath
		icon.PackagePath = filepath.Join(destFolder, icon.Category, icon.Name, icon.Type.ext())
		// Create the icons.go file
		if err := generateIconVG(icon); err != nil {
			errs = append(errs, IconError{Icon: icon, Error: err})
		}
	}

	// If there are errors, print them
	if len(errs) > 0 {
		fmt.Println(strings.Repeat("-", 80))
		log.Printf("There are %d errors\n", len(errs))
		for _, err := range errs {
			log.Printf("%s %s %s: %s\n", err.Icon.Category, err.Icon.Name, err.Icon.Type.ext(), err.Error)
		}
		fmt.Println(strings.Repeat("-", 80))
	}
}

// collectIcons collect the icons from the assets folder
func collectIcons() ([]*Icon, error) {
	var icons []*Icon

	// List all files in the assets folder and its subfolders recursively which have
	// the .svg extension using filepath.WalkDir
	err := filepath.WalkDir(srcFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// If the file is not a directory and is named 24px.svg
		if !d.IsDir() && filepath.Base(path) == "24px.svg" {
			// Get the icon info
			icon := &Icon{Path: path}
			// Split the PackagePath
			split := strings.Split(filepath.ToSlash(path), "/")
			// Get the icon category
			icon.Category = strings.ToLower(split[len(split)-4])
			// Get the icon name
			icon.Name = strings.ToLower(split[len(split)-3])
			// Get the icon type
			icon.Type = toIconType(split[len(split)-2])
			// Get the icon size
			icon.Size = strings.TrimSuffix(filepath.Base(path), filepath.Ext(filepath.Base(path)))
			// Set the icon name
			icon.IconName = icoName(*icon)
			// Set the package name
			icon.PackageName = pkgNane(*icon)
			// Add the icon to the icons slice
			icons = append(icons, icon)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return icons, nil
}

func icoName(icon Icon) string {
	// Split the icon name whenever there is a dash or an underscore
	split := strings.FieldsFunc(icon.Name, func(r rune) bool {
		return r == '-' || r == '_'
	})
	// Loop through the split name to capitalize each word
	for i, s := range split {
		// lowercase the string
		split[i] = strings.ToLower(s)
		// capitalize the first letter
		split[i] = strings.ToUpper(string(s[0])) + s[1:]
	}
	// Join the split name
	return strings.Join(split, "")
}

func pkgNane(icon Icon) string {
	pn := fmt.Sprintf("mdi%s%s%s",
		cases.Title(language.English).String(icon.Category),
		icon.IconName,
		cases.Title(language.English).String(icon.Type.ext()),
	)
	return pn
}

func generateIconVG(icon *Icon) error {
	// Convert the SVG file to IconVG
	ivgData, err := ivgconv.FromFile(icon.Path)
	if err != nil {
		return err
	}

	// Prepare the buffer
	out := new(bytes.Buffer)

	// Write the header
	fmt.Fprintf(out, "// generated by go run genicons.go; DO NOT EDIT\n\n")
	fmt.Fprintf(out, "// Package %s provides the (%s) %s icon\n", icon.PackageName, icon.Type.desc(), icon.Name)
	fmt.Fprintf(out, "//\n")
	fmt.Fprintf(out, "package %s\n\n", icon.PackageName)

	// Store the image bytes in a variable
	fmt.Fprintf(out, "var Ivg = []byte{")
	for i, b := range ivgData {
		if i%16 == 0 {
			fmt.Fprintf(out, "\n\t")
		}
		fmt.Fprintf(out, "0x%02x, ", b)
	}
	fmt.Fprintf(out, "}\n")

	// Format the code
	src, err := format.Source(out.Bytes())
	if err != nil {
		return err
	}

	// Write icon file
	if err := os.WriteFile(filepath.Join(icon.PackagePath, "icon.go"), src, os.ModePerm); err != nil {
		return err
	}
	return nil
}
