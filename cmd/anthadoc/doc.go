// antha/cmd/anthadoc/doc.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

anthadoc extracts and generates documentation for antha programs.

It has two modes.

Without the -http flag, it runs in command-line mode and prints plain text
documentation to standard output and exits. If both a library package and
a command with the same name exists, using the prefix cmd/ will force
documentation on the command rather than the library package. If the -src
flag is specified, anthadoc prints the exported interface of a package in Go/Antha
source form, or the implementation of a specific exported language entity:

	anthadoc fmt                # documentation for package fmt
	anthadoc fmt Printf         # documentation for fmt.Printf
	anthadoc cmd/go             # force documentation for the antha command
	anthadoc -src fmt           # fmt package interface in Antha source form
	anthadoc -src fmt Printf    # implementation of fmt.Printf

In command-line mode, the -q flag enables search queries against a anthadoc running
as a webserver. If no explicit server address is specified with the -server flag,
anthadoc first tries localhost:6060 and then http://golang.org.

	anthadoc -q Reader
	anthadoc -q math.Sin
	anthadoc -server=:6060 -q sin

With the -http flag, it runs as a web server and presents the documentation as a
web page.

	anthadoc -http=:6060

Usage:
	anthadoc [flag] package [name ...]

The flags are:
	-v
		verbose mode
	-q
		arguments are considered search queries: a legal query is a
		single identifier (such as ToLower) or a qualified identifier
		(such as math.Sin)
	-src
		print (exported) source in command-line mode
	-tabwidth=4
		width of tabs in units of spaces
	-timestamps=true
		show timestamps with directory listings
	-index
		enable identifier and full text search index
		(no search box is shown if -index is not set)
	-index_files=""
		glob pattern specifying index files; if not empty,
		the index is read from these files in sorted order
	-index_throttle=0.75
		index throttle value; a value of 0 means no time is allocated
		to the indexer (the indexer will never finish), a value of 1.0
		means that index creation is running at full throttle (other
		goroutines may get no time while the index is built)
	-links=true:
		link identifiers to their declarations
	-write_index=false
		write index to a file; the file name must be specified with
		-index_files
	-maxresults=10000
		maximum number of full text search results shown
		(no full text index is built if maxresults <= 0)
	-notes="BUG"
		regular expression matching note markers to show
		(e.g., "BUG|TODO", ".*")
	-html
		print HTML in command-line mode
	-goroot=$GOROOT
		Go root directory
	-http=addr
		HTTP service address (e.g., '127.0.0.1:6060' or just ':6060')
	-server=addr
		webserver address for command line searches
	-analysis=type,pointer
		comma-separated list of analyses to perform
    		"type": display identifier resolution, type info, method sets,
			'implements', and static callees
		"pointer" display channel peers, callers and dynamic callees
			(significantly slower)
		See http://golang.org/lib/anthadoc/analysis/help.html for details.
	-templates=""
		directory containing alternate template files; if set,
		the directory may provide alternative template files
		for the files in $GOROOT/lib/anthadoc
	-url=path
		print to standard output the data that would be served by
		an HTTP request for path
	-zip=""
		zip file providing the file system to serve; disabled if empty

By default, anthadoc looks at the packages it finds via $GOROOT and $GOPATH (if set).
This behavior can be altered by providing an alternative $GOROOT with the -goroot
flag.

When anthadoc runs as a web server and -index is set, a search index is maintained.
The index is created at startup.

The index contains both identifier and full text search information (searchable
via regular expressions). The maximum number of full text search results shown
can be set with the -maxresults flag; if set to 0, no full text results are
shown, and only an identifier index but no full text search index is created.

The presentation mode of web pages served by anthadoc can be controlled with the
"m" URL parameter; it accepts a comma-separated list of flag names as value:

	all	show documentation for all declarations, not just the exported ones
	methods	show all embedded methods, not just those of unexported anonymous fields
	src	show the original source code rather then the extracted documentation
	text	present the page in textual (command-line) form rather than HTML
	flat	present flat (not indented) directory listings using full paths

For instance, http://golang.org/pkg/math/big/?m=all,text shows the documentation
for all (not just the exported) declarations of package big, in textual form (as
it would appear when using anthadoc from the command line: "anthadoc -src math/big .*").

By default, anthadoc serves files from the file system of the underlying OS.
Instead, a .zip file may be provided via the -zip flag, which contains
the file system to serve. The file paths stored in the .zip file must use
slash ('/') as path separator; and they must be unrooted. $GOROOT (or -goroot)
must be set to the .zip file directory path containing the Go root directory.
For instance, for a .zip file created by the command:

	zip go.zip $HOME/go

one may run anthadoc as follows:

	anthadoc -http=:6060 -zip=go.zip -goroot=$HOME/go

Anthadoc documentation is converted to HTML or to text using the antha/doc package;
see http://golang.org/pkg/go/doc/#ToHTML for the exact rules.
See "Godoc: documenting Go code" for how to write good comments for anthadoc:
http://golang.org/doc/articles/anthadoc_documenting_go_code.html

*/
package main
