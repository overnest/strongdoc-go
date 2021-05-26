package docidx

import (
	"io"
	"math/rand"
	"sort"
	"testing"

	"github.com/overnest/strongdoc-go-sdk/client"
	"github.com/overnest/strongdoc-go-sdk/search/index/docidx/docidxv1"
	"github.com/overnest/strongdoc-go-sdk/utils"
	sscrypto "github.com/overnest/strongsalt-crypto-go"
	"gotest.tools/assert"
)

func TestTools(t *testing.T) {
	// ================================ Prev Test ================================
	sdc := prevTest(t)
	defer CleanupTemporaryDocumentIndex()
	// ================================ Generate Doc Indexes ================================
	key, err := sscrypto.GenerateKey(sscrypto.Type_XChaCha20)
	assert.NilError(t, err)

	docs, err := InitTestDocuments(10, false)
	assert.NilError(t, err)

	// Create the indexes
	for _, doc := range docs {
		assert.NilError(t, doc.CreateDoiAndDti(sdc, key))
	}
	// ================================ Validate Doc Indexes ================================
	for _, doc := range docs {
		terms := validateDocDoi(t, sdc, key, doc)
		validateDocDti(t, sdc, key, doc, terms)
	}
	// ================================ Remove Doc Indexes ================================
	for _, doc := range docs {
		assert.NilError(t, doc.RemoveAllVersionsIndexes(sdc))
	}
}

// validate doi, return sorted term list
func validateDocDoi(t *testing.T, sdc client.StrongDocClient, key *sscrypto.StrongSaltKey, doc *TestDocumentIdxV1) []string {
	// Open the DOI
	doiVersion, err := doc.OpenDoi(sdc, key)
	assert.NilError(t, err)

	doi, ok := doiVersion.(*docidxv1.DocOffsetIdxV1)
	assert.Check(t, ok == true)

	termloc := make(map[string][]uint64)
	for err == nil {
		var blk *docidxv1.DocOffsetIdxBlkV1
		blk, err = doi.ReadNextBlock()
		if err != nil {
			assert.Equal(t, err, io.EOF)
		}
		if blk != nil {
			for term, locs := range blk.TermLoc {
				termloc[term] = append(termloc[term], locs...)
			}
		}
	}

	sourceData, err := utils.OpenLocalFile(doc.DocFilePath)
	assert.NilError(t, err)
	tokenizer, err := utils.OpenRawFileTokenizer(sourceData)
	assert.NilError(t, err)
	err = doiVersion.Close()
	assert.NilError(t, err)

	// Validate the DOI
	i := uint64(0)
	for token, _, wordCounter, err := tokenizer.NextToken(); err != io.EOF; token, _, wordCounter, err = tokenizer.NextToken() {
		if err != nil {
			assert.Equal(t, err, io.EOF)
		}

		if token != "" {
			locs, exist := termloc[token]
			assert.Assert(t, exist)
			assert.Assert(t, len(locs) > 0)
			assert.Equal(t, i, wordCounter)
			termloc[token] = locs[1:]
			i++
		}
	}
	err = sourceData.Close()
	assert.NilError(t, err)

	terms := make([]string, 0, len(termloc))
	for term := range termloc {
		terms = append(terms, term)
		delete(termloc, term)
	}
	sort.Strings(terms)
	return terms
}

// validate dti
func validateDocDti(t *testing.T, sdc client.StrongDocClient, key *sscrypto.StrongSaltKey, doc *TestDocumentIdxV1, terms []string) {
	// Open the DTI
	dtiVerison, err := doc.OpenDti(sdc, key)
	assert.NilError(t, err)

	defer dtiVerison.Close()

	dti, ok := dtiVerison.(*docidxv1.DocTermIdxV1)
	assert.Assert(t, ok)

	// Validate the DTI
	for err == nil {
		var blk *docidxv1.DocTermIdxBlkV1
		blk, err = dti.ReadNextBlock()
		if err != nil {
			assert.Equal(t, err, io.EOF)
		}
		if blk != nil {
			for _, term := range blk.Terms {
				assert.Assert(t, len(terms) > 0)
				assert.Equal(t, term, terms[0])
				terms = terms[1:]
			}
		}
	}

	assert.Equal(t, len(terms), 0)
	//doc.Close()
}

