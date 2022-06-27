package cli

import (
	"aurox_task/internal/app"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

func Run() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	flag.CommandLine.SetOutput(os.Stdout)

	url := flag.String("url", "", "url of the site to generate the sitemap(required)")
	maxDepth := flag.Int("max-depth", 2, "max depth of url navigation recursion")
	parallel := flag.Int("parallel", 1, "number of parallel workers to navigate through site")
	outputFile := flag.String("outputFile", "sitemap.xml", "output file path")
	help := flag.Bool("help", false, "show options")
	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if *url == "" {
		log.Fatal("url argument is mandatory")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	settings := map[string]interface{}{
		"Url":        url,
		"Parallel":   parallel,
		"MaxDepth":   maxDepth,
		"OutputFile": outputFile,
	}

	app.Run(ctx, settings)
}
