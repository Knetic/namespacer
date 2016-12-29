package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"io"
	"io/ioutil"
	"path/filepath"
)

func main() {

	var settings RunSettings
	var err error

	settings, err = parseRunSettings()
	if(err != nil) {
		fatal(err, 1)
	}

	err = filepath.Walk(settings.targetPath, createWalker(settings))
	if err != nil {
		fatal(err, 1)
	}
}

// creates a file walker that rewrites the C# namespace.
func createWalker(settings RunSettings) func(path string, info os.FileInfo, err error) (error) {

	return func(path string, info os.FileInfo, _ error) error {

		var ext string
		var err error

		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		if info.Mode().IsDir() {
			return nil
		}

		ext = filepath.Ext(path)
		if ext != ".cs" {
			return nil
		}

		// well it's definitely a C# file. modify it.
		return ensureNamespace(path, settings.namespace)
	}
}

// opens the file at the given path and makes sure it has the expected namespace.
func ensureNamespace(path string, namespace string) error {

	var target, temp *os.File
	var err error
	var lineNumber int
	var hasNamespace bool

	// already contains. deal with it.
	hasNamespace, lineNumber, err = containsNamespace(path, namespace)
	if err != nil {
		return err
	}
	if hasNamespace	{
		return nil
	}

	target, err = os.Open(path)
	if err != nil {
		return err
	}
	defer target.Close()

	temp, err = ioutil.TempFile(os.TempDir(), "namespacer")
	if err != nil {
		return err
	}
	defer temp.Close()

	err = insertLine(target, temp, ("namespace " + namespace + "\n{\n"), lineNumber)
	if err != nil {
		return err
	}

	// move on over.
	temp.Close()
	target.Close()
	return moveFile(temp.Name(), path)
}

// Returns true if the given file contains the given namespace.
// false otherwise.
// If there is no namespace, this also returns the line number which a namespace should be inserted.
func containsNamespace(path string, namespace string) (bool, int, error) {

	var file *os.File
	var scanner *bufio.Scanner
	var line string
	var err error
	var lineNumber, desiredNamespaceLocation int

	file, err = os.Open(path)
	if err != nil {
		return false, -1, err
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	desiredNamespaceLocation = -1

	for scanner.Scan() {

		line = scanner.Text()
		lineNumber++

		// probably an autodoc comment before the thing we're after.
		if strings.Contains(line, "///") && desiredNamespaceLocation == -1 {
			desiredNamespaceLocation = lineNumber-1
			continue
		}

		// comment.
		if strings.HasPrefix(line, "//") {
			continue
		}

		// already got a namespace?
		if strings.HasPrefix(line, "namespace") {
			return true, -1, nil
		}

		// see if this has a signature
		if !containsSignature(line) {
			continue
		}

		// alright, found something. Call it a day.
		if desiredNamespaceLocation == -1 {
			desiredNamespaceLocation = lineNumber-1
		}

		return false, desiredNamespaceLocation, nil
	}

	return false, -1, nil
}

func containsSignature(line string) bool {

	return strings.Contains(line, "public class ") ||
		strings.Contains(line, "public abstract class ") ||
		strings.Contains(line, "public sealed class ") ||
		strings.Contains(line, "public struct ") ||
		strings.Contains(line, "public enum ") ||
		strings.Contains(line, "public interface ") ||
		strings.Contains(line, "public delegate ") ||
		strings.Contains(line, "internal class ") ||
		strings.Contains(line, "internal abstract class ") ||
		strings.Contains(line, "internal sealed class ") ||
		strings.Contains(line, "internal struct ") ||
		strings.Contains(line, "internal enum ") ||
		strings.Contains(line, "internal interface ") ||
		strings.Contains(line, "internal delegate  ")
}

/*
	copies all of [inFile] to [outFile], line by line.
	except that on the given [lineNumber], the line will be the [desired] string.
*/
func insertLine(inFile *os.File, outFile *os.File, desired string, lineNumber int) error {

	var scanner *bufio.Scanner
	var writer *bufio.Writer
	var line string
	var err error

	scanner = bufio.NewScanner(inFile)
	writer = bufio.NewWriter(outFile)

	// copy up to the line
	for i := 0; i < lineNumber; i++ {

		scanner.Scan()
		err = scanner.Err()
		if err != nil {
			return err
		}

		writeLine(writer, scanner.Text())
	}

	// insert
	writer.WriteString(desired)

	// do the rest, indenting each.
	for scanner.Scan() {

		line = "\t" + scanner.Text()
		writeLine(writer, line)
	}

	writeLine(writer, "}")

	err = scanner.Err()
	if err != nil {
		return err
	}

	writer.Flush()
	return nil
}

func moveFile(sourcePath string, targetPath string) error {

	var err error

	// try to rename first. If they're on different partitions, this'll choke.
	err = os.Rename(sourcePath, targetPath)
	if err == nil {
		return nil
	}

	err = copyFile(sourcePath, targetPath)
	if err != nil {
		return err
	}

	return nil
	//return os.Remove(sourcePath)
}

func copyFile(sourcePath string, targetPath string) error {

	var source, target *os.File
	var err error

	source, err = os.OpenFile(sourcePath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err = os.OpenFile(targetPath, os.O_CREATE | os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	if err != nil {
		return err
	}

	return nil
}

func writeLine(writer *bufio.Writer, line string) {
	writer.WriteString(line + "\n")
}

func fatal(fault error, code int) {

	fmt.Fprintf(os.Stderr, "%s\n", fault.Error())
	os.Exit(code)
}
