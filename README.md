# slf-corpora
manage user corpora (get, set, list, write json)

### methods
- `TryGetCorpusName(uuid uint64) -> (corpusName string, found bool)`
- `GetCorpusName(uuid uint64) -> (corpusName string)`
- `PutCorpusName(uuid uint64, corpusName string) -> error`
- `ListCorpora() -> corporaNames []string`
- `ListNgramNames() -> ngramNames []string`
- `WriteJSON()`
- `GetNgrams(n int, uuid uint64) -> corpus []Ngram{chars: []rune, freq: float64}`
