package mocks

import (
	"io"
	"slices"
)

type FileWriterMock struct {
	mkaAllDirErrorMap    map[string]error
	createErrorMap       map[string]error
	writeErrorMap        map[string]error
	saveJsonFileErrorMap map[string]error
	filesWritten         map[string]*WriteCloserMock
	directoriesCreated   map[string]bool
	jsonFilesSaved       map[string]interface{}
}

type WriteCloserMock struct {
	FileName       string
	FileContent    []byte
	Closed         bool
	WriteError     error
	FileWriterMock *FileWriterMock
}

func (writerCloserMock *WriteCloserMock) Close() error {
	writerCloserMock.Closed = true
	return nil
}

func (writerCloserMock *WriteCloserMock) Write(p []byte) (n int, err error) {
	if writerCloserMock.WriteError != nil {
		return 0, writerCloserMock.WriteError
	}
	writerCloserMock.FileContent = slices.Concat(writerCloserMock.FileContent, p)
	return len(p), nil
}

func NewFilewriterMock() *FileWriterMock {
	return &FileWriterMock{
		mkaAllDirErrorMap:    make(map[string]error),
		createErrorMap:       make(map[string]error),
		writeErrorMap:        make(map[string]error),
		saveJsonFileErrorMap: make(map[string]error),
		directoriesCreated:   map[string]bool{},
		filesWritten:         make(map[string]*WriteCloserMock),
		jsonFilesSaved:       make(map[string]interface{}),
	}
}

func (fileWriterMock *FileWriterMock) Create(name string) (io.WriteCloser, error) {
	if err := fileWriterMock.createErrorMap[name]; err != nil {
		return nil, err
	}
	fWritten := &WriteCloserMock{
		FileName:       name,
		FileWriterMock: fileWriterMock,
		WriteError:     fileWriterMock.writeErrorMap[name],
	}
	fileWriterMock.filesWritten[fWritten.FileName] = fWritten
	return fWritten, nil
}

func (fileWriterMock *FileWriterMock) MkdirAll(name string) error {
	if err := fileWriterMock.mkaAllDirErrorMap[name]; err != nil {
		return err
	}
	fileWriterMock.directoriesCreated[name] = true
	return nil
}

func (fileWriterMock *FileWriterMock) SaveJsonFile(name string, value interface{}) error {
	if err := fileWriterMock.saveJsonFileErrorMap[name]; err != nil {
		return err
	}
	fileWriterMock.jsonFilesSaved[name] = value
	return nil
}

func (fileWriterMock *FileWriterMock) AllFilesClosed() (allFilesClosed bool) {
	allFilesClosed = true
	for _, file := range fileWriterMock.filesWritten {
		allFilesClosed = allFilesClosed && file.Closed
	}
	return
}

func (fileWriterMock *FileWriterMock) FileNamesWritten() (fNames []string) {
	for fName := range fileWriterMock.filesWritten {
		fNames = append(fNames, fName)
	}
	for fname := range fileWriterMock.jsonFilesSaved {
		fNames = append(fNames, fname)
	}
	slices.Sort(fNames)
	return
}

func (fileWriterMock *FileWriterMock) JsonFileSaved(fname string) interface{} {
	return fileWriterMock.jsonFilesSaved[fname]
}

func (fileWriterMock *FileWriterMock) FileWritten(fname string) *WriteCloserMock {
	return fileWriterMock.filesWritten[fname]
}

func (fileWriterMock *FileWriterMock) DirectoriesCreated() (dirNames []string) {
	for fName := range fileWriterMock.directoriesCreated {
		dirNames = append(dirNames, fName)
	}
	slices.Sort(dirNames)
	return
}

func (fileWriterMock *FileWriterMock) ErrorOnCreate(name string, err error) {
	fileWriterMock.createErrorMap[name] = err
}

func (fileWriterMock *FileWriterMock) ErrorOnMkdirAll(name string, err error) {
	fileWriterMock.mkaAllDirErrorMap[name] = err
}

func (fileWriterMock *FileWriterMock) ErrorOnWrite(name string, err error) {
	fileWriterMock.writeErrorMap[name] = err
}

func (fileWriterMock *FileWriterMock) ErrorOnSaveJson(name string, err error) {
	fileWriterMock.saveJsonFileErrorMap[name] = err
}
