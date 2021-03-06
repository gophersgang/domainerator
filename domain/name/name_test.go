package name

import (
	"reflect"
	"sort"
	"testing"

	"github.com/hgfischer/domainerator/tests"
)

var (
	accepted = map[string]bool{
		"com": true, "net": true, "org": true,
		"us": true, "im": true, "io": true,
		"ca": true, "co": true, "in": true,
		"com.br": true, "org.br": true, "co.uk": true,
	}
)

func TestParsePublicSuffixCSV(t *testing.T) {
	csv := "us, im, in, io, ,, ca,co,co ,,ca , com"
	expected := []string{"com", "net", "org", "us", "im", "in", "io", "ca", "co"}
	sort.Strings(expected)
	psl, err := ParsePublicSuffixCSV(csv, accepted, true)
	if err != nil {
		t.Fatalf(tests.ErrFmtExpectedGot, "ParsePublicSuffixCSV", "No Error", err)
	}
	if !reflect.DeepEqual(psl, expected) {
		t.Errorf(tests.ErrFmtExpectedGot, "ParsePublicSuffixCSV", expected, psl)
	}
}

func TestParsePublicSuffixCSVForUnknownSuffix(t *testing.T) {
	csv := "com,net,org,unk"
	_, err := ParsePublicSuffixCSV(csv, accepted, false)
	if err == nil {
		t.Errorf(tests.ErrFmtExpectedGot, "ParsePublicSuffixCSV", "Unknown Public Suffix Error", "No Error")
	}
}

func TestCombinePhraseAndPublicSuffixes(t *testing.T) {
	psl := []string{"ex", "nd", "com"}
	expected := []string{"index.ex", "ind.ex", "index.nd", "index.com"}
	domains := CombinePhraseAndPublicSuffixes("index", psl, true)
	if !reflect.DeepEqual(expected, domains) {
		t.Errorf(tests.ErrFmtExpectedGot, "CombineWordAndPublixSuffixes", expected, domains)
	}
}

func TestCombinePhraseAndPublicSuffixesWithSmallPhrase(t *testing.T) {
	psl := []string{"ex"}
	expected := []string{"ex.ex"}
	domains := CombinePhraseAndPublicSuffixes("ex", psl, true)
	if !reflect.DeepEqual(expected, domains) {
		t.Errorf(tests.ErrFmtExpectedGot, "CombineWordAndPublixSuffixes", expected, domains)
	}
}

var (
	prefixes = []string{"go", "py"}
	suffixes = []string{"lang", "coder"}
	psl      = []string{"com", "er"}
)

func TestCombineFull(t *testing.T) {
	expected := []string{
		"go.com", "go.er", "lang.com", "lang.er",
		"golang.com", "golang.er", "go-lang.com", "go-lang.er",
		"gocoder.com", "gocoder.er", "gocod.er", "go-coder.com", "go-coder.er", "go-cod.er",
		"py.com", "py.er", "coder.com", "coder.er", "cod.er",
		"pylang.com", "pylang.er", "py-lang.com", "py-lang.er",
		"pycoder.com", "pycoder.er", "pycod.er", "py-coder.com", "py-coder.er", "py-cod.er",
	}
	sort.Strings(expected)
	words := Combine(prefixes, suffixes, psl, true, true, true, true, false, 3)
	sort.Strings(words)
	if !reflect.DeepEqual(expected, words) {
		t.Errorf(tests.ErrFmtExpectedGot, "Combine", expected, words)
	}
}

func TestCombineSimple(t *testing.T) {
	expected := []string{
		"golang.com", "golang.er", "gocoder.com", "gocoder.er",
		"pylang.com", "pylang.er", "pycoder.com", "pycoder.er",
	}
	sort.Strings(expected)
	words := Combine(prefixes, suffixes, psl, false, false, false, false, false, 3)
	sort.Strings(words)
	if !reflect.DeepEqual(expected, words) {
		t.Errorf(tests.ErrFmtExpectedGot, "Combine", expected, words)
	}
}