func TestCreateModifiedDoc(t *testing.T) {
	documents := 30
	versions := 3

	// ================================ Prev Test ================================
	sdc := prevTest(t)
	defer CleanupTemporaryDocumentIndex()
	// ================================ Test Doc Modified Indexes ================================
	key, err := sscrypto.GenerateKey(sscrypto.Type_XChaCha20)
	assert.NilError(t, err)

	docs, err := InitTestDocuments(documents, false)
	assert.NilError(t, err)
	oldDoc := docs[0]
	defer oldDoc.RemoveAllVersionsIndexes(sdc)
	err = oldDoc.CreateDoiAndDti(sdc, key)
	assert.NilError(t, err)

	docVers := make([][]*TestDocumentIdxV1, len(docs))
	for i, doc := range docs {
		docVers[i] = make([]*TestDocumentIdxV1, versions+1)
		docVers[i][0] = doc
	}

	for v := 0; v < versions; v++ {
		for _, oldDocs := range docVers {
			oldDoc := oldDocs[v]
			newDoc := testCreateModifiedDoc(t, sdc, oldDoc, key, v == 0)
			oldDocs[v+1] = newDoc
		}
	}
}

func testCreateModifiedDoc(t *testing.T, sdc client.StrongDocClient, oldDoc *TestDocumentIdxV1, key *sscrypto.StrongSaltKey, createOldIdx bool) *TestDocumentIdxV1 {
	addedTerms, deletedTerms := rand.Intn(99)+1, rand.Intn(99)+1

	// create modified doc
	newDoc, err := oldDoc.CreateModifiedDoc(addedTerms, deletedTerms)
	assert.NilError(t, err)
	assert.Equal(t, len(newDoc.AddedTerms), addedTerms)
	assert.Equal(t, len(newDoc.DeletedTerms), deletedTerms)

	if createOldIdx {
		err = oldDoc.CreateDoiAndDti(sdc, key)
		assert.NilError(t, err)
	}

	// create doi and dti for new doc
	err = newDoc.CreateDoiAndDti(sdc, key)
	assert.NilError(t, err)

	diffDocs(t, sdc, key, oldDoc, newDoc)

	return newDoc
}

func diffDocs(t *testing.T, sdc client.StrongDocClient, key *sscrypto.StrongSaltKey, oldDoc *TestDocumentIdxV1, newDoc *TestDocumentIdxV1) {
	// open dti of old doc
	oldTermList, oldTermMap, err := readAllTermsFromDti(oldDoc, sdc, key)
	assert.NilError(t, err)

	// open dti of new doc
	newTermList, newTermMap, err := readAllTermsFromDti(newDoc, sdc, key)
	assert.NilError(t, err)

	var dtiAdded []string
	for _, term := range newTermList {
		if !oldTermMap[term] {
			dtiAdded = append(dtiAdded, term)
		}
	}

	var docAdded []string
	for term := range newDoc.AddedTerms {
		docAdded = append(docAdded, term)
	}
	sort.Strings(dtiAdded)
	sort.Strings(docAdded)
	assert.DeepEqual(t, dtiAdded, docAdded)

	var dtiDeleted []string
	for _, term := range oldTermList {
		if !newTermMap[term] {
			dtiDeleted = append(dtiDeleted, term)
		}
	}

	var docDeleted []string
	for term := range newDoc.DeletedTerms {
		docDeleted = append(docDeleted, term)
	}
	sort.Strings(dtiDeleted)
	sort.Strings(docDeleted)
	assert.DeepEqual(t, dtiDeleted, docDeleted)
}

func readAllTermsFromDti(doc *TestDocumentIdxV1, sdc client.StrongDocClient, key *sscrypto.StrongSaltKey) ([]string, map[string]bool, error) {
	dtiVersion, err := doc.OpenDti(sdc, key)
	if err != nil {
		return nil, nil, err
	}
	defer dtiVersion.Close()

	dti := dtiVersion.(*docidxv1.DocTermIdxV1)
	return dti.ReadAllTerms()
}
