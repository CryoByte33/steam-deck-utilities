package internal

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var _ logger = &testLogger{}

type testLogger struct {
	entries [][]any
}

func (l *testLogger) Println(v ...any) {
	l.entries = append(l.entries, v)
}

var _ execCommander = &testCommander{}

type testCommander struct {
	// map key is the name + " " + arg[0]
	shouldReturn     map[string]testCommanderReturn
	commandsExecuted []struct {
		name string
		args []string
	}
}

type testCommanderReturn struct {
	value string
	err   error
}

func (c *testCommander) ExecAndOutput(name string, arg ...string) ([]byte, error) {
	c.commandsExecuted = append(c.commandsExecuted, struct {
		name string
		args []string
	}{name: name, args: arg})

	searchBy := name
	if len(arg) > 0 {
		searchBy += " " + strings.Join(arg, " ")
	}
	shouldReturn, ok := c.shouldReturn[searchBy]
	if !ok {
		panic(fmt.Errorf("failed to find a result to return for %q", searchBy))
	}

	return []byte(shouldReturn.value), shouldReturn.err
}

type SwapTestSuite struct {
	suite.Suite
	validProcsFileContent string
}

func Test_SwapTestSuite(t *testing.T) {
	s := new(SwapTestSuite)
	s.validProcsFileContent = `Filename				Type		Size	Used	Priority
/home/swapfile			file		8388604	0	-2`

	suite.Run(t, s)
}

func (s *SwapTestSuite) getValidFS() afero.Fs {
	validFS := afero.NewMemMapFs()
	err := validFS.MkdirAll("/proc/", 0755)
	require.NoError(s.T(), err, "creating mem fs dir /proc/")
	err = afero.WriteFile(validFS, "/proc/swaps", []byte(s.validProcsFileContent), 0444)
	require.NoError(s.T(), err, "writing mem fs file /proc/swaps")

	return validFS
}

