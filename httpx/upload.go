package httpx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
)

type FormType = int

const (
	ContentTypePattern        = "^multipart/form-data; boundary=(.*)$"
	ContentDispositionPattern = "^Content-Disposition: form-data; name=\"([^\"]*)\"(; filename=\"([^\"]*)\".*)?$"

	PLAIN FormType = 1
	FILE  FormType = 2
)

var (
	newLineMarker                  = []byte{13, 10}
	RegexContentTypePattern        = regexp.MustCompile(ContentTypePattern)
	RegexContentDispositionPattern = regexp.MustCompile(ContentDispositionPattern)
	formReaderType                 = reflect.TypeOf(&FileFormReader{})
)

type FileFormReader struct {
	request          *http.Request
	unReadableBuffer *bytes.Buffer
	atomByte         []byte
	newLineBytesPair []byte
	buffer           []byte
	newLineBuffer    *bytes.Buffer
}

type FileUploadHandler struct {
	Request              *http.Request
	paraBoundary         string
	endParaBoundary      string
	separator            []byte
	separatorTestBuffer  []byte
	separatorMergeBuffer []byte
	formReader           *FileFormReader
	// call when read a plain text field.
	OnFormField func(name string, value string, formType FormType)
	// call when about begin to read file body from form, need to provide an io.WriteCloser to write file bytes.
	OnFileField func(name string) io.WriteCloser
	// call when a file is end.
	OnEndFile func(name string, out io.WriteCloser)
}

// Unread return extra read bytes for next read.
func (reader *FileFormReader) Unread(read []byte) {
	reader.unReadableBuffer.Write(read)
}

// Read reads bytes from stream, if the buffer has bytes remain, read it first, then read from form body.
func (reader *FileFormReader) Read(buff []byte) (int, error) {
	// if buffer of FileFormReader has bytes cached before(from Unread()), then read it first.
	if reader.unReadableBuffer.Len() > 0 {
		if len(buff) <= reader.unReadableBuffer.Len() {
			return reader.unReadableBuffer.Read(buff)
		} else {
			offsetPos, err := reader.unReadableBuffer.Read(buff)
			if err != nil {
				return 0, err
			}
			// read directly from reader
			len, err := reader.request.Body.Read(buff[offsetPos:])
			if err != nil && err != io.EOF {
				return 0, err
			}
			return offsetPos + len, err
		}
	}
	// read directly from reader
	return reader.request.Body.Read(buff)
}

// beginUpload begin to read request entity and parse form field
func (handler *FileUploadHandler) Parse() error {
	handler.formReader = &FileFormReader{
		request:          handler.Request,
		unReadableBuffer: new(bytes.Buffer),
		atomByte:         make([]byte, 1),
		newLineBytesPair: make([]byte, 2),
		newLineBuffer:    new(bytes.Buffer),
		buffer:           make([]byte, 1024*30),
	}

	var fileIndex = 0

	headerContentType := handler.Request.Header["Content-Type"]
	contentType := ""
	if headerContentType != nil && len(headerContentType) > 0 {
		contentType = headerContentType[0]
	}
	if RegexContentTypePattern.Match([]byte(contentType)) {
		boundary := RegexContentTypePattern.ReplaceAllString(contentType, "${1}")
		handler.paraBoundary = "--" + boundary
		handler.endParaBoundary = "--" + boundary + "--"
		handler.separator = []byte("\r\n" + handler.paraBoundary)
		handler.separatorTestBuffer = make([]byte, len(handler.separator))
		handler.separatorMergeBuffer = make([]byte, len(handler.separator)*2)
		for {
			line, err := handler.formReader.readNextLine()
			if err != nil {
				return err
			}
			// if it is paraSeparator, then start read new form text field or file field
			if handler.paraBoundary == line {
				contentDisposition, err := handler.formReader.readNextLine()
				if err != nil {
					return err
				} else {
					mat1 := RegexContentDispositionPattern.Match([]byte(contentDisposition))
					paramName := ""
					paramValue := ""
					if mat1 {
						paramName = RegexContentDispositionPattern.ReplaceAllString(contentDisposition, "${1}")
					}

					paramContentType, err := handler.formReader.readNextLine()
					if err != nil {
						return err
					} else {
						if paramContentType == "" { // read text parameter field
							param, err := handler.formReader.readNextLine()
							if err != nil {
								return err
							} else {
								paramValue = param
								handler.OnFormField(paramName, paramValue, PLAIN)
							}
						} else { // parse content type
							mat2 := RegexContentDispositionPattern.Match([]byte(contentDisposition))
							fileName := ""
							if mat2 {
								fileName = RegexContentDispositionPattern.ReplaceAllString(contentDisposition, "${3}")
							}
							_, err = handler.formReader.readNextLine() // read blank line
							if err != nil {
								return err
							} else { // read file body
								handler.OnFormField(paramName, fileName, FILE)
								err := handler.readFileBody(fileName, handler.OnFileField(fileName))
								if err != nil {
									return err
								}
								fileIndex++
							}
						}
					}

				}
			} else if handler.endParaBoundary == line {
				// form stream hit end
				break
			} else {
				fmt.Println("unknown line")
			}
		}
	}

	return nil
}

