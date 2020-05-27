package utils

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type CommandUtilsTestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *CommandUtilsTestSuite) SetupTest() {
}

var mockedExitStatus int
var mockedStdout string


func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	es := strconv.Itoa(mockedExitStatus)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
		"STDOUT=" + mockedStdout,
		"EXIT_STATUS=" + es}
	return cmd
}

func TestExecCommandHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, os.Getenv("STDOUT"))
	i, _ := strconv.Atoi(os.Getenv("EXIT_STATUS"))
	os.Exit(i)
}


func(suite *CommandUtilsTestSuite) Test_FindCommand_WhenCommandDoesNotExist_ReturnsError() {

	command := "command"

	mockedExitStatus = 1
	//mockedStdout     = ""

	cmdService := CommandService{
		ExecCommand: fakeExecCommand,
	}

	location, err := cmdService.FindCommand(command)

	suite.NotNil(err)
	suite.Equal(true, len(location)==0)
}


func(suite *CommandUtilsTestSuite) Test_FindCommand_WhenCommandExists_ReturnCommandLocation() {
	command := "command"
	mockedStdout = "/usr/bin/command"
	mockedExitStatus = 0

	cmdService := CommandService{
		ExecCommand: fakeExecCommand,
	}

	location, err := cmdService.FindCommand(command)

	suite.Nil(err)
	suite.Equal(true, len(location)>0)
	suite.Equal(mockedStdout, location)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCommandUtilsSuite(t *testing.T) {
	suite.Run(t, new(CommandUtilsTestSuite))
}