func (s *SwapTestSuite) Test_NewSwap() {
	validFS := s.getValidFS()

	tt := []struct {
		name                    string
		defaultSwapSizeBytes    int64
		availableSwapSizes      []string
		oldSwappinessUnitFile   string
		defaultSwapFileLocation string
		fs                      afero.Fs
		loggerInfo              logger
		expectedResult          *Swap
		expectedError           error
	}{
		{
			name:                    "happy path",
			defaultSwapSizeBytes:    1,
			availableSwapSizes:      []string{"2", "4"},
			oldSwappinessUnitFile:   "/etc/sysctl.d/zzz-custom-swappiness.conf",
			defaultSwapFileLocation: "/home/swapfile",
			fs:                      validFS,
			loggerInfo:              &testLogger{},
			expectedResult: &Swap{
				defaultSwapSizeBytes:  1,
				availableSwapSizes:    []string{"2", "4"},
				oldSwappinessUnitFile: "/etc/sysctl.d/zzz-custom-swappiness.conf",
				swapFileLocation:      "/home/swapfile",
				fs:                    validFS,
				execCommander:         realExecCommander{},
				loggerInfo:            &testLogger{},
			},
			expectedError: nil,
		},
		{
			name:                    "zero default swap file size - error",
			defaultSwapSizeBytes:    0,
			availableSwapSizes:      []string{"2", "4"},
			oldSwappinessUnitFile:   "/etc/sysctl.d/zzz-custom-swappiness.conf",
			defaultSwapFileLocation: "/home/swapfile",
			fs:                      validFS,
			loggerInfo:              &testLogger{},
			expectedResult:          nil,
			expectedError:           errors.New("defaultSwapSizeBytes is required"),
		},
		{
			name:                    "no available swap sizes - error",
			defaultSwapSizeBytes:    1,
			availableSwapSizes:      nil,
			oldSwappinessUnitFile:   "/etc/sysctl.d/zzz-custom-swappiness.conf",
			defaultSwapFileLocation: "/home/swapfile",
			fs:                      validFS,
			loggerInfo:              &testLogger{},
			expectedResult:          nil,
			expectedError:           errors.New("availableSwapSizes is required"),
		},
		{
			name:                    "no oldSwappinessUnitFile - error",
			defaultSwapSizeBytes:    1,
			availableSwapSizes:      []string{"2", "4"},
			oldSwappinessUnitFile:   "",
			defaultSwapFileLocation: "/home/swapfile",
			fs:                      validFS,
			loggerInfo:              &testLogger{},
			expectedResult:          nil,
			expectedError:           errors.New("oldSwappinessUnitFile is required"),
		},
		{
			name:                    "no defaultSwapFileLocation - error",
			defaultSwapSizeBytes:    1,
			availableSwapSizes:      []string{"2", "4"},
			oldSwappinessUnitFile:   "/etc/sysctl.d/zzz-custom-swappiness.conf",
			defaultSwapFileLocation: "",
			fs:                      validFS,
			loggerInfo:              &testLogger{},
			expectedResult:          nil,
			expectedError:           errors.New("default swap location is required"),
		},
		{
			name:                    "no fs - error",
			defaultSwapSizeBytes:    1,
			availableSwapSizes:      []string{"2", "4"},
			oldSwappinessUnitFile:   "/etc/sysctl.d/zzz-custom-swappiness.conf",
			defaultSwapFileLocation: "/home/swapfile",
			fs:                      nil,
			loggerInfo:              &testLogger{},
			expectedResult:          nil,
			expectedError:           errors.New("fs is required"),
		},
		{
			name:                    "wrong /proc/swaps - error",
			defaultSwapSizeBytes:    1,
			availableSwapSizes:      []string{"2", "4"},
			oldSwappinessUnitFile:   "/etc/sysctl.d/zzz-custom-swappiness.conf",
			defaultSwapFileLocation: "/home/swapfile",
			fs: func() afero.Fs {
				fs := s.getValidFS()
				err := afero.WriteFile(fs, "/proc/swaps", []byte("wrong file"), 0444)
				require.NoError(s.T(), err, "writing mem fs file custom /proc/swaps")
				return fs
			}(),
			loggerInfo:     &testLogger{},
			expectedResult: nil,
			expectedError:  errors.New("getting swapfile location: no swapfile found in /proc/swaps"),
		},
		{
			name:                    "no info logger - error",
			defaultSwapSizeBytes:    1,
			availableSwapSizes:      []string{"2", "4"},
			oldSwappinessUnitFile:   "/etc/sysctl.d/zzz-custom-swappiness.conf",
			defaultSwapFileLocation: "/home/swapfile",
			fs:                      validFS,
			loggerInfo:              nil,
			expectedResult:          nil,
			expectedError:           errors.New("info logger is required"),
		},
	}

	for _, tc := range tt {
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := NewSwap(
				tc.defaultSwapSizeBytes,
				tc.availableSwapSizes,
				tc.oldSwappinessUnitFile,
				tc.defaultSwapFileLocation,
				tc.fs,
				tc.loggerInfo,
			)

			require.Equal(t, fmt.Sprintf("%v", tc.expectedError), fmt.Sprintf("%v", err))
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func (s *SwapTestSuite) Test_getSwapFileLocation() {
	tt := []struct {
		name                    string
		fileContent             []byte
		defaultSwapFileLocation string
		expectedResult          string
		expectedError           error
	}{
		{
			name:                    "happy path",
			fileContent:             []byte(s.validProcsFileContent),
			defaultSwapFileLocation: "/home/swapfile__default",
			expectedResult:          "/home/swapfile", // comes from procs file
			expectedError:           nil,
		},
		{
			name:                    "no procs file - error",
			fileContent:             nil,
			defaultSwapFileLocation: "/home/swapfile__default",
			expectedResult:          "",
			expectedError:           errors.New("open /proc/swaps: file does not exist"),
		},
		{
			name: "swapfile is a partition - error",
			fileContent: func() []byte {
				data := strings.Replace(s.validProcsFileContent, "/home/swapfile", "/dev/1", -1)
				return []byte(data)
			}(),
			defaultSwapFileLocation: "/home/swapfile__default",
			expectedResult:          "",
			expectedError:           errors.New("no swapfile found in /proc/swaps"),
		},
		{
			name:                    "no swapfile mentioned in /proc/swaps and default not exists - error",
			fileContent:             []byte("some_data"),
			defaultSwapFileLocation: "/home/swapfile__default",
			expectedResult:          "",
			expectedError:           errors.New("no swapfile found in /proc/swaps"),
		},
	}

	for _, tc := range tt {
		s.T().Run(tc.name, func(t *testing.T) {
			fs := s.getValidFS()
			if tc.fileContent == nil {
				err := fs.Remove("/proc/swaps")
				require.NoError(t, err, "removing from mem fs file /proc/swaps")
			} else {
				err := afero.WriteFile(fs, "/proc/swaps", tc.fileContent, 0444)
				require.NoError(t, err, "writing mem fs file /proc/swaps")
			}

			res, err := getSwapFileLocation(fs, tc.defaultSwapFileLocation)

			require.Equal(t, fmt.Sprintf("%v", tc.expectedError), fmt.Sprintf("%v", err))
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func (s *SwapTestSuite) Test_Swap_getSwappinessValue() {
	tt := []struct {
		name                      string
		testCommanderShouldReturn map[string]testCommanderReturn
		expectedResult            int
		expectedError             error
		expectedLogs              [][]any
	}{
		{
			name: "happy path",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sysctl vm.swappiness": {
					value: "vm.swappiness = 10",
					err:   nil,
				},
			},
			expectedResult: 10,
			expectedError:  nil,
			expectedLogs:   [][]any{{"Found a swappiness of", "10"}},
		},
		{
			name: "commander error - error",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sysctl vm.swappiness": {
					value: "",
					err:   errors.New("commander error"),
				},
			},
			expectedResult: 100, // default
			expectedError:  errors.New("error getting current swappiness: commander error"),
			expectedLogs:   nil,
		},
		{
			name: "unexpected value from commander - error",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sysctl vm.swappiness": {
					value: "",
					err:   nil,
				},
			},
			expectedResult: 100, // default
			expectedError:  errors.New("unexpected swappiness returned: \"\""),
			expectedLogs:   nil,
		},
		{
			name: "non-int swappiness value from commander - error",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sysctl vm.swappiness": {
					value: "vm.swappiness = ten",
					err:   nil,
				},
			},
			expectedResult: 100, // default
			expectedError:  errors.New("converting swappiness of \"vm.swappiness = ten\": strconv.Atoi: parsing \"ten\": invalid syntax"),
			expectedLogs:   nil,
		},
	}

	for _, tc := range tt {
		s.T().Run(tc.name, func(t *testing.T) {
			fs := s.getValidFS()
			lgr := &testLogger{}
			swap, err := NewSwap(DefaultSwapSizeBytes, AvailableSwapSizes, OldSwappinessUnitFile, DefaultSwapFileLocation, fs, lgr)
			require.NoError(t, err, "instantiating a new swap")
			commander := &testCommander{
				shouldReturn:     tc.testCommanderShouldReturn,
				commandsExecuted: nil,
			}
			swap.execCommander = commander

			res, err := swap.getSwappinessValue()

			require.Equal(t, fmt.Sprintf("%v", tc.expectedError), fmt.Sprintf("%v", err))
			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedLogs, lgr.entries)
		})
	}
}

