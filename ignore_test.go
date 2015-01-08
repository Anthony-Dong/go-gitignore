// Implement tests for the `ignore` library
package ignore

import (
    "fmt"

    "testing"
    "github.com/stretchr/testify/assert"
)

// Validate `CompileIgnoreLines()`
func TestCompileIgnoreLines(test *testing.T) {
    fmt.Println(" * Testing CompileIgnoreLines()")

    lines := []string{"abc/def", "a/b/c", "b"}
    object2, error := CompileIgnoreLines(lines...)

    // Validate no error
    assert.Nil(test, error, "error from CompileIgnoreLines should be nil")

    // IncludesPath
    // Paths which should not be ignored
    assert.Equal(test, true, object2.IncludesPath("abc"), "abc should be accepted")
    assert.Equal(test, true, object2.IncludesPath("def"), "def should be accepted")
    assert.Equal(test, true, object2.IncludesPath("bd"),  "bd should be accepted")

    // Paths which should be ignored
    assert.Equal(test, false, object2.IncludesPath("abc/def/child"), "abc/def/child should be rejected")
    assert.Equal(test, false, object2.IncludesPath("a/b/c/d"),       "a/b/c/d should be rejected")

    // IgnorePath
    assert.Equal(test, false, object2.IgnoresPath("abc"), "abc should not be ignored")
    assert.Equal(test, false, object2.IgnoresPath("def"), "def should not be ignored")
    assert.Equal(test, false, object2.IgnoresPath("bd"),  "bd should not be ignored")

    // Paths which should be ignored
    assert.Equal(test, true, object2.IgnoresPath("abc/def/child"), "abc/def/child should be ignored")
    assert.Equal(test, true, object2.IgnoresPath("a/b/c/d"),       "a/b/c/d should be ignored")
}

// Validate `CompileIgnoreFile()` for invalid files
func TestCompileIgnoreFile_InvalidFile(test *testing.T) {
    fmt.Println(" * Testing CompileIgnoreFile() for invalid file")

    object, error := CompileIgnoreFile("./test_fixtures/invalid.file")
    assert.Nil(test, object, "object should be nil")
    assert.NotNil(test, error, "error should be unknown file / dir")
}

// Validate `CompileIgnoreFile()` for an empty files
func TestCompileIgnoreLines_EmptyFile(test *testing.T) {
    fmt.Println(" * Testing CompileIgnoreFile() for empty file")

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore.0")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, false, object.IgnoresPath("a"),       "should accept all paths")
    assert.Equal(test, false, object.IgnoresPath("a/b"),     "should accept all paths")
    assert.Equal(test, false, object.IgnoresPath(".foobar"), "should accept all paths")
}

// Validate `CompileIgnoreFile()` 1
func TestCompileIgnoreLines_1(test *testing.T) {
    fmt.Println(" * Testing CompileIgnoreFile() exclude all but...")

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore.1")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, false, object.IgnoresPath("a"),        "a should not be ignored")
    //assert.Equal(test, false, object.IgnoresPath("/foo"),     "/foo should not match")
    //assert.Equal(test, false, object.IgnoresPath("/foo/baz"), "/foo/baz should not match")
    assert.Equal(test, true,  object.IgnoresPath("/foo/bar"), "/foo/bar should be ignored")
}

// Validate `CompileIgnoreFile()` 2
func TestCompileIgnoreLines_2(test *testing.T) {
    fmt.Println(" * Testing CompileIgnoreFile() file space stuffs")

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore.2")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, 1, len(object.patterns), "should have single regex pattern")
    assert.Equal(test, false, object.IgnoresPath("/abc/abc"), "/abc/abc should not be ignored")
    //assert.Equal(test, true,  object.IgnoresPath("/abc/def"), "/abc/def should be ignored")
}
