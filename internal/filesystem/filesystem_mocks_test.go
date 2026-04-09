package filesystem

import (
	"os"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockDir struct {
	mock.Mock
}

func (m *MockDir) Readdir(count int) ([]os.FileInfo, error) {
	args := m.Called(count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

type MockOSOperations struct {
	mock.Mock
}

func (m *MockOSOperations) Mkdir(name string, perm os.FileMode) error {
	args := m.Called(name, perm)
	return args.Error(0)
}

func (m *MockOSOperations) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockOSOperations) Open(name string) (File, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(File), args.Error(1)
}

func (m *MockOSOperations) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockOSOperations) WriteFile(name string, data []byte, perm os.FileMode) error {
	args := m.Called(name, data, perm)
	return args.Error(0)
}

type MockFile struct {
	mock.Mock
}

func (m *MockFile) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockFile) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockFile) Readdir(count int) ([]os.FileInfo, error) {
	args := m.Called(count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

type FakeFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (f FakeFileInfo) Name() string       { return f.name }
func (f FakeFileInfo) Size() int64        { return f.size }
func (f FakeFileInfo) Mode() os.FileMode  { return 0 }
func (f FakeFileInfo) ModTime() time.Time { return time.Now() }
func (f FakeFileInfo) IsDir() bool        { return f.isDir }
func (f FakeFileInfo) Sys() any           { return nil }

type FileSystemTestSuite struct {
	suite.Suite
	fs       *FileSystemWrapper
	mockOS   *MockOSOperations
	mockFile *MockFile
	path     string
}