func TestCombineHyphenation(t *testing.T) {
	expected := []string{
		"golang.com", "golang.er", "gocoder.com", "gocoder.er",
		"go-lang.com", "go-lang.er", "go-coder.com", "go-coder.er",
		"pylang.com", "pylang.er", "pycoder.com", "pycoder.er",
		"py-lang.com", "py-lang.er", "py-coder.com", "py-coder.er",
	}
	sort.Strings(expected)
	words := Combine(prefixes, suffixes, psl, false, true, false, false, false, 3)
	sort.Strings(words)
	if !reflect.DeepEqual(expected, words) {
		t.Errorf(tests.ErrFmtExpectedGot, "Combine", expected, words)
	}
}

func TestCombinePrefixAndSuffix(t *testing.T) {
	expected := []string{"prefixsuffix", "prefix-suffix"}
	sort.Strings(expected)
	words := CombinePrefixAndSuffix("prefix", "suffix", false, true, false, 3)
	sort.Strings(words)
	if !reflect.DeepEqual(expected, words) {
		t.Errorf(tests.ErrFmtExpectedGot, "CombinePrefixAndSuffix", expected, words)
	}
}

func TestCombinePrefixAndSuffixWithoutHyphenation(t *testing.T) {
	expected := []string{"prefixsuffix"}
	words := CombinePrefixAndSuffix("prefix", "suffix", false, false, false, 3)
	if !reflect.DeepEqual(expected, words) {
		t.Errorf(tests.ErrFmtExpectedGot, "CombinePrefixAndSuffix", expected, words)
	}
}

func TestCombinePrefixAndSuffixWithItself(t *testing.T) {
	expected := []string{"itselfitself", "itself-itself"}
	sort.Strings(expected)
	words := CombinePrefixAndSuffix("itself", "itself", true, true, false, 3)
	sort.Strings(words)
	if !reflect.DeepEqual(expected, words) {
		t.Errorf(tests.ErrFmtExpectedGot, "CombinePrefixAndSuffix", expected, words)
	}
}

func TestCombinePrefixAndSuffixWithFusion(t *testing.T) {
	expected := []string{"fusionnation", "fusionation"}
	sort.Strings(expected)
	words := CombinePrefixAndSuffix("fusion", "nation", false, false, true, 3)
	sort.Strings(words)
	if !reflect.DeepEqual(expected, words) {
		t.Errorf(tests.ErrFmtExpectedGot, "CombinePrefixAndSuffix", expected, words)
	}
}

func TestCombinePrefixAndSuffixWithFusionLatestTwo(t *testing.T) {
	expected := []string{"flamingogorilla", "flamingorilla"}
	sort.Strings(expected)
	words := CombinePrefixAndSuffix("flamingo", "gorilla", false, false, true, 3)
	sort.Strings(words)
	if !reflect.DeepEqual(expected, words) {
		t.Errorf(tests.ErrFmtExpectedGot, "CombinePrefixAndSuffix", expected, words)
	}
}

func TestCombinePrefixAndSuffixWithMinLength(t *testing.T) {
	expected := []string{"a-b"}
	sort.Strings(expected)
	words := CombinePrefixAndSuffix("a", "b", true, true, false, 3)
	sort.Strings(words)
	if !reflect.DeepEqual(expected, words) {
		t.Errorf(tests.ErrFmtExpectedGot, "CombinePrefixAndSuffix", expected, words)
	}
}

func TestFilterStrictDomains(t *testing.T) {
	expected := []string{"lalalala.com"}
	domains := []string{"co.com.br", "us.com.br", "lalalala.com"}
	domains = FilterStrictDomains(domains, accepted)
	if !reflect.DeepEqual(expected, domains) {
		t.Errorf(tests.ErrFmtExpectedGot, "FilterStrictDomains", expected, domains)
	}
}
