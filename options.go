package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alsm/ioprogress"
	"github.com/davidjpeacock/cli"
)

type Options struct {
	outputFilename string
	fileUpload     string
	remoteName     bool
	verbose        bool
	maxTime        uint
	remoteTime     bool
	followRedirect bool
	maxRedirects   uint
	redirectsTaken uint
	silent         bool
	method         string
	headers        []string
	agent          string
	user           string
	expectTimeout  uint
	data           []string
	dataAscii      []string
	dataRaw        []string
	dataBinary     []string
	dataURLEncode  []string
	head           bool
}

func (o *Options) getOptions(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "Filename to name url content to",
			Destination: &o.outputFilename,
		},
		cli.StringFlag{
			Name:        "upload-file, T",
			Usage:       "File to upload",
			Destination: &o.fileUpload,
		},
		cli.BoolFlag{
			Name:        "remote-name, O",
			Usage:       "Save output to file named with file part of URL",
			Destination: &o.remoteName,
		},
		cli.BoolFlag{
			Name:        "v",
			Usage:       "Verbose output",
			Destination: &o.verbose,
		},
		cli.UintFlag{
			Name:        "max-time, m",
			Usage:       "Maximum time to wait for an operation to complete in seconds",
			Destination: &o.maxTime,
		},
		cli.BoolFlag{
			Name:        "R",
			Usage:       "Set the timestamp of the local file to that of the remote file, if available",
			Destination: &o.remoteTime,
		},
		cli.BoolFlag{
			Name:        "location, L",
			Usage:       "Follow 3xx redirects",
			Destination: &o.followRedirect,
		},
		cli.UintFlag{
			Name:        "max-redirs",
			Usage:       "Maximum number of 3xx redirects to follow",
			Destination: &o.maxRedirects,
			Value:       10,
		},
		cli.BoolFlag{
			Name:        "silent, s",
			Usage:       "Mute kurly entirely, operation without any output",
			Destination: &o.silent,
		},
		cli.StringFlag{
			Name:        "request, X",
			Usage:       "HTTP method to use",
			Destination: &o.method,
			Value:       "GET",
		},
		cli.StringFlag{
			Name:        "user-agent, A",
			Usage:       "User agent to set for this request",
			Destination: &o.agent,
			Value:       "Kurly/1.0",
		},
		cli.StringFlag{
			Name:        "user, u",
			Usage:       "User authentication data to set for this request",
			Destination: &o.user,
		},
		cli.StringSliceFlag{
			Name:  "header, H",
			Usage: "Extra headers to be sent with the request",
		},
		cli.UintFlag{
			Name:        "expect100-timeout",
			Usage:       "Timeout in seconds for Expect: 100-continue wait period",
			Destination: &o.expectTimeout,
			Value:       1,
		},
		cli.StringSliceFlag{
			Name:  "data, d",
			Usage: "Sends the specified data in a POST request to the server",
		},
		cli.StringSliceFlag{
			Name:  "data-ascii",
			Usage: "The same as --data, -d",
		},
		cli.StringSliceFlag{
			Name:  "data-raw",
			Usage: "Basically the same as --data-binary (no @ interpretation)",
		},
		cli.StringSliceFlag{
			Name:  "data-binary",
			Usage: "Sends the data as binary",
		},
		cli.StringSliceFlag{
			Name:  "data-urlencode",
			Usage: "Sends the data as urlencoded ascii",
		},
		cli.BoolFlag{
			Name:        "head, I",
			Usage:       "Get HEAD from URL only",
			Destination: &o.head,
		},
	}
}

func (o *Options) checkRedirect(req *http.Request, via []*http.Request) error {
	o.redirectsTaken++

	if !o.followRedirect || o.redirectsTaken >= o.maxRedirects {
		return http.ErrUseLastResponse
	}

	return nil
}

func (o *Options) ProcessData() {
	var uriEncodes url.Values
	for _, d := range o.dataAscii {
		parts := strings.SplitN(d, "=", 2)
		if len(parts) == 1 {
			o.data = append(o.data, d)
			continue
		}
		if strings.HasPrefix(parts[1], "@") {
			data, err := ioutil.ReadFile(strings.TrimPrefix(parts[1], "@"))
			if err != nil {
				Status.Fatalf("Unable to read file %s for data element %s\n", strings.TrimPrefix(parts[1], "@"), parts[0])
			}
			data = []byte(strings.Replace(string(data), "\r", "", -1))
			data = []byte(strings.Replace(string(data), "\n", "", -1))
			o.data = append(o.data, fmt.Sprintf("%s=%s", parts[0], string(data)))
		} else {
			o.data = append(o.data, d)
		}
	}
	for _, d := range o.dataRaw {
		parts := strings.SplitN(d, "=", 2)
		o.data = append(o.data, fmt.Sprintf("%s=%s", parts[0], parts[1]))
	}
	for _, d := range o.dataBinary {
		parts := strings.SplitN(d, "=", 2)
		o.data = append(o.data, fmt.Sprintf("%s=%s", parts[0], parts[1]))
	}
	for _, d := range o.dataURLEncode {
		parts := strings.SplitN(d, "=", 2)
		uriEncodes.Add(parts[0], parts[1])
	}
	if len(uriEncodes) > 0 {
		o.data = append(o.data, uriEncodes.Encode())
	}
}

func (o *Options) openOutputFile() *os.File {
	var err error
	var outputFile *os.File
	if o.outputFilename != "" {
		if outputFile, err = os.Create(o.outputFilename); err != nil {
			Status.Fatalf("Error: Unable to create file '%s' for output\n", o.outputFilename)
		}
	} else {
		outputFile = os.Stdout
	}
	return outputFile
}

func (o *Options) uploadFile() {
	o.method = "PUT"

	tr := &http.Transport{
		ExpectContinueTimeout: time.Duration(o.expectTimeout) * time.Second,
	}
	client.Transport = tr
	o.headers = append(o.headers, "Expect: 100-continue")

	reader, err := os.Open(o.fileUpload)
	if err != nil {
		Status.Fatalf("Error opening %s\n", o.fileUpload)
	}

	if !o.silent {
		fi, err := reader.Stat()
		if err != nil {
			Status.Fatalf("Unable to get file stats for %v\n", o.fileUpload)
		}
		body = &ioprogress.Reader{
			Reader: reader,
			Size:   fi.Size(),
			DrawFunc: ioprogress.DrawTerminalf(os.Stderr, func(progress, total int64) string {
				return fmt.Sprintf(
					"%s %s",
					(ioprogress.DrawTextFormatBarWithIndicator(40, '>'))(progress, total),
					ioprogress.DrawTextFormatBytes(progress, total))
			}),
		}
	} else {
		body = reader
	}
}
