package zip

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func Archive(inPath string, outPath string, includeRootDir bool) error {
	inFileInfo, err := os.Stat(inPath)
	if err != nil {
		return err
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)

	inIsDir := inFileInfo.IsDir()
	archivePath := ""
	if includeRootDir {
		archivePath = inFileInfo.Name()
	}

	err = archive(zipWriter, inPath, inIsDir, archivePath)
	if err != nil {
		return err
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

func archive(zipWriter *zip.Writer, inPath string, inIsDir bool, archivePath string) error {
	if inIsDir {
		return archiveDir(zipWriter, inPath, archivePath)
	} else {
		return archiveFile(zipWriter, inPath, archivePath)
	}
}

func archiveDir(zipWriter *zip.Writer, inPath string, archivePath string) error {
	childFileInfos, err := ioutil.ReadDir(inPath)
	if err != nil {
		return err
	}

	for _, childFileInfo := range childFileInfos {
		childFileName := childFileInfo.Name()
		childInPath := filepath.Join(inPath, childFileName)
		childArchivePath := path.Join(archivePath, childFileName)
		childIsDir := childFileInfo.IsDir()
		err = archive(zipWriter, childInPath, childIsDir, childArchivePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func archiveFile(zipWriter *zip.Writer, inPath string, archivePath string) error {
	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()
	writer, err := zipWriter.Create(archivePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, inFile)
	if err != nil {
		return err
	}
	return nil
}

func Unarchive(archivePath string, filePath string) error {
	//TODO
	return nil
}
