package smg

import (
	"net/url"
	"testing"
)

func Test_digSite(t *testing.T) {
	data := []struct {
		url   string
		depth int
		error bool
	}{
		{
			url:   "https://lavoz.com.ar",
			depth: 0,
		},

		{
			url:   "https://familia.edu.gva.es/wf-front/myitaca/",
			depth: 0,
			error: true,
		},
	}

	for _, test := range data {
		u, _ := url.Parse(test.url)
		wd := digSite(test.depth, u)
		if wd.err != nil && !test.error {
			t.Fatalf("failed to dig %s\n", u.String())
		}

		if u.String() != wd.sourceURL.String() {
			t.Fatalf("unknown source url found: %s\n", wd.sourceURL.String())
		}

		if wd.depth != test.depth+1 {
			t.Fatalf("expected depth %d but got %d\n", test.depth+1, wd.depth)
		}
	}
}
