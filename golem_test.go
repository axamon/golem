package golem_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aaaton/golem"

	itasset "github.com/axamon/golem/dicts/IT"
)

var en *golem.Lemmatizer
var fr *golem.Lemmatizer
var it *golem.Lemmatizer
var ge *golem.Lemmatizer
var sp *golem.Lemmatizer
var sw *golem.Lemmatizer

// Languages
var languages = []struct {
	lang       string
	lemmatizer *golem.Lemmatizer
}{
	{"english", en},
	//{"italian", it},
	{"french", fr},
	{"german", ge},
	{"spanish", sp},
	{"swedish", sw},
}

func init() {
	var err error
	for _, language := range languages {
		language.lemmatizer, err = golem.New(language.lang)
		if err != nil {
			log.Fatal("Error in init function: ", err.Error())
		}
	}
}

var flagtests = []struct {
	name       string
	language   string
	lemmatizer *golem.Lemmatizer
	in         string
	out        string
}{
	{"Italian Verb", "italian", it, "lavorerai", "lavorare"},
	{"Italian Plural Noun", "italian", it, "bicchieri", "bicchiere"},
	{"Italian FirstName", "italian", it, "Alberto", "Alberto"},
	{"Italian Plural Adjective", "italian", it, "lunghi", "lungo"},
	{"Swedish Example1", "swedish", sw, "Avtalet", "avtal"},
	{"Swedish Example2", "swedish", sw, "avtalets", "avtal"},
	{"Swedish Example3", "swedish", sw, "avtalens", "avtal"},
	{"Swedish Example4", "swedish", sw, "Avtaletsadlkj", "Avtaletsadlkj"},
	{"English Verb", "english", en, "goes", "go"},
	{"English Noun", "english", en, "wolves", "wolf"},
	{"English FirstName", "english", en, "Edward", "Edward"},
	{"French Example1", "french", fr, "avait", "avoir"},
	{"Spanish Example1", "spanish", sp, "Buenas", "bueno"},
	{"German Example1", "german", ge, "Hast", "haben"},
}

func TestLemmatizer_Lemma_All(t *testing.T) {
	for _, tt := range flagtests {
		t.Run(tt.in, func(t *testing.T) {
			l := tt.lemmatizer
			got := l.Lemma(tt.in)
			if got != tt.out {
				t.Errorf("%s Lemmatizer.Lemma() = %v, want %v", tt.name, got, tt.out)
			}
			got = l.LemmaLower(strings.ToLower(tt.in))
			if got != strings.ToLower(tt.out) {
				t.Errorf("%s Lemmatizer.LemmaLower() = %v, want %v", tt.name, got, tt.out)
			}

		})
	}
}

func TestReadBinary_IT(t *testing.T) {
	b, err := itasset.Asset("data/it.gz")
	if err != nil {
		t.Fatal(err)
	}
	_, err = gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}
}

var exampleDataLemma = []struct {
	language string
	word     string
}{
	{"english", "agreed"},
	{"italian", "armadi"},
	{"swedish", "Avtalet"},
}

func ExampleLemmatizer_Lemma() {
	for _, element := range exampleDataLemma {
		l, err := golem.New(element.language)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(l.Lemma(element.word))
	}
	// Output:
	// agree
	// armadio
	// avtal
}

var exampleDataInDict = []struct {
	language string
	word     string
	result   bool
}{
	{"italian", "armadio", true},
	{"italian", "ammaccabanane", false},
	{"swedish", "Avtalet", true},
	{"swedish", "Avtalt", false},
}

func ExampleLemmatizer_InDict() {
	for _, element := range exampleDataInDict {
		l, err := golem.New(element.language)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(l.InDict(element.word))
	}
	// Output:
	// true
	// false
	// true
	// false
}

var exampleDataLemmas = []struct {
	language string
	word     string
	result   []string
}{
	{"italian", "soli", []string{"sole", "solo"}},
}

func ExampleLemmatizer_Lemmas() {
	for _, element := range exampleDataLemmas {
		l, err := golem.New(element.language)
		if err != nil {
			log.Fatal(err)
		}
		lemmas := l.Lemmas(element.word)
		for _, lemma := range lemmas {
			fmt.Println(lemma)
		}
	}
	// Unordered output:
	// solare
	// solere
	// solo
	// sole
}

func BenchmarkLookup(b *testing.B) {
	l, err := golem.New("swedish")
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N/2; i++ {
		l.Lemma("Avtalet")
	}
}

func BenchmarkLookupLower(b *testing.B) {
	l, err := golem.New("swedish")
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N/2; i++ {
		l.LemmaLower("avtalet")
	}
}
