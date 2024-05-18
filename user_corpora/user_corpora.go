package user_corpora

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/waterdragen/slf-corpora/v2/assert"
)

const DefaultCorpus string = "shai"

func TryGetCorpusName(uuid uint64) (string, bool) {
	return _UserCorporaData.TryGet(uuid)
}

func GetCorpusName(uuid uint64) string {
	return _UserCorporaData.Get(uuid)
}

func PutCorpusName(uuid uint64, corpusName string) error {
	return _UserCorporaData.Put(uuid, corpusName)
}

func ListCorpora() []string {
	corporaNames := make([]string, len(_CorporaNames))
	copy(corporaNames, _CorporaNames)
	return corporaNames
}

func ListNgramNames() []string {
	ngramNames := make([]string, len(_NgramNames))
	copy(ngramNames, _NgramNames)
	return ngramNames
}

func WriteJSON() {
	_UserCorporaData.WriteJSON()
}

func GetNgrams(n int, uuid uint64) []Ngram {
	corpusName := _UserCorporaData.Get(uuid)
	ngramName := _NgramNames[n-1]
	path := fmt.Sprintf("./corpora/%s/%s.json", corpusName, ngramName)
	corpus := _CorporaData.loadCorpus(path)
	return corpus
}

var _NgramNames = []string{"monograms", "bigrams", "trigrams"}
var _CorporaNames = listCorporaOnce()
var _UserCorporaData = loadUserCorporaOnce()
var _CorporaData = newCorpora()

func listCorporaOnce() []string {
	corporaPath := "./corpora/"
	entries, err := os.ReadDir(corporaPath)
	assert.Ok(err)
	var corporaNames []string
	for _, entry := range entries {
		corporaNames = append(corporaNames, entry.Name())
	}
	return corporaNames
}

// UserCorpora
// - A sync map[uuid]corpusName
// - Caches corpora.json
type UserCorpora struct {
	inner map[uint64]string
	lock  sync.RWMutex
}

func loadUserCorporaOnce() UserCorpora {
	userCorpora := make(map[uint64]string)
	data, err := os.ReadFile("./corpora.json")
	assert.Ok(err)
	err = json.Unmarshal(data, &userCorpora)
	assert.Ok(err)

	return UserCorpora{
		inner: userCorpora,
		lock:  sync.RWMutex{},
	}
}

func (uc *UserCorpora) TryGet(uuid uint64) (string, bool) {
	uc.lock.RLock()
	defer uc.lock.RUnlock()
	corpusName, found := uc.inner[uuid]
	return corpusName, found
}

func (uc *UserCorpora) Get(uuid uint64) string {
	corpusName, found := uc.TryGet(uuid)
	if !found {
		uc.lock.Lock()
		defer uc.lock.Unlock()
		uc.inner[uuid] = DefaultCorpus
		return DefaultCorpus
	}
	return corpusName
}

func (uc *UserCorpora) Put(uuid uint64, corpusName string) error {
	exists := false
	for _, existedCorpusName := range _CorporaNames {
		if corpusName == existedCorpusName {
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf("the corpus `%s` doesn't exist", corpusName)
	}
	uc.lock.Lock()
	defer uc.lock.Unlock()
	uc.inner[uuid] = corpusName
	return nil
}

func (uc *UserCorpora) WriteJSON() {
	data, err := json.MarshalIndent(uc.inner, "", "    ")
	assert.Ok(err)
	err = os.WriteFile("./corpora.json", data, 0644)
	assert.Ok(err)
}

type Ngram struct {
	chars []rune
	freq  float64
}

// Corpora
// - A sync map[corpusName]Ngrams.
// - Caches the corpora lazily, loads a new corpus if not found.
// - Assuming corpora are in ./corpora/<corpus name>/trigrams.json.
type Corpora struct {
	inner map[string][]Ngram
	lock  sync.RWMutex
}

func newCorpora() Corpora {
	return Corpora{
		inner: make(map[string][]Ngram),
		lock:  sync.RWMutex{},
	}
}

func (c *Corpora) tryLoadCorpus(path string) ([]Ngram, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	corpus, found := c.inner[path]
	return corpus, found
}

func (c *Corpora) putCorpus(path string, ngrams []Ngram) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner[path] = ngrams
}

func (c *Corpora) loadCorpus(path string) []Ngram {
	corpus, found := c.tryLoadCorpus(path)
	if found {
		return corpus
	}
	jsonData, err := os.ReadFile(path)
	assert.Ok(err)

	var rawCorpus map[string]float64
	err = json.Unmarshal(jsonData, &rawCorpus)
	assert.Ok(err)

	corpus = make([]Ngram, 0, len(rawCorpus))
	for chars, freq := range rawCorpus {
		corpus = append(corpus, Ngram{[]rune(chars), freq})
	}
	c.putCorpus(path, corpus)
	return corpus
}