func (s *SwapTestSuite) Test_Swap_getSwapFileSize() {
	validSwapFileContent := "swap file content to have non-zero size"
	tt := []struct {
		name            string
		swapfileContent []byte
		expectedResult  int64
		expectedError   error
		expectedLogs    [][]any
	}{
		{
			name:            "happy path",
			swapfileContent: []byte(validSwapFileContent),
			expectedResult:  int64(len(validSwapFileContent)),
			expectedError:   nil,
			expectedLogs:    [][]any{{"Found a swap file with a size of", int64(len(validSwapFileContent))}},
		},
		{
			name:            "no swap file - error",
			swapfileContent: nil,
			expectedResult:  0,
			expectedError:   errors.New("error getting current swap file size: open /home/swapfile: file does not exist"),
			expectedLogs:    nil,
		},
	}

	for _, tc := range tt {
		s.T().Run(tc.name, func(t *testing.T) {
			fs := s.getValidFS()
			if tc.swapfileContent != nil {
				err := afero.WriteFile(fs, "/home/swapfile", tc.swapfileContent, 0444)
				require.NoError(t, err, "writing mem fs file /home/swapfile")
			}
			lgr := &testLogger{}
			swap, err := NewSwap(DefaultSwapSizeBytes, AvailableSwapSizes, OldSwappinessUnitFile, DefaultSwapFileLocation, fs, lgr)
			require.NoError(t, err, "instantiating a new swap")

			res, err := swap.getSwapFileSize()

			require.Equal(t, fmt.Sprintf("%v", tc.expectedError), fmt.Sprintf("%v", err))
			if err == nil {
				assert.Equal(t, tc.expectedResult, res)
			}
			assert.Equal(t, tc.expectedLogs, lgr.entries)
		})
	}
}

