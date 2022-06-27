package smg

import "net/url"

// processor defines coordinatorPayload process
type processor interface {
	process(c *Coordinator, wd *coordinatorPayload) (proceed bool)
}

type processorFunc func(c *Coordinator, wd *coordinatorPayload) (proceed bool)

func (pf processorFunc) process(c *Coordinator, wd *coordinatorPayload) (proceed bool) {
	return pf(c, wd)
}

// uniqueURLProcessor adds source url to unique crawled
func uniqueURLProcessor() processor {
	return processorFunc(func(c *Coordinator, wd *coordinatorPayload) (proceed bool) {
		c.scrappedUnique[wd.sourceURL.String()]++
		var unique []*url.URL
		for _, u := range wd.urls {
			if _, ok := c.scrappedUnique[u.String()]; !ok {
				unique = append(unique, u)
				continue
			}

			c.scrappedUnique[u.String()]++
		}
		wd.urls = unique
		return true
	})
}

// errorCheckProcessor check if the url scrape failed for any reason
func errorCheckProcessor() processor {
	return processorFunc(func(c *Coordinator, wd *coordinatorPayload) (proceed bool) {
		if wd.err == nil {
			return true
		}

		c.errorURLs[wd.sourceURL.String()] = wd.err
		return false
	})
}

// skippedURLProcessor add the unknown urls to skipped map
func skippedURLProcessor() processor {
	return processorFunc(func(c *Coordinator, wd *coordinatorPayload) (proceed bool) {
		c.skippedURLs[wd.sourceURL.String()] = append(c.skippedURLs[wd.sourceURL.String()], wd.invalidURLs...)
		return true
	})
}

// maxDepthCheckProcessor add the unscrapped urls to scrapped if the max depth has been reached
func maxDepthCheckProcessor() processor {
	return processorFunc(func(c *Coordinator, wd *coordinatorPayload) (proceed bool) {
		if c.maxDepth == -1 || wd.depth < c.maxDepth {
			return true
		}
		if len(wd.urls) < 1 {
			return false
		}
		c.scrapped[wd.depth] = append(c.scrapped[wd.depth], wd.urls...)
		for _, u := range wd.urls {
			if c.domainRegex.MatchString(u.Hostname()) {
				c.scrappedUnique[u.String()]++
				continue
			}
		}
		return false
	})
}

// domainFilterProcessor filter the wd.urls and update skipped urls with unmatched urls
func domainFilterProcessor() processor {
	return processorFunc(func(c *Coordinator, wd *coordinatorPayload) (proceed bool) {
		if c.domainRegex == nil {
			return true
		}
		var m []*url.URL
		var um []string
		for _, u := range wd.urls {
			if c.domainRegex.MatchString(u.Hostname()) {
				m = append(m, u)
				continue
			}
			um = append(um, u.String())
		}
		wd.urls = m
		c.skippedURLs[wd.sourceURL.String()] = append(c.skippedURLs[wd.sourceURL.String()], um...)
		return true
	})
}
