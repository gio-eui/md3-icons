# Material design 3 icons for [gioui](https://gioui.org)

[![GoDoc](https://godoc.org/github.com/gio-eui/md3-icons?status.svg)](https://godoc.org/github.com/gio-eui/md3-icons)
[![Go Report Card](https://goreportcard.com/badge/github.com/gio-eui/md3-icons)](https://goreportcard.com/report/github.com/gio-eui/md3-icons)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://raw.githubusercontent.com/gio-eui/md3-icons/master/LICENSE)


This repository contains the material design 3 icons converted to [IconVG](https://github.com/golang/exp/tree/master/shiny/iconvg) format.

## Usage

Each icon have its own package, so you can import only those icons that you need.

```go
    import "github.com/gio-eui/md3-icons/icons/toggle/check_box"

    var CheckBox *widget.Icon
    CheckBox, _ = widget.NewIcon(mdiToggleCheckBox.Ivg)
```
