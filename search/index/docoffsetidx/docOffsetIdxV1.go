package docoffsetidx

import (
	"io"

	"github.com/go-errors/errors"

	"github.com/overnest/strongdoc-go-sdk/search/index/crypto"
	ssblocks "github.com/overnest/strongsalt-common-go/blocks"
	ssheaders "github.com/overnest/strongsalt-common-go/headers"
	sscrypto "github.com/overnest/strongsalt-crypto-go"
	sscryptointf "github.com/overnest/strongsalt-crypto-go/interfaces"
)

// The format off Document Offset Index
//
// --------------------------------------------------------------------------
// |   Unencrypted    |                   Encrypted                         |
// --------------------------------------------------------------------------
// | Plaintext Header | Ciphertext Header | Block Header | .... Blocks .... |
// --------------------------------------------------------------------------

//////////////////////////////////////////////////////////////////
//
//                   Document Offset Index
//
//////////////////////////////////////////////////////////////////

// DocOffsetIdxV1 is the Document Offset Index V1
type DocOffsetIdxV1 struct {
	DoiVersionS
	DocID         string
	DocVer        uint64
	Key           *sscrypto.StrongSaltKey
	Nonce         []byte
	InitOffset    uint64
	PlainHdrBody  *DoiPlainHdrBodyV1
	CipherHdrBody *DoiCipherHdrBodyV1
	Writer        ssblocks.BlockListWriterV1
	Reader        ssblocks.BlockListReaderV1
	Block         *DocOffsetIdxBlkV1
}

// CreateDocOffsetIdxV1 creates a document offset index writer V1
func CreateDocOffsetIdxV1(docID string, docVer uint64, key *sscrypto.StrongSaltKey,
	store interface{}, initOffset int64) (*DocOffsetIdxV1, error) {

	var err error
	writer, ok := store.(io.Writer)
	if !ok {
		return nil, errors.Errorf("The passed in storage does not implement io.Writer")
	}

	if key.Type != sscrypto.Type_XChaCha20 {
		return nil, errors.Errorf("Key type %v is not supported. The only supported key type is %v",
			key.Type.Name, sscrypto.Type_XChaCha20.Name)
	}

	// Create plaintext and ciphertext headers
	plainHdrBody := &DoiPlainHdrBodyV1{
		DoiVersionS: DoiVersionS{DoiVer: DOI_V1},
		KeyType:     key.Type.Name,
		DocID:       docID,
		DocVer:      docVer,
	}

	cipherHdrBody := &DoiCipherHdrBodyV1{
		BlockVersion: BlockVersion{BlockVer: DOI_BLOCK_V1},
	}

	if midStreamKey, ok := key.Key.(sscryptointf.KeyMidstream); ok {
		plainHdrBody.Nonce, err = midStreamKey.GenerateNonce()
		if err != nil {
			return nil, errors.New(err)
		}
	} else {
		return nil, errors.Errorf("The key type %v is not a midstream key", key.Type.Name)
	}

	plainHdrBodySerial, err := plainHdrBody.serialize()
	if err != nil {
		return nil, errors.New(err)
	}

	plainHdr := ssheaders.CreatePlainHdr(ssheaders.HeaderTypeJSONGzip, plainHdrBodySerial)
	plainHdrSerial, err := plainHdr.Serialize()
	if err != nil {
		return nil, errors.New(err)
	}

	cipherHdrBodySerial, err := cipherHdrBody.serialize()
	if err != nil {
		return nil, errors.New(err)
	}

	cipherHdr := ssheaders.CreateCipherHdr(ssheaders.HeaderTypeJSONGzip, cipherHdrBodySerial)
	cipherHdrSerial, err := cipherHdr.Serialize()
	if err != nil {
		return nil, errors.New(err)
	}

	// Write the plaintext header to storage
	n, err := writer.Write(plainHdrSerial)
	if err != nil {
		return nil, errors.New(err)
	}
	if n != len(plainHdrSerial) {
		return nil, errors.Errorf("Failed to write the entire plaintext header")
	}

	// Initialize the streaming crypto to encrypt ciphertext header and the
	// blocks after that
	streamCrypto, err := crypto.CreateStreamCrypto(key, plainHdrBody.Nonce, store,
		initOffset+int64(n))
	if err != nil {
		return nil, errors.New(err)
	}

	// Write the ciphertext header to storage
	n, err = streamCrypto.Write(cipherHdrSerial)
	if err != nil {
		return nil, errors.New(err)
	}
	if n != len(cipherHdrSerial) {
		return nil, errors.Errorf("Failed to write the entire ciphertext header")
	}

	// Create a block list writer using the streaming crypto so the blocks will be
	// encrypted.
	blockWriter, err := ssblocks.NewBlockListWriterV1(streamCrypto, 0,
		uint64(initOffset+int64(len(plainHdrSerial)+len(cipherHdrSerial))))
	if err != nil {
		return nil, errors.New(err)
	}

	index := &DocOffsetIdxV1{DoiVersionS{DoiVer: DOI_V1},
		docID, docVer, key, plainHdrBody.Nonce, uint64(initOffset),
		plainHdrBody, cipherHdrBody, blockWriter, nil, nil}
	return index, nil
}