func (s *SwapTestSuite) Test_Swap_getAvailableSwapSizes() {
	// TODO - now it depends on utilities and thus depends on real os; need to refactor before testing
}

func (s *SwapTestSuite) Test_Swap_disableSwap() {
	tt := []struct {
		name                      string
		testCommanderShouldReturn map[string]testCommanderReturn
		expectedError             error
		expectedLogs              [][]any
	}{
		{
			name: "happy path",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sudo swapoff -a": {
					value: "",
					err:   nil,
				},
			},
			expectedError: nil,
			expectedLogs:  [][]any{{"Disabling swap temporarily..."}},
		},
		{
			name: "error from commander - return error",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sudo swapoff -a": {
					value: "",
					err:   errors.New("commander error"),
				},
			},
			expectedError: errors.New("error disabling swap: commander error"),
			expectedLogs:  [][]any{{"Disabling swap temporarily..."}},
		},
	}

	for _, tc := range tt {
		s.T().Run(tc.name, func(t *testing.T) {
			fs := s.getValidFS()
			lgr := &testLogger{}
			swap, err := NewSwap(DefaultSwapSizeBytes, AvailableSwapSizes, OldSwappinessUnitFile, DefaultSwapFileLocation, fs, lgr)
			require.NoError(t, err, "instantiating a new swap")
			commander := &testCommander{
				shouldReturn:     tc.testCommanderShouldReturn,
				commandsExecuted: nil,
			}
			swap.execCommander = commander

			err = swap.disableSwap()

			require.Equal(t, fmt.Sprintf("%v", tc.expectedError), fmt.Sprintf("%v", err))
			assert.Equal(t, tc.expectedLogs, lgr.entries)
		})
	}
}

func (s *SwapTestSuite) Test_Swap_resizeSwapFile() {
	tt := []struct {
		name                      string
		size                      int
		testCommanderShouldReturn map[string]testCommanderReturn
		expectedError             error
		expectedLogs              [][]any
	}{
		{
			name: "happy path",
			size: 16,
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sudo dd if=/dev/zero of=/home/swapfile bs=1G count=16 status=progress": {
					value: "",
					err:   nil,
				},
			},
			expectedError: nil,
			expectedLogs:  [][]any{{"Resizing swap to", 16, "GB..."}},
		},
		{
			name: "commander error - return error",
			size: 16,
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sudo dd if=/dev/zero of=/home/swapfile bs=1G count=16 status=progress": {
					value: "",
					err:   errors.New("commander error"),
				},
			},
			expectedError: errors.New("error resizing /home/swapfile: commander error"),
			expectedLogs:  [][]any{{"Resizing swap to", 16, "GB..."}},
		},
	}

	for _, tc := range tt {
		s.T().Run(tc.name, func(t *testing.T) {
			fs := s.getValidFS()
			lgr := &testLogger{}
			swap, err := NewSwap(DefaultSwapSizeBytes, AvailableSwapSizes, OldSwappinessUnitFile, DefaultSwapFileLocation, fs, lgr)
			require.NoError(t, err, "instantiating a new swap")
			commander := &testCommander{
				shouldReturn:     tc.testCommanderShouldReturn,
				commandsExecuted: nil,
			}
			swap.execCommander = commander

			err = swap.resizeSwapFile(tc.size)

			require.Equal(t, fmt.Sprintf("%v", tc.expectedError), fmt.Sprintf("%v", err))
			assert.Equal(t, tc.expectedLogs, lgr.entries)
		})
	}
}

