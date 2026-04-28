package dotenv

import (
	"os"
	"testing"
)

func TestLoadSource_File(t *testing.T) {
	f := writeTempEnv(t, "APP=hello\nDEBUG=true\n")
	src, err := LoadSource(SourceSpec{Type: SourceTypeFile, Ref: f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Name != f {
		t.Errorf("expected name %q, got %q", f, src.Name)
	}
	if src.Vars["APP"] != "hello" {
		t.Errorf("expected APP=hello, got %q", src.Vars["APP"])
	}
	if src.Vars["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", src.Vars["DEBUG"])
	}
}

func TestLoadSource_FileNotFound(t *testing.T) {
	_, err := LoadSource(SourceSpec{Type: SourceTypeFile, Ref: "/no/such/file.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSource_Env_NoPrefix(t *testing.T) {
	t.Setenv("_ENVOY_TEST_KEY", "world")
	src, err := LoadSource(SourceSpec{Type: SourceTypeEnv, Ref: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Name != "<env>" {
		t.Errorf("expected name <env>, got %q", src.Name)
	}
	if src.Vars["_ENVOY_TEST_KEY"] != "world" {
		t.Errorf("expected _ENVOY_TEST_KEY=world, got %q", src.Vars["_ENVOY_TEST_KEY"])
	}
}

func TestLoadSource_Env_WithPrefix(t *testing.T) {
	os.Setenv("MYAPP_HOST", "localhost")
	os.Setenv("MYAPP_PORT", "8080")
	os.Setenv("OTHER_KEY", "ignored")
	t.Cleanup(func() {
		os.Unsetenv("MYAPP_HOST")
		os.Unsetenv("MYAPP_PORT")
		os.Unsetenv("OTHER_KEY")
	})

	src, err := LoadSource(SourceSpec{Type: SourceTypeEnv, Ref: "MYAPP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Vars["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", src.Vars["HOST"])
	}
	if src.Vars["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", src.Vars["PORT"])
	}
	if _, ok := src.Vars["OTHER_KEY"]; ok {
		t.Error("OTHER_KEY should have been filtered out")
	}
}

func TestLoadSource_UnknownType(t *testing.T) {
	_, err := LoadSource(SourceSpec{Type: "s3", Ref: "bucket/path"})
	if err == nil {
		t.Fatal("expected error for unknown source type")
	}
}
