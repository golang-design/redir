// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var xTmpl = template.Must(template.New("x").ParseFiles("public/x.html"))

type x struct {
	ImportRoot string
	VCS        string
	VCSRoot    string
	Suffix     string
}

// xHandler redirect returns an HTTP handler that redirects requests for
// the tree rooted at importPath to pkg.go.dev pages for those import paths.
// The redirections include headers directing `go get.` to satisfy the
// imports by checking out code from repoPath using the given VCS.
// As a special case, if both importPath and repoPath end in /*, then
// the matching element in the importPath is substituted into the repoPath
// specified for `go get.`
func (s *server) xHandler(vcs, importPath, repoPath string) http.Handler {
	wildcard := false
	if strings.HasSuffix(importPath, "/*") && strings.HasSuffix(repoPath, "/*") {
		wildcard = true
		importPath = strings.TrimSuffix(importPath, "/*")
		repoPath = strings.TrimSuffix(repoPath, "/*")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/.ping") {
			fmt.Fprintf(w, "pong")
			return
		}
		path := strings.TrimSuffix(req.Host+req.URL.Path, "/")
		var importRoot, repoRoot, suffix string
		if wildcard {
			if path == importPath {
				http.Redirect(w, req, "https://pkg.go.dev/"+importPath, 302)
				return
			}
			if !strings.HasPrefix(path, importPath+"/") {
				http.NotFound(w, req)
				return
			}
			elem := path[len(importPath)+1:]
			if i := strings.Index(elem, "/"); i >= 0 {
				elem, suffix = elem[:i], elem[i:]
			}
			importRoot = importPath + "/" + elem
			repoRoot = repoPath + "/" + elem
		} else {
			if path != importPath && !strings.HasPrefix(path, importPath+"/") {
				http.NotFound(w, req)
				return
			}
			importRoot = importPath
			repoRoot = repoPath
			suffix = path[len(importPath):]
		}
		d := &x{
			ImportRoot: importRoot,
			VCS:        vcs,
			VCSRoot:    repoRoot,
			Suffix:     suffix,
		}
		var buf bytes.Buffer
		err := xTmpl.Execute(&buf, d)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.Write(buf.Bytes())
	})
}
