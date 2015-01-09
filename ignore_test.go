// Implement tests for the `ignore` library
package ignore

import (
    "os"

    "io/ioutil"
    "path/filepath"

    "fmt"
    "testing"

    "github.com/stretchr/testify/assert"
)

const (
    TEST_DIR = "test_fixtures"
)

// Helper function to setup a test fixture dir and write to
// it a file with the name "fname" and content "content"
func writeFileToTestDir(fname, content string) {
    testDirPath := "." + string(filepath.Separator) + TEST_DIR
    testFilePath := testDirPath + string(filepath.Separator) + fname

    _ = os.MkdirAll(testDirPath, 0755)
    _ = ioutil.WriteFile(testFilePath, []byte(content), os.ModePerm)
}

func cleanupTestDir() {
    _ = os.RemoveAll(fmt.Sprintf(".%s%s", string(filepath.Separator), TEST_DIR))
}

// Validate "CompileIgnoreLines()"
func TestCompileIgnoreLines(test *testing.T) {
    lines := []string{"abc/def", "a/b/c", "b"}

    object2, error := CompileIgnoreLines(lines...)
    assert.Nil(test, error, "error from CompileIgnoreLines should be nil")

    // IncludesPath
    // Paths which should not be ignored
    assert.Equal(test, true, object2.IncludesPath("abc"), "abc should not be ignored")
    assert.Equal(test, true, object2.IncludesPath("def"), "def should not be ignored")
    assert.Equal(test, true, object2.IncludesPath("bd"),  "bd should not be ignored")
    // Paths which should be ignored
    assert.Equal(test, false, object2.IncludesPath("abc/def/child"), "abc/def/child should be ignored")
    assert.Equal(test, false, object2.IncludesPath("a/b/c/d"),       "a/b/c/d should be ignored")

    // IgnorePath
    assert.Equal(test, false, object2.IgnoresPath("abc"), "abc should not be ignored")
    assert.Equal(test, false, object2.IgnoresPath("def"), "def should not be ignored")
    assert.Equal(test, false, object2.IgnoresPath("bd"),  "bd should not be ignored")

    // Paths which should be ignored
    assert.Equal(test, true, object2.IgnoresPath("abc/def/child"), "abc/def/child should be ignored")
    assert.Equal(test, true, object2.IgnoresPath("a/b/c/d"),       "a/b/c/d should be ignored")
}

// Validate "CompileIgnoreFile()" for invalid files
func TestCompileIgnoreFile_InvalidFile(test *testing.T) {
    object, error := CompileIgnoreFile("./test_fixtures/invalid.file")
    assert.Nil(test, object, "object should be nil")
    assert.NotNil(test, error, "error should be unknown file / dir")
}

// Validate "CompileIgnoreFile()" for an empty files
func TestCompileIgnoreLines_EmptyFile(test *testing.T) {
    writeFileToTestDir("test.gitignore", ``)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, false, object.IgnoresPath("a"),       "should accept all paths")
    assert.Equal(test, false, object.IgnoresPath("a/b"),     "should accept all paths")
    assert.Equal(test, false, object.IgnoresPath(".foobar"), "should accept all paths")
}

//
// FOLDER based path checking tests
//

// Validate "CompileIgnoreFile()" for correct handling of the negation operator "!"
func TestCompileIgnoreLines_HandleIncludePattern(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
# exclude everything except directory foo/bar
/*
!/foo
/foo/*
!/foo/bar`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    //assert.Equal(test, false, object.IgnoresPath("a"),        "a should not be ignored")
    //assert.Equal(test, false, object.IgnoresPath("/foo"),     "/foo should not match")
    //assert.Equal(test, false, object.IgnoresPath("/foo/baz"), "/foo/baz should not match")
    assert.Equal(test, true,  object.IgnoresPath("/foo/bar"), "/foo/bar should be ignored")
}

// Validate "CompileIgnoreFile()" for correct handling of comments and empty lines
func TestCompileIgnoreLines_HandleSpaces(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
#
# A comment

# Another comment


          # Invalid Comment

abc/def
`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, 2, len(object.patterns), "should have two regex pattern")
    assert.Equal(test, false, object.IgnoresPath("abc/abc"), "/abc/abc should not be ignored")
    assert.Equal(test, true,  object.IgnoresPath("abc/def"), "/abc/def should be ignored")
}

// Validate "CompileIgnoreFile()" for correct handling of leading / chars
func TestCompileIgnoreLines_HandleLeadingSlash(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
/a/b/c
d/e/f
/g
`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, 3, len(object.patterns), "should have 3 regex patterns")
    assert.Equal(test, true,  object.IgnoresPath("a/b/c"),   "a/b/c should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("a/b/c/d"), "a/b/c/d should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("d/e/f"),   "d/e/f should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("g"),       "g should be ignored")
}

//
// FILE based path checking tests
//

// Validate "CompileIgnoreFile()" for correct handling of files starting with # or !
func TestCompileIgnoreLines_HandleLeadingSpecialChars(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
# Comment
\#file.txt
\!file.txt
file.txt
`)
    defer cleanupTestDir()

    object, _ := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, true,  object.IgnoresPath("#file.txt"),  "#file.txt should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("!file.txt"),  "!file.txt should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("file.txt"),   "file.txt should be ignored")
    assert.Equal(test, false, object.IgnoresPath("file2.txt"),  "file2.txt should not be ignored")
    assert.Equal(test, false, object.IgnoresPath("a/file.txt"), "a/file.txt should not be ignored")
}
