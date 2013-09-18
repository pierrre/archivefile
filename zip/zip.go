package zip

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func Archive(filePath string, includeRootDir bool, writer io.Writer) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(writer)

	isDir := fileInfo.IsDir()
	archivePath := ""
	if !isDir || includeRootDir {
		archivePath = fileInfo.Name()
	}

	err = archive(zipWriter, filePath, isDir, archivePath)
	if err != nil {
		return err
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

func ArchiveFile(filePath string, includeRootDir bool, outFilePath string) error {
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = Archive(filePath, includeRootDir, outFile)
	if err != nil {
		return err
	}

	return nil
}

func archive(zipWriter *zip.Writer, filePath string, isDir bool, archivePath string) error {
	if isDir {
		return archiveDir(zipWriter, filePath, archivePath)
	} else {
		return archiveFile(zipWriter, filePath, archivePath)
	}
}

func archiveDir(zipWriter *zip.Writer, filePath string, archivePath string) error {
	childFileInfos, err := ioutil.ReadDir(filePath)
	if err != nil {
		return err
	}

	for _, childFileInfo := range childFileInfos {
		childFileName := childFileInfo.Name()
		childFilePath := filepath.Join(filePath, childFileName)
		childArchivePath := path.Join(archivePath, childFileName)
		childIsDir := childFileInfo.IsDir()
		err = archive(zipWriter, childFilePath, childIsDir, childArchivePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func archiveFile(zipWriter *zip.Writer, filePath string, archivePath string) error {
	file, err := os.Open(filePath)
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

func Unarchive(archivePath string, filePath string) error {
	//TODO
	return nil
}
