package core

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Zip(source, target string) error {
	file, err := os.Create(target)
	if err != nil {
		return err
	}

	writer := zip.NewWriter(file)
	defer writer.Close()

	if err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filePath := strings.Replace(path, source, "", 1)
		if filePath == "" {
			return nil
		}

		if strings.HasPrefix(filePath, "/backup") {
			return nil
		}

		log.Println("Compressing ", filePath)

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filePath

		if info.IsDir() {
			header.Name += "/"
		}

		w, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		_, err = io.Copy(w, srcFile)

		return err
	}); err != nil {
		return err
	}

	if _, err := os.Stat(Conf.DataDir + "/backup/database.sql"); err == nil {
		if err := addFileToZip(Conf.DataDir+"/backup/database.sql", "database.sql", writer); err != nil {
			return err
		}
	}

	return addFileToZip("config.json", "config.json", writer)
}

func addFileToZip(file, targetName string, writer *zip.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Println("Compressing ", targetName)

	w, err := writer.Create(targetName)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)

	return err
}
