package internal

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testValidProcsFileContent = `Filename				Type		Size	Used	Priority
/home/swapfile			file		8388604	0	-2`

var _ logger = &testLogger{}

type testLogger struct {
	entries [][]any
}

func (l *testLogger) Println(v ...any) {
	l.entries = append(l.entries, v)
}

func Test_NewSwap(t *testing.T) {
	validFS := afero.NewMemMapFs()
	err := validFS.MkdirAll("/proc/", 0755)
	require.NoError(t, err, "creating mem fs dir /proc/")
	err = afero.WriteFile(validFS, "/proc/swaps", []byte(testValidProcsFileContent), 0444)
	require.NoError(t, err, "writing mem fs file /proc/swaps")

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
				fs := afero.NewMemMapFs()
				err := fs.MkdirAll("/proc/", 0755)
				require.NoError(t, err, "creating mem fs dir custom /proc/")
				err = afero.WriteFile(fs, "/proc/swaps", []byte("wrong file"), 0444)
				require.NoError(t, err, "writing mem fs file custom /proc/swaps")
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
		t.Run(tc.name, func(t *testing.T) {
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

func Test_getSwapFileLocation(t *testing.T) {
	tt := []struct {
		name                    string
		fileContent             []byte
		defaultSwapFileLocation string
		expectedResult          string
		expectedError           error
	}{
		{
			name:                    "happy path",
			fileContent:             []byte(testValidProcsFileContent),
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
				data := strings.Replace(testValidProcsFileContent, "/home/swapfile", "/dev/1", -1)
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
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			err := fs.MkdirAll("/proc/", 0755)
			require.NoError(t, err, "creating mem fs dir /proc/")
			if tc.fileContent != nil {
				err = afero.WriteFile(fs, "/proc/swaps", tc.fileContent, 0444)
				require.NoError(t, err, "writing mem fs file /proc/swaps")
			}

			res, err := getSwapFileLocation(fs, tc.defaultSwapFileLocation)

			require.Equal(t, fmt.Sprintf("%v", tc.expectedError), fmt.Sprintf("%v", err))
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}
