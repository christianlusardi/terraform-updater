package utils

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type FileUtilsTestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *FileUtilsTestSuite) SetupTest() {
}


func (suite *FileUtilsTestSuite) Test_MoveFile_WhenFileDoesNotExists_ShouldReturnError() {

	fileService := FileService{
		Os: afero.NewMemMapFs(),
	}

	err := fileService.MoveFile("a", "b")

	suite.NotNil(err)
	suite.Contains(err.Error(), "file does not exist")
}


func (suite *FileUtilsTestSuite) Test_MoveFile_WhenAllIsOk_ShouldWork() {

	fileService := FileService{
		Os: afero.NewMemMapFs(),
	}

	fileService.Os.MkdirAll("/home/test/source/", 0777)
	fileService.Os.Create("/home/test/source/source.txt")
	fileService.Os.MkdirAll("/home/test/dest/", 0777)


	err := fileService.MoveFile("/home/test/source/source.txt", "/home/test/dest/source.txt")

	suite.Nil(err)


}




// ---------- FileExists ----------

func (suite *FileUtilsTestSuite) Test_FileExists_WhenFileExists_ReturnTrue() {

	fileService := FileService{
		Os: afero.NewMemMapFs(),
	}

	fileService.Os.MkdirAll("/home/test/source/", 0777)
	f, err := fileService.Os.Create("/home/test/source/source.txt")

	if err != nil {
		suite.Error(err)
	}

	f.WriteString("Hello World")

	if err != nil {
		f.Close()
		suite.Error(err)
	}

	suite.Equal(true, fileService.FileExists("/home/test/source/source.txt"))

}


func (suite *FileUtilsTestSuite) Test_FileExists_WhenFileDoesNotExists_ReturnFalse() {

	fileService := FileService{
		Os: afero.NewMemMapFs(),
	}

	suite.Equal(false, fileService.FileExists("/home/test/source/source.txt"))

}



// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFileServiceSuite(t *testing.T) {
	suite.Run(t, new(FileUtilsTestSuite))
}