func (s *SwapTestSuite) Test_Swap_setSwapPermissions() {
	tt := []struct {
		name                      string
		testCommanderShouldReturn map[string]testCommanderReturn
		expectedError             error
		expectedLogs              [][]any
	}{
		{
			name: "happy path",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sudo chmod 600 /home/swapfile": {
					value: "",
					err:   nil,
				},
			},
			expectedError: nil,
			expectedLogs:  [][]any{{"Setting permissions on", "/home/swapfile", "to 0600..."}},
		},
		{
			name: "commander error - return error",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sudo chmod 600 /home/swapfile": {
					value: "",
					err:   errors.New("commander error"),
				},
			},
			expectedError: errors.New("error setting permissions on /home/swapfile: commander error"),
			expectedLogs:  [][]any{{"Setting permissions on", "/home/swapfile", "to 0600..."}},
		},
	}

	for _, tc := range tt {
		s.T().Run(tc.name, func(t *testing.T) {
			fs := s.getValidFS()
			lgr := &testLogger{}
			swap, err := NewSwap(DefaultSwapSizeBytes, AvailableSwapSizes, OldSwappinessUnitFile, DefaultSwapFileLocation, fs, lgr)
			require.NoError(t, err, "instantiating a new swap")
			commander := &testCommander{
				shouldReturn:     tc.testCommanderShouldReturn,
				commandsExecuted: nil,
			}
			swap.execCommander = commander

			err = swap.setSwapPermissions()

			require.Equal(t, fmt.Sprintf("%v", tc.expectedError), fmt.Sprintf("%v", err))
			assert.Equal(t, tc.expectedLogs, lgr.entries)
		})
	}
}

func (s *SwapTestSuite) Test_Swap_initNewSwapFile() {
	tt := []struct {
		name                      string
		testCommanderShouldReturn map[string]testCommanderReturn
		expectedError             error
		expectedLogs              [][]any
	}{
		{
			name: "happy path",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sudo mkswap /home/swapfile": {
					value: "",
					err:   nil,
				},
				"sudo swapon /home/swapfile": {
					value: "",
					err:   nil,
				},
			},
			expectedError: nil,
			expectedLogs:  [][]any{{"Enabling swap on", "/home/swapfile", "..."}},
		},
		{
			name: "mkswap error - return error",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sudo mkswap /home/swapfile": {
					value: "",
					err:   errors.New("mkswap error"),
				},
				"sudo swapon /home/swapfile": {
					value: "",
					err:   nil,
				},
			},
			expectedError: errors.New("error creating swap on /home/swapfile: mkswap error"),
			expectedLogs:  [][]any{{"Enabling swap on", "/home/swapfile", "..."}},
		},
		{
			name: "swapon error - return error",
			testCommanderShouldReturn: map[string]testCommanderReturn{
				"sudo mkswap /home/swapfile": {
					value: "",
					err:   nil,
				},
				"sudo swapon /home/swapfile": {
					value: "",
					err:   errors.New("swapon error"),
				},
			},
			expectedError: errors.New("error enabling swap on /home/swapfile: swapon error"),
			expectedLogs:  [][]any{{"Enabling swap on", "/home/swapfile", "..."}},
		},
	}

	for _, tc := range tt {
		s.T().Run(tc.name, func(t *testing.T) {
			fs := s.getValidFS()
			lgr := &testLogger{}
			swap, err := NewSwap(DefaultSwapSizeBytes, AvailableSwapSizes, OldSwappinessUnitFile, DefaultSwapFileLocation, fs, lgr)
			require.NoError(t, err, "instantiating a new swap")
			commander := &testCommander{
				shouldReturn:     tc.testCommanderShouldReturn,
				commandsExecuted: nil,
			}
			swap.execCommander = commander

			err = swap.initNewSwapFile()

			require.Equal(t, fmt.Sprintf("%v", tc.expectedError), fmt.Sprintf("%v", err))
			assert.Equal(t, tc.expectedLogs, lgr.entries)
		})
	}
}

func (s *SwapTestSuite) Test_Swap_ChangeSwappiness() {
	// TODO - now it depends on utilities and thus depends on real os; need to refactor before testing
}
