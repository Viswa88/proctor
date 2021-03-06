package execution

import (
	"errors"
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/gojektech/proctor/daemon"
	"github.com/gojektech/proctor/io"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExecutionCmdTestSuite struct {
	suite.Suite
	mockPrinter             *io.MockPrinter
	mockProctorEngineClient *daemon.MockClient
	testExecutionCmd        *cobra.Command
}

func (s *ExecutionCmdTestSuite) SetupTest() {
	s.mockPrinter = &io.MockPrinter{}
	s.mockProctorEngineClient = &daemon.MockClient{}
	s.testExecutionCmd = NewCmd(s.mockPrinter, s.mockProctorEngineClient)
}

func (s *ExecutionCmdTestSuite) TestExecutionCmdUsage() {
	assert.Equal(s.T(), "execute", s.testExecutionCmd.Use)
}

func (s *ExecutionCmdTestSuite) TestExecutionCmdHelp() {
	assert.Equal(s.T(), "Execute a proc with arguments given", s.testExecutionCmd.Short)
	assert.Equal(s.T(), "To execute a proc, this command helps communicate with `proctord` and streams to logs of proc in execution", s.testExecutionCmd.Long)
	assert.Equal(s.T(), "proctor execute proc-one SOME_VAR=foo ANOTHER_VAR=bar\nproctor execute proc-two ANY_VAR=baz", s.testExecutionCmd.Example)
}

func (s *ExecutionCmdTestSuite) TestExecutionCmd() {
	args := []string{"say-hello-world", "SAMPLE_ARG_ONE=any", "SAMPLE_ARG_TWO=variable"}
	procArgs := make(map[string]string)
	procArgs["SAMPLE_ARG_ONE"] = "any"
	procArgs["SAMPLE_ARG_TWO"] = "variable"

	s.mockPrinter.On("Println", fmt.Sprintf("%-40s %-100s", "Executing Proc", "say-hello-world"), color.Reset).Once()
	s.mockPrinter.On("Println", "With Variables", color.FgMagenta).Once()
	s.mockPrinter.On("Println", fmt.Sprintf("%-40s %-100s", "SAMPLE_ARG_ONE", "any"), color.Reset).Once()
	s.mockPrinter.On("Println", fmt.Sprintf("%-40s %-100s", "SAMPLE_ARG_TWO", "variable"), color.Reset).Once()

	s.mockProctorEngineClient.On("ExecuteProc", "say-hello-world", procArgs).Return("executed-proc-name", nil).Once()

	s.mockPrinter.On("Println", "Proc execution successful. \nStreaming logs:", color.FgGreen).Once()

	s.mockProctorEngineClient.On("StreamProcLogs", "executed-proc-name").Return(nil).Once()
	s.mockPrinter.On("Println", "Log stream of proc completed.", color.FgGreen).Once()

	s.testExecutionCmd.Run(&cobra.Command{}, args)

	s.mockProctorEngineClient.AssertExpectations(s.T())
	s.mockPrinter.AssertExpectations(s.T())
}

func (s *ExecutionCmdTestSuite) TestExecutionCmdForIncorrectUsage() {
	s.mockPrinter.On("Println", "Incorrect command. See `proctor execute --help` for usage", color.FgRed).Once()

	s.testExecutionCmd.Run(&cobra.Command{}, []string{})

	s.mockPrinter.AssertExpectations(s.T())
}

func (s *ExecutionCmdTestSuite) TestExecutionCmdForNoProcVariables() {
	args := []string{"say-hello-world"}

	s.mockPrinter.On("Println", fmt.Sprintf("%-40s %-100s", "Executing Proc", "say-hello-world"), color.Reset).Once()
	s.mockPrinter.On("Println", "With No Variables", color.FgRed).Once()

	procArgs := make(map[string]string)
	s.mockProctorEngineClient.On("ExecuteProc", "say-hello-world", procArgs).Return("executed-proc-name", nil).Once()

	s.mockPrinter.On("Println", "Proc execution successful. \nStreaming logs:", color.FgGreen).Once()

	s.mockProctorEngineClient.On("StreamProcLogs", "executed-proc-name").Return(nil).Once()
	s.mockPrinter.On("Println", "Log stream of proc completed.", color.FgGreen).Once()

	s.testExecutionCmd.Run(&cobra.Command{}, args)

	s.mockProctorEngineClient.AssertExpectations(s.T())
	s.mockPrinter.AssertExpectations(s.T())
}

func (s *ExecutionCmdTestSuite) TestExecutionCmdForIncorrectVariableFormat() {
	args := []string{"say-hello-world", "incorrect-format"}

	s.mockPrinter.On("Println", fmt.Sprintf("%-40s %-100s", "Executing Proc", "say-hello-world"), color.Reset).Once()
	s.mockPrinter.On("Println", "With Variables", color.FgMagenta).Once()
	s.mockPrinter.On("Println", fmt.Sprintf("%-40s %-100s", "\nIncorrect variable format\n", "incorrect-format"), color.FgRed).Once()

	procArgs := make(map[string]string)
	s.mockProctorEngineClient.On("ExecuteProc", "say-hello-world", procArgs).Return("executed-proc-name", nil).Once()

	s.mockPrinter.On("Println", "Proc execution successful. \nStreaming logs:", color.FgGreen).Once()

	s.mockProctorEngineClient.On("StreamProcLogs", "executed-proc-name").Return(nil).Once()
	s.mockPrinter.On("Println", "Log stream of proc completed.", color.FgGreen).Once()

	s.testExecutionCmd.Run(&cobra.Command{}, args)

	s.mockProctorEngineClient.AssertExpectations(s.T())
	s.mockPrinter.AssertExpectations(s.T())
}

func (s *ExecutionCmdTestSuite) TestExecutionCmdForProctorEngineExecutionFailure() {
	args := []string{"say-hello-world"}

	s.mockPrinter.On("Println", fmt.Sprintf("%-40s %-100s", "Executing Proc", "say-hello-world"), color.Reset).Once()
	s.mockPrinter.On("Println", "With No Variables", color.FgRed).Once()

	procArgs := make(map[string]string)
	s.mockProctorEngineClient.On("ExecuteProc", "say-hello-world", procArgs).Return("", errors.New("test error")).Once()

	s.mockPrinter.On("Println", "test error", color.FgRed).Once()

	s.testExecutionCmd.Run(&cobra.Command{}, args)

	s.mockProctorEngineClient.AssertExpectations(s.T())
	s.mockPrinter.AssertExpectations(s.T())
}

func (s *ExecutionCmdTestSuite) TestExecutionCmdForProctorEngineLogStreamingFailure() {
	args := []string{"say-hello-world"}

	s.mockPrinter.On("Println", fmt.Sprintf("%-40s %-100s", "Executing Proc", "say-hello-world"), color.Reset).Once()
	s.mockPrinter.On("Println", "With No Variables", color.FgRed).Once()

	procArgs := make(map[string]string)
	s.mockProctorEngineClient.On("ExecuteProc", "say-hello-world", procArgs).Return("executed-proc-name", nil).Once()

	s.mockPrinter.On("Println", "Proc execution successful. \nStreaming logs:", color.FgGreen).Once()

	s.mockProctorEngineClient.On("StreamProcLogs", "executed-proc-name").Return(errors.New("error")).Once()
	s.mockPrinter.On("Println", "Error Streaming Logs", color.FgRed).Once()

	s.testExecutionCmd.Run(&cobra.Command{}, args)

	s.mockProctorEngineClient.AssertExpectations(s.T())
	s.mockPrinter.AssertExpectations(s.T())
}

func TestExecutionCmdTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutionCmdTestSuite))
}
