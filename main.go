package main

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func main() {


	regexSemVer         := `(\d+\.)?(\d+\.)?(\*|\d+)`
	regexCurrentVersion := `v`+regexSemVer+`$`
	regexNextVersion    := regexSemVer+`[.]\s`

	fmt.Println("Operative System -> "+runtime.GOOS)
	fmt.Println("Architecture -> "+runtime.GOARCH)

	terraformFolder, errFin := findTerraform(runtime.GOOS)

	if errFin != nil {
		return
	}

	basePath := "https://releases.hashicorp.com/terraform/"

	// https://releases.hashicorp.com/terraform/0.12.25/terraform_0.12.25_windows_amd64.zip
	// https://releases.hashicorp.com/terraform/0.12.25/terraform_0.12.25_linux_amd64.zip

	fmt.Println("Terraform is available -> "+strconv.FormatBool(isCommandAvailable("terraform")))

	cmd := exec.Command("terraform", "version")

	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}

	var nextVer bool
	var version []string

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			//fmt.Printf("\t > %s\n", scanner.Text())

			semVer := regexp.MustCompile(regexSemVer)

			r := regexp.MustCompile(regexCurrentVersion)
			res := r.MatchString(scanner.Text())

			if res == true {
				match := semVer.FindStringSubmatch(scanner.Text())
				fmt.Println("current version -> "+match[0])
			}

			rNext := regexp.MustCompile(regexNextVersion)
			nextVer = rNext.MatchString(scanner.Text())

			if nextVer == true {
				version = semVer.FindStringSubmatch(scanner.Text())
				fmt.Println("new version -> "+version[0])
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return
	}

	if nextVer == true {
		fmt.Println("Downloading newer version of Terraform")

		url := basePath+version[0]+"/terraform_"+version[0]+"_"+runtime.GOOS+"_"+runtime.GOARCH+".zip"

		errDown := DownloadFile("terraform.zip", url)
		//https://releases.hashicorp.com/terraform/0.12.25/terraform_0.12.25_linux_amd64.zip

		if errDown != nil {
			fmt.Println("errore durante download")
		}

		here,_ := filepath.Abs("terraform.zip")
		file := filepath.Base(here)
		dir, file := filepath.Split(here)
		fmt.Println("Dir:", dir)   //Dir: /some/path/to/remove/
		fmt.Println("File:", file)

		files, unzipErr := Unzip("terraform.zip", dir)

		if unzipErr != nil {
			fmt.Println(unzipErr.Error())
			fmt.Println("errore durane unzip")
		}

		defer os.Remove("terraform.zip")
		defer os.Remove(files[0])


		// sposto il file
		oldLocation, newExecutable := filepath.Split(files[0])
		fmt.Println("Old Location:", oldLocation)   //Dir: /some/path/to/remove/
		fmt.Println("Executable:", newExecutable)


		newLocation, oldExecutable := filepath.Split(terraformFolder)
		fmt.Println("New Location:", newLocation)   //Dir: /some/path/to/remove/
		fmt.Println("Old Executable:", oldExecutable)


		fmt.Println("")

		terraformFolder = strings.Trim(terraformFolder, " ")
		terraformFolder = strings.TrimSuffix(terraformFolder, "\n")
		terraformFolder = strings.TrimSuffix(terraformFolder, "\r")
		fmt.Println("Sposto "+files[0]+" in "+terraformFolder)
		err := MoveFile(files[0], terraformFolder)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println()
	}
}



func isCommandAvailable(name string) bool {
	cmd := exec.Command(name, "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}



// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}


func findTerraform(os string) (string, error) {

	location := ""

	if os == "linux" {
		out, err := exec.Command("whereis","terraform").Output()

		if err != nil {
			return location, err
		}

		fmt.Println(out)

		return location, nil

	} else if os == "windows" {

		out, err := exec.Command("where","terraform").Output()

		if err != nil {
			return location, err
		}

		fmt.Println("terraform installation --> "+string(out))

		return string(out), nil

	} else {
		return location, errors.New("not supported")
	}


}

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}


func MoveFile(source, destination string) (err error) {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	fi, err := src.Stat()
	if err != nil {
		return err
	}
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := fi.Mode() & os.ModePerm
	dst, err := os.OpenFile(destination, flag, perm)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		dst.Close()
		os.Remove(destination)
		return err
	}
	err = dst.Close()
	if err != nil {
		return err
	}
	err = src.Close()
	if err != nil {
		return err
	}
	err = os.Remove(source)
	if err != nil {
		return err
	}
	return nil
}