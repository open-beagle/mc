// Copyright (c) 2015-2024 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"strings"

	"github.com/minio/cli"
	"github.com/minio/pkg/v2/console"
)

// get command flags.
var (
	getFlags = []cli.Flag{}
)

// Get command.
var getCmd = cli.Command{
	Name:         "get",
	Usage:        "get s3 object to local",
	Action:       mainGet,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        append(append(ioFlags, globalFlags...), getFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS] SOURCE TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
ENVIRONMENT VARIABLES:
  MC_ENCRYPT:      list of comma delimited prefixes
  MC_ENCRYPT_KEY:  list of comma delimited prefix=secret values

EXAMPLES:
  1. Get an object from S3 storage to local file system 
    {{.Prompt}} {{.HelpName}} ALIAS/BUCKET/object path-to/object 
`,
}

// mainGet is the entry point for get command.
func mainGet(cliCtx *cli.Context) (e error) {
	ctx, cancelGet := context.WithCancel(globalContext)
	defer cancelGet()

	encKeyDB, err := getEncKeys(cliCtx)
	fatalIf(err, "Unable to parse encryption keys.")

	args := cliCtx.Args()
	if len(args) != 2 {
		showCommandHelpAndExit(cliCtx, 1) // last argument is exit code.
	}

	// get source and target
	sourceURLs := args[:len(args)-1]
	targetURL := args[len(args)-1]

	getURLsCh := make(chan URLs, 10000)
	var totalObjects, totalBytes int64

	// Store a progress bar or an accounter
	var pg ProgressReader
	// Enable progress bar reader only during default mode.
	if !globalQuiet && !globalJSON { // set up progress bar
		pg = newProgressBar(totalBytes)
	} else {
		pg = newAccounter(totalBytes)
	}
	go func() {
		opts := prepareCopyURLsOpts{
			sourceURLs:              sourceURLs,
			targetURL:               targetURL,
			encKeyDB:                encKeyDB,
			ignoreBucketExistsCheck: true,
		}

		for getURLs := range prepareGetURLs(ctx, opts) {
			if getURLs.Error != nil {
				getURLsCh <- getURLs
				break
			}
			totalObjects++
			getURLsCh <- getURLs
		}
		close(getURLsCh)
	}()
	for {
		select {
		case <-ctx.Done():
			showLastProgressBar(pg, nil)
			return
		case getURLs, ok := <-getURLsCh:
			if !ok {
				showLastProgressBar(pg, nil)
				return
			}
			if getURLs.Error != nil {
				printGetURLsError(&getURLs)
				showLastProgressBar(pg, getURLs.Error.ToGoError())
				return
			}
			urls := doCopy(ctx, doCopyOpts{
				cpURLs:              getURLs,
				pg:                  pg,
				encKeyDB:            encKeyDB,
				updateProgressTotal: true,
			})
			if urls.Error != nil {
				e = urls.Error.ToGoError()
				showLastProgressBar(pg, e)
				return
			}
		}
	}
}

func printGetURLsError(cpURLs *URLs) {
	// Print in new line and adjust to top so that we
	// don't print over the ongoing scan bar
	if !globalQuiet && !globalJSON {
		console.Eraseline()
	}

	if strings.Contains(cpURLs.Error.ToGoError().Error(),
		" is a folder.") {
		errorIf(cpURLs.Error.Trace(),
			"Folder cannot be copied. Please use `...` suffix.")
	} else {
		errorIf(cpURLs.Error.Trace(),
			"Unable to download.")
	}
}