// OpenDocOffsetIdxV1 opens a document offset index reader V1
func OpenDocOffsetIdxV1(key *sscrypto.StrongSaltKey, plainHdrBody *DoiPlainHdrBodyV1,
	store interface{}, initOffset int64) (*DocOffsetIdxV1, error) {

	if key.Type != sscrypto.Type_XChaCha20 {
		return nil, errors.Errorf("Key type %v is not supported. The only supported key type is %v",
			key.Type.Name, sscrypto.Type_XChaCha20.Name)
	}

	_, ok := store.(io.Reader)
	if !ok {
		return nil, errors.Errorf("The passed in storage does not implement io.Reader")
	}

	// Initialize the streaming crypto to decrypt ciphertext header and the blocks after that
	streamCrypto, err := crypto.CreateStreamCrypto(key, plainHdrBody.Nonce, store, initOffset)
	if err != nil {
		return nil, errors.New(err)
	}

	// Read the ciphertext header from storage
	cipherHdr, parsed, err := ssheaders.DeserializeCipherHdrStream(streamCrypto)
	if err != nil {
		return nil, errors.New(err)
	}

	cipherHdrBodyData, err := cipherHdr.GetBody()
	if err != nil {
		return nil, errors.New(err)
	}

	cipherHdrBody := &DoiCipherHdrBodyV1{}
	cipherHdrBody, err = cipherHdrBody.deserialize(cipherHdrBodyData)
	if err != nil {
		return nil, errors.New(err)
	}

	// Create a block list reader using the streaming crypto so the blocks will be
	// decrypted.
	reader, err := ssblocks.NewBlockListReader(streamCrypto,
		uint64(initOffset+int64(parsed)), 0)
	if err != nil {
		return nil, errors.New(err)
	}
	blockReader, ok := reader.(ssblocks.BlockListReaderV1)
	if !ok {
		return nil, errors.Errorf("Block list reader is not BlockListReaderV1")
	}

	index := &DocOffsetIdxV1{DoiVersionS{DoiVer: plainHdrBody.GetDoiVersion()},
		plainHdrBody.DocID, plainHdrBody.DocVer, key, plainHdrBody.Nonce,
		uint64(initOffset), plainHdrBody, cipherHdrBody, nil, blockReader, nil}
	return index, nil
}

func (idx *DocOffsetIdxV1) AddTermOffset(term string, offset uint64) error {
	if idx.Writer == nil {
		return errors.Errorf("The document offset index is not open for writing")
	}

	if idx.Block == nil {
		idx.Block = &DocOffsetIdxBlkV1{
			TermLoc: make(map[string][]uint64),
		}
	}

	idx.Block.AddTermOffset(term, offset)
	serial, err := idx.Block.Serialize()
	if err != nil {
		return errors.New(err)
	}

	if len(serial) > DOI_BLOCK_SIZE_MAX {
		return idx.flush(serial)
	}

	return nil
}

func (idx *DocOffsetIdxV1) ReadNextBlock() (*DocOffsetIdxBlkV1, error) {
	if idx.Reader == nil {
		return nil, errors.Errorf("The document offset index is not open for reading")
	}

	b, err := idx.Reader.ReadNextBlock()
	if err != nil {
		return nil, errors.New(err)
	}

	block := &DocOffsetIdxBlkV1{}
	return block.Deserialize(b.GetData())
}

func (idx *DocOffsetIdxV1) Close() error {
	if idx.Block != nil {
		serial, err := idx.Block.Serialize()
		if err != nil {
			return errors.New(err)
		}
		return idx.flush(serial)
	}
	return nil
}

func (idx *DocOffsetIdxV1) flush(data []byte) error {
	_, err := idx.Writer.WriteBlockData(data)
	if err != nil {
		return errors.New(err)
	}
	idx.Block = nil
	return nil
}
