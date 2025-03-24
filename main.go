package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	for i := 3; i < 28; i++ {
		var comprocs []string

		bytes, err := os.ReadFile(fmt.Sprintf(`pkg\webview2\ICoreWebView2_%v.go`, i))
		if err != nil {
			continue
		}

		firstSplit := fmt.Sprintf(`type ICoreWebView2_%vVtbl struct {`, i)
		secondSplit := `}`

		fileString := string(bytes)

		fileSplit := strings.Split(fileString, firstSplit)
		if len(fileSplit) < 2 {
			continue
		}

		fileSplit = strings.Split(fileSplit[1], secondSplit)
		if len(fileSplit) < 2 {
			continue
		}

		comprocsString := fileSplit[0]

		comprocsString = strings.ReplaceAll(comprocsString, "IUnknownVtbl", "")
		comprocsString = strings.TrimSpace(comprocsString)

		comprocsString = strings.ReplaceAll(comprocsString, "\r", "")
		comprocsString = strings.ReplaceAll(comprocsString, "ComProc", "")
		//comprocsString = strings.ReplaceAll(comprocsString, "\n", "")
		comprocsString = strings.ReplaceAll(comprocsString, "\t", "")
		comprocsString = strings.ReplaceAll(comprocsString, ` `, "")

		lines := strings.Split(comprocsString, "\n")

		for _, line := range lines {
			if line == "" {
				continue
			}

			comprocs = append(comprocs, line)
		}

		var comprocsStructTypes string
		for _, comproc := range comprocs {
			comprocsStructTypes += fmt.Sprintf("%v ComProc\n", comproc)
		}

		goFileString := fmt.Sprintf(fileTemplate, i, i-1, comprocsStructTypes, i, i, i, i, i, i, i, i, i, i, i, i)

		os.WriteFile(fmt.Sprintf("ICoreWebView2_%v.go", i), []byte(goFileString), 0644)
	}
}

var fileTemplate = `//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_%vVtbl struct {
	iCoreWebView2_%vVtbl
	%v
}

type ICoreWebView2_%v struct {
	vtbl *iCoreWebView2_%vVtbl
}

func (i *ICoreWebView2_%v) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_%v) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_%v() *ICoreWebView2_%v {
	var result *ICoreWebView2_%v

	iidICoreWebView2_%v := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_%v)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_%v() *ICoreWebView2_%v {
	return e.webview.GetICoreWebView2_%v()
}
`
