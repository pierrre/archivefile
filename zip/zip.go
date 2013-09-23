// (Un)archive file/directory to/from file/writer/reader using "archive/zip" package
package zip

import (
	zip_impl "archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Archive a file/directory to a writer
//
// If inFilePath is a file, the archive will contain this file at the root.
// If inFilePath is a directory, the archive will contain the directory's content if includeRootDir is false, or the directory if includeRootDir is true.
//
// If progress is not nil, it receives the path of each file added to the archive.
func Archive(inFilePath string, includeRootDir bool, writer io.Writer, progress ProgressFunc) error {
	fileInfo, err := os.Stat(inFilePath)
	if err != nil {
		return err
	}

	zipWriter := zip_impl.NewWriter(writer)

	isDir := fileInfo.IsDir()
	archivePath := ""
	if !isDir || includeRootDir {
		archivePath = fileInfo.Name()
	}

	err = archive(zipWriter, inFilePath, isDir, archivePath, progress)
	if err != nil {
		return err
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

// Archive a file/directory to a file
//
// See Archive() doc
func ArchiveFile(inFilePath string, includeRootDir bool, outFilePath string, progress ProgressFunc) error {
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = Archive(inFilePath, includeRootDir, outFile, progress)
	if err != nil {
		return err
	}

	return nil
}

func archive(zipWriter *zip_impl.Writer, inFilePath string, isDir bool, archivePath string, progress ProgressFunc) error {
	if isDir {
		return archiveDir(zipWriter, inFilePath, archivePath, progress)
	} else {
		return archiveFile(zipWriter, inFilePath, archivePath, progress)
	}
}

func archiveDir(zipWriter *zip_impl.Writer, inFilePath string, archivePath string, progress ProgressFunc) error {
	childFileInfos, err := ioutil.ReadDir(inFilePath)
	if err != nil {
		return err
	}

	for _, childFileInfo := range childFileInfos {
		childFileName := childFileInfo.Name()
		childFilePath := filepath.Join(inFilePath, childFileName)
		childArchivePath := path.Join(archivePath, childFileName)
		childIsDir := childFileInfo.IsDir()
		err = archive(zipWriter, childFilePath, childIsDir, childArchivePath, progress)
		if err != nil {
			return err
		}
	}

	return nil
}

func archiveFile(zipWriter *zip_impl.Writer, inFilePath string, archivePath string, progress ProgressFunc) error {
	if progress != nil {
		progress(archivePath)
	}

	file, err := os.Open(inFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(archivePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}

// Unarchive a reader to a directory
//
// The data's size is required because the zip reader needs it.
//
// The archive's content will be extracted directly to outFilePath.
//
// If progress is not nil, it receives the path of each file extracted from the archive.
func Unarchive(reader io.ReaderAt, readerSize int64, outFilePath string, progress ProgressFunc) error {
	zipReader, err := zip_impl.NewReader(reader, readerSize)
	if err != nil {
		return err
	}

	for _, zipFile := range zipReader.File {
		err := unarchiveFile(zipFile, outFilePath, progress)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unarchive a file to a directory
//
// See Unarchive() doc
func UnarchiveFile(inFilePath string, outFilePath string, progress ProgressFunc) error {
	inFile, err := os.Open(inFilePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	inFileInfo, err := inFile.Stat()
	if err != nil {
		return err
	}
	inFileSize := inFileInfo.Size()

	err = Unarchive(inFile, inFileSize, outFilePath, progress)
	if err != nil {
		return err
	}

	return nil
}

func unarchiveFile(zipFile *zip_impl.File, outFilePath string, progress ProgressFunc) error {
	if zipFile.FileInfo().IsDir() {
		return nil
	}

	if progress != nil {
		progress(zipFile.Name)
	}

	reader, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	filePath := filepath.Join(outFilePath, filepath.Join(strings.Split(zipFile.Name, "/")...))

	err = os.MkdirAll(filepath.Dir(filePath), os.FileMode(0755))
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

type ProgressFunc func(filePath string)
