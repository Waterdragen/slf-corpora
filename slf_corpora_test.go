package slf_corpora

import (
	"os"
	"os/signal"
	"sort"
	"syscall"
	"testing"

	"github.com/waterdragen/slf-corpora/v2/assert"
	"github.com/waterdragen/slf-corpora/v2/cron_job"
	"github.com/waterdragen/slf-corpora/v2/user_corpora"
)

func TestCronJob(t *testing.T) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	cronJob := cron_job.NewDailyCronJob(user_corpora.WriteJSON)
	cronJob.Start()

	// Uncomment the following:
	// <-done
	user_corpora.WriteJSON()
}

func TestGetCorpusName(t *testing.T) {
	monkeyracerUser := uint64(1007355784830652507)
	corpusName := user_corpora.GetCorpusName(monkeyracerUser)
	assert.Eq(corpusName, "monkeyracer")
}

func TestTryGetCorpusName(t *testing.T) {
	shaiUser := uint64(1010618301895934034)
	corpusName, found := user_corpora.TryGetCorpusName(shaiUser)
	assert.Eq(corpusName, "shai")
	assert.Eq(found, true)

	notAUser := uint64(69420)
	corpusName, found = user_corpora.TryGetCorpusName(notAUser)
	assert.Eq(corpusName, "")
	assert.Eq(found, false)
}

func TestGetNgrams(t *testing.T) {
	shaiUser := uint64(1010618301895934034)
	shaiLength := 163381
	ngrams := user_corpora.GetNgrams(3, shaiUser)
	assert.Eq(len(ngrams), shaiLength)
}

func TestPutCorpusName(t *testing.T) {
	newUser := uint64(69420)
	defaultCorpusName := user_corpora.GetCorpusName(newUser)
	assert.Eq(defaultCorpusName, "shai")

	validCorpusName := "monkeyracer"
	err := user_corpora.PutCorpusName(newUser, validCorpusName)
	assert.Ok(err)

	invalidCorpusName := "amogus"
	err = user_corpora.PutCorpusName(newUser, invalidCorpusName)
	assert.Ne(err, nil)

	corpusName := user_corpora.GetCorpusName(newUser)
	assert.Eq(corpusName, validCorpusName)
}

func TestListCorpora(t *testing.T) {
	corpusNames := user_corpora.ListCorpora()
	sort.Strings(corpusNames)
	assert.Eq(len(corpusNames), 18)
	expectedCorpusNames := []string{
		"akl",
		"chained-bigrams",
		"dutch",
		"english-10k", "english-1k", "english-200", "english-450k", "english-5k",
		"french",
		"gutenberg",
		"keymash",
		"monkeyracer", "movie", "mt-quotes",
		"reddit",
		"shai",
		"tr-quotes",
		"wordle"}
	for index, _ := range expectedCorpusNames {
		assert.Eq(corpusNames[index], expectedCorpusNames[index])
	}
}

func TestListNgramNames(t *testing.T) {
	ngramNames := user_corpora.ListNgramNames()
	sort.Strings(ngramNames)
	assert.Eq(len(ngramNames), 3)
	expectedNgramNames := []string{"bigrams", "monograms", "trigrams"}
	for index, _ := range expectedNgramNames {
		assert.Eq(ngramNames[index], expectedNgramNames[index])
	}
}
