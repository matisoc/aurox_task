package modules

import (
	"aurox_task/internal/modules/smg"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

func NewSiteMapGenerator(ctx context.Context, settings map[string]interface{}) {
	txt := fmt.Sprintf("%v", "NewSiteMapGenerator()")
	fmt.Println(txt)
	jsonStr, err := json.Marshal(settings)
	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct
	var g SiteMapGenerator
	if err := json.Unmarshal(jsonStr, &g); err != nil {
		fmt.Println(err)
	}
	err = g.Execute(ctx)
	if err != nil {
		fmt.Println(err)
	}
}

type SiteMapGenerator struct {
	Url        string
	Parallel   int
	MaxDepth   int
	OutputFile string
}

func (sm *SiteMapGenerator) Execute(ctx context.Context) error {
	url, err := url.Parse(sm.Url)
	if err != nil {
		return fmt.Errorf("error parsing url: %v\n", err)
	}
	c := smg.NewCoordinator(ctx, url, sm.MaxDepth, sm.Parallel)
	c.Start(ctx)
	data := c.Data()
	return sm.generateSiteMap(data)
}

func (sm *SiteMapGenerator) generateSiteMap(data *smg.DataResponse) error {
	err := smg.BuildSitemapFile(sm.OutputFile, data.UniqueURLs)
	return err
}
