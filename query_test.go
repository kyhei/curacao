package curacao

import (
	"testing"
)

func TestQuery(t *testing.T) {
	q := map[string][]string{
		"q1": []string{"hoge", "fuga", "moja"},
		"q2": []string{"fuga", "moja", "hoge"},
		"q3": []string{"moja", "hoge", "fuga"},
	}

	httpQuery := NewHTTPQuery(q)

	q1, _ := httpQuery.Get("q1")
	if q1 != "hoge" {
		t.Errorf("q1 is expected 'hoge', but get: %v", q1)
	}

	q2, _ := httpQuery.Get("q2")
	if q2 != "fuga" {
		t.Errorf("q2 is expected 'fuga', but get: %v", q2)
	}

	q3, _ := httpQuery.Get("q3")
	if q3 != "moja" {
		t.Errorf("q3 is expected 'moja', but get: %v", q3)
	}

	_, err := httpQuery.Get("q4")
	if err == nil {
		t.Error("q4 query does not exsist. But found it")
	}

}
