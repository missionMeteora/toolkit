package crypty

import (
	"io"
	"sync/atomic"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

const (
	// BlockSize is the standard block size
	BlockSize = aes.BlockSize

	// ErrIsClosed is returned when an action is attempted on a closed instance
	ErrIsClosed = Error("cannot perform an action on a closed instance")
)

// NewReaderPair will return a new Reader using a newly created Crypty
func NewReaderPair(in io.Reader, key, iv []byte) (nr *Reader, err error) {
	var (
		r   Reader
		blk cipher.Block
	)

	if blk, err = aes.NewCipher(key); err != nil {
		return
	}

	r.dec = cipher.NewCFBDecrypter(blk, iv)
	r.rdr = in
	nr = &r

	return
}

// NewReader will return a new Reader using a newly created Crypty
func NewReader(in io.Reader, key []byte) (nr *Reader, err error) {
	var (
		r   Reader
		blk cipher.Block

		iv = make([]byte, BlockSize)
	)

	if _, err = io.ReadFull(in, iv); err != nil {
		return
	}

	if blk, err = aes.NewCipher(key); err != nil {
		return
	}

	r.dec = cipher.NewCFBDecrypter(blk, iv)
	r.rdr = in
	nr = &r
	return
}

// Reader is an io.Reader interface for Crypty
type Reader struct {
	rdr io.Reader
	dec cipher.Stream

	closed uint32
}

func (r *Reader) isClosed() bool {
	return atomic.LoadUint32(&r.closed) == 1
}

func (r *Reader) Read(out []byte) (n int, err error) {
	if r.isClosed() {
		err = ErrIsClosed
		return
	}

	if n, err = r.rdr.Read(out); err != nil {
		return
	}

	r.dec.XORKeyStream(out, out)
	return
}

// Close will close a particular instance of Reader
func (r *Reader) Close() (err error) {
	if !atomic.CompareAndSwapUint32(&r.closed, 0, 1) {
		err = ErrIsClosed
	}

	return
}

// NewWriter will return a new Writer using a newly created Crypty
func NewWriter(out io.Writer, key []byte) (nw *Writer, err error) {
	var (
		w   Writer
		blk cipher.Block
		iv  []byte
	)

	if iv, err = GetIV(); err != nil {
		return
	}

	if _, err = out.Write(iv); err != nil {
		return
	}

	if blk, err = aes.NewCipher(key); err != nil {
		return
	}

	w.enc = cipher.NewCFBEncrypter(blk, iv)
	w.wtr = out
	nw = &w
	return
}

// NewWriterPair will return a new Writer using a newly created Crypty
func NewWriterPair(out io.Writer, key, iv []byte) (nw *Writer, err error) {
	var (
		w   Writer
		blk cipher.Block
	)

	if blk, err = aes.NewCipher(key); err != nil {
		return
	}

	w.enc = cipher.NewCFBEncrypter(blk, iv)
	w.wtr = out
	nw = &w
	return
}

// Writer is an io.Writer interface for Crypty
type Writer struct {
	wtr io.Writer
	enc cipher.Stream

	closed uint32
}

func (w *Writer) isClosed() bool {
	return atomic.LoadUint32(&w.closed) == 1
}

func (w *Writer) Write(in []byte) (n int, err error) {
	if w.isClosed() {
		err = ErrIsClosed
		return
	}

	out := make([]byte, len(in))
	w.enc.XORKeyStream(out, in)
	return w.wtr.Write(out)
}

// Close will close a particular instance of Writer
func (w *Writer) Close() (err error) {
	if !atomic.CompareAndSwapUint32(&w.closed, 0, 1) {
		err = ErrIsClosed
	}

	return
}

// Error is the constant error type
type Error string

func (e Error) Error() string {
	return string(e)
}

// GetIV will return an AES IV
func GetIV() (iv []byte, err error) {
	iv = make([]byte, BlockSize)
	_, err = rand.Read(iv)
	return
}

// GetIVFromFile will return an AES IV from an io.Reader
func GetIVFromFile(in io.Reader) (iv []byte, err error) {
	iv = make([]byte, BlockSize)
	_, err = io.ReadFull(in, iv)
	return
}