// readNextLine reads next form field meta string.
func (reader *FileFormReader) readNextLine() (string, error) {
	reader.newLineBuffer.Reset()
	for {
		len, err := reader.Read(reader.atomByte)
		if err != nil && err != io.EOF {
			return "", err
		}
		if len != 1 {
			return "", errors.New("error read from stream[0]")
		}
		reader.newLineBytesPair[0] = reader.newLineBytesPair[1]
		reader.newLineBytesPair[1] = reader.atomByte[0]
		reader.newLineBuffer.Write(reader.atomByte)
		if bytes.Equal(newLineMarker, reader.newLineBytesPair) {
			return string(reader.newLineBuffer.Bytes()[0 : reader.newLineBuffer.Len()-2]), nil
		}
	}
}

// readFileBody reads a file body part.
func (handler *FileUploadHandler) readFileBody(fileName string, out io.WriteCloser) error {
	separatorLength := len(handler.separator)
	for {
		len1, err := handler.formReader.Read(handler.formReader.buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if len1 == 0 {
			return errors.New("read file body failed[0]")
		}
		// whether buff1 contains separator
		pos := bytes.Index(handler.formReader.buffer, handler.separator)
		if pos != -1 {
			out.Write(handler.formReader.buffer[0:pos])
			handler.formReader.Unread(handler.formReader.buffer[pos+2 : len1]) // skip "\r\n"
			handler.OnEndFile(fileName, out)
			break
		} else {
			len2, err := handler.formReader.Read(handler.separatorTestBuffer)
			if err != nil {
				if err != io.EOF {
					return err
				}
			}
			if len2 == 0 {
				return errors.New("read file body failed[1]")
			}
			// []byte tail is last bytes of buff1 and first bytes of buff2 in case broken separator.
			//
			if len1 >= separatorLength {
				ByteCopy(handler.separatorMergeBuffer, 0, separatorLength, handler.formReader.buffer[len1-separatorLength:len1])
			}
			if len2 >= separatorLength {
				ByteCopy(handler.separatorMergeBuffer, separatorLength, len(handler.separatorMergeBuffer), handler.separatorTestBuffer[0:separatorLength])
			}

			i2 := bytes.Index(handler.separatorMergeBuffer, handler.separator)
			if i2 != -1 {
				if i2 < separatorLength {
					_, err := out.Write(handler.formReader.buffer[0 : len1-i2])
					if err != nil {
						return err
					}
					handler.formReader.Unread(handler.formReader.buffer[len1-i2+2 : len1])
					handler.formReader.Unread(handler.separatorTestBuffer[0:len2])
				} else {
					_, err := out.Write(handler.formReader.buffer[0:len1])
					if err != nil {
						return err
					}
					handler.formReader.Unread(handler.separatorTestBuffer[i2-separatorLength+2 : len2])
				}
				handler.OnEndFile(fileName, out)
				break
			} else {
				_, err := out.Write(handler.formReader.buffer[0:len1])
				if err != nil {
					return err
				}
				handler.formReader.Unread(handler.separatorTestBuffer[0:len2])
			}
		}
	}
	return nil
}

// ByteCopy copies bytes
func ByteCopy(src []byte, start int, end int, cp []byte) {
	for i := range src {
		if i >= start && i < end {
			src[i] = cp[i]
		} else {
			break
		}
	}
}
