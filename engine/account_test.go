package engine

import "testing"

func TestDomainFromUrl(t *testing.T) {
	for _, tc := range []struct {
		url    string
		domain string
	}{
		{"https://github.com/vqhuy/luxio", "github.com"},
		{"https://gist.github.com/", "github.com"},
		{"github.com", "github.com"},
		{"my.special.service", "my.special.service"},
	} {
		if got := DomainFromUrl(tc.url); got != tc.domain {
			t.Error("Expect", tc.domain,
				"got", got)
		}
	}
}

func TestSerialize(t *testing.T) {
	acc := &AccountInfo{
		Domain:   "random.org",
		Username: "tester",
	}
	out := acc.Serialize()
	if t1, t2 := acc.Deserialize(out); t1 != acc.Domain || t2 != acc.Username {
		t.Error("bad format")
	}
}
