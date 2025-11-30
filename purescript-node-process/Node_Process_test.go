package purescript_node_process

import (
	"os"
	"runtime"
	"testing"

	. "github.com/purescript-native/go-runtime"
)

func TestArgv(t *testing.T) {
	exports := Foreign("Node.Process")
	argv := exports["argv"].(func() Any)().([]Any)
	
	if len(argv) == 0 {
		t.Error("argv should not be empty")
	}
	
	// First arg should be executable path
	if _, ok := argv[0].(string); !ok {
		t.Error("argv[0] should be a string")
	}
}

func TestExecPath(t *testing.T) {
	exports := Foreign("Node.Process")
	execPath := exports["execPath"].(func() Any)().(string)
	
	if execPath == "" {
		t.Error("execPath should not be empty")
	}
}

func TestCwd(t *testing.T) {
	exports := Foreign("Node.Process")
	cwd := exports["cwd"].(func() Any)().(string)
	
	if cwd == "" {
		t.Error("cwd should not be empty")
	}
	
	// Should match os.Getwd()
	expected, _ := os.Getwd()
	if cwd != expected {
		t.Errorf("cwd = %s, want %s", cwd, expected)
	}
}

func TestLookupEnv(t *testing.T) {
	exports := Foreign("Node.Process")
	
	// Set a test env var
	testKey := "TEST_GO_FFI_VAR"
	testValue := "test_value_123"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)
	
	lookupEnv := exports["lookupEnv"].(func(Any) Any)
	result := lookupEnv(testKey).(func() Any)().(Dict)
	
	// Should be Just value
	if val, ok := result["value0"].(string); ok {
		if val != testValue {
			t.Errorf("lookupEnv returned %s, want %s", val, testValue)
		}
	} else {
		t.Error("lookupEnv should return Just for existing variable")
	}
	
	// Test non-existent variable
	result2 := lookupEnv("NONEXISTENT_VAR_XYZ").(func() Any)().(Dict)
	if len(result2) != 0 {
		t.Error("lookupEnv should return Nothing for non-existent variable")
	}
}

func TestSetEnv(t *testing.T) {
	exports := Foreign("Node.Process")
	
	testKey := "TEST_SET_ENV_VAR"
	testValue := "set_value_456"
	
	setEnv := exports["setEnv"].(func(Any, Any) Any)
	setEnv(testKey, testValue).(func() Any)()
	
	// Verify it was set
	if val := os.Getenv(testKey); val != testValue {
		t.Errorf("setEnv failed: got %s, want %s", val, testValue)
	}
	
	os.Unsetenv(testKey)
}

func TestPlatform(t *testing.T) {
	exports := Foreign("Node.Process")
	platform := exports["platform"].(string)
	
	if platform != runtime.GOOS {
		t.Errorf("platform = %s, want %s", platform, runtime.GOOS)
	}
}

func TestArch(t *testing.T) {
	exports := Foreign("Node.Process")
	arch := exports["arch"].(string)
	
	if arch != runtime.GOARCH {
		t.Errorf("arch = %s, want %s", arch, runtime.GOARCH)
	}
}

func TestPid(t *testing.T) {
	exports := Foreign("Node.Process")
	pid := exports["pid"].(int)
	
	if pid != os.Getpid() {
		t.Errorf("pid = %d, want %d", pid, os.Getpid())
	}
}

func TestVersion(t *testing.T) {
	exports := Foreign("Node.Process")
	version := exports["version"].(string)
	
	if version != runtime.Version() {
		t.Errorf("version = %s, want %s", version, runtime.Version())
	}
}

func TestStdStreams(t *testing.T) {
	exports := Foreign("Node.Process")
	
	if exports["stdin"] != os.Stdin {
		t.Error("stdin should be os.Stdin")
	}
	
	if exports["stdout"] != os.Stdout {
		t.Error("stdout should be os.Stdout")
	}
	
	if exports["stderr"] != os.Stderr {
		t.Error("stderr should be os.Stderr")
	}
}

