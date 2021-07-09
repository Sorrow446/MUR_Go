package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	jar, _ = cookiejar.New(nil)
	client = &http.Client{Jar: jar, Transport: &myTransport{}}
)

const (
	outFolder = "MUR downloads"
	tmpFolder = "MUR_tmp"
)

type myTransport struct{}

func (t *myTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"+
		" (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	)
	return http.DefaultTransport.RoundTrip(req)
}

func getScriptDir() (string, error) {
	var (
		ok    bool
		err   error
		fname string
	)
	if filepath.IsAbs(os.Args[0]) {
		_, fname, _, ok = runtime.Caller(0)
		if !ok {
			return "", errors.New("Failed to get script filename.")
		}
	} else {
		fname, err = os.Executable()
		if err != nil {
			return "", err
		}
	}
	scriptDir := filepath.Dir(fname)
	return scriptDir, nil
}

func readTxtFile(path string) ([]string, error) {
	var lines []string
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return lines, nil
}

func contains(lines []string, value string) bool {
	for _, line := range lines {
		if strings.EqualFold(line, value) {
			return true
		}
	}
	return false
}

func processUrls(urls []string) ([]string, error) {
	var processed []string
	var txtPaths []string
	for _, url := range urls {
		if strings.HasSuffix(url, ".txt") && !contains(txtPaths, url) {
			txtLines, err := readTxtFile(url)
			if err != nil {
				return nil, err
			}
			for _, txtLine := range txtLines {
				if !contains(processed, txtLine) {
					processed = append(processed, txtLine)
				}
			}
			txtPaths = append(txtPaths, url)
		} else {
			if !contains(processed, url) {
				processed = append(processed, url)
			}
		}
	}
	return processed, nil
}

func parseCookies() ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	replacer := strings.NewReplacer("&quot;", "", "\"", "")
	f, err := os.Open("cookies.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		splitLine := strings.Split(line, "\t")
		secure, err := strconv.ParseBool(splitLine[3])
		if err != nil {
			return nil, err
		}
		cookie := &http.Cookie{
			Domain: ".marvel.com",
			Name:   splitLine[5],
			Path:   "/",
			Secure: secure,
			Value:  replacer.Replace(splitLine[6]),
		}
		cookies = append(cookies, cookie)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return cookies, nil
}

func setCookies(cookies []*http.Cookie) {
	urlObj, _ := url.Parse("https://www.marvel.com/")
	client.Jar.SetCookies(urlObj, cookies)
}

func checkUrl(url string) bool {
	regexes := [2]string{
		`^https://www.marvel.com/comics/issue/\d+/[\d\w-]+$`,
		`^https://read.marvel.com/#/book/\d+$`,
	}
	for _, regexString := range regexes {
		regex := regexp.MustCompile(regexString)
		match := regex.MatchString(url)
		if match {
			return true
		}
	}
	return false
}

func getId(url string) (string, error) {
	req, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return "", errors.New(req.Status)
	}
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	regex := regexp.MustCompile(`digital_comic_id: "(\d+)"`)
	match := regex.FindStringSubmatch(bodyString)
	if match != nil {
		return match[1], nil
	}
	return "", errors.New("No regex match.")
}

func randInt() string {
	i := 10000 + rand.Intn(99999-9999)
	return strconv.Itoa(i)
}

func getMeta(id string) (*Meta, error) {
	req, err := client.Get(
		"https://bifrost.marvel.com/v1/catalog/digital-comics/metadata/" + id,
	)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return nil, errors.New(req.Status)
	}
	var obj Meta
	err = json.NewDecoder(req.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}
	if obj.Code != 200 {
		return nil, errors.New("Bad response.")
	}
	return &obj, nil
}

func getAssetMeta(id string) (*AssetMeta, error) {
	req, err := http.NewRequest(
		"GET", "https://bifrost.marvel.com/v1/catalog/digital-comics/web/assets/"+id, nil,
	)
	if err != nil {
		return nil, err
	}
	query := url.Values{}
	query.Set("rand", randInt())
	req.URL.RawQuery = query.Encode()
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(do.Status)
	}
	var obj AssetMeta
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}
	if obj.Code != 200 {
		return nil, errors.New("Bad response.")
	}
	if !obj.Data.Results[0].AuthState.Subscriber {
		panic("An account subscription is required.")
	}
	return &obj, nil
}

func sanitize(filename string) string {
	regex := regexp.MustCompile(`[\/:*?"><|]`)
	sanitized := regex.ReplaceAllString(filename, "_")
	return sanitized
}

func fileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func downloadPages(meta *AssetMeta) (int, error) {
	pageCount := len(meta.Data.Results[0].Pages)
	for num, page := range meta.Data.Results[0].Pages {
		num++
		fmt.Printf("\rPage %d of %d.", num, pageCount)
		pageOutPath := filepath.Join(tmpFolder, fmt.Sprintf("%04d.jpg", num))
		f, err := os.OpenFile(pageOutPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return -1, err
		}
		defer f.Close()
		req, err := client.Get(page.Assets.Source)
		if err != nil {
			return -1, err
		}
		defer req.Body.Close()
		_, err = io.Copy(f, req.Body)
		if err != nil {
			return -1, err
		}
	}
	return pageCount, nil
}

func createCbz(pageCount int, outPath string) error {
	newZipFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	for i := 1; i <= pageCount; i++ {
		pagePath := filepath.Join(tmpFolder, fmt.Sprintf("%04d.jpg", i))
		err := addImagetoZip(zipWriter, pagePath)
		if err != nil {
			return err
		}
		_ = os.Remove(pagePath)
	}
	return nil
}
func addImagetoZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filename)
	header.Method = zip.Store
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func init() {
	fmt.Println(`
 _____ _____ _____ 
|     |  |  | __  |
| | | |  |  |    -|
|_|_|_|_____|__|__|
	`)
	if len(os.Args) == 1 {
		fmt.Println("At least one URL or text file filename/path is required.")
		os.Exit(1)
	}
	scriptDir, err := getScriptDir()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(scriptDir)
	if err != nil {
		panic(err)
	}
	for _, path := range [2]string{outFolder, tmpFolder} {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil && !errors.Is(err, os.ErrExist) {
			panic(err)
		}
	}
	cookies, err := parseCookies()
	if err != nil {
		panic(err)
	}
	setCookies(cookies)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	urls, err := processUrls(os.Args[1:])
	if err != nil {
		errString := fmt.Sprintf("Failed to process URLs. %s", err)
		panic(errString)
	}
	total := len(urls)
	for num, url := range urls {
		fmt.Printf("URL %d of %d:\n", num+1, total)
		ok := checkUrl(url)
		if !ok {
			fmt.Println("Invalid URL:", url)
			continue
		}
		id, err := getId(url)
		if err != nil {
			fmt.Println("Failed to extract comic ID.", err)
			continue
		}
		meta, err := getMeta(id)
		if err != nil {
			fmt.Println("Failed to get comic metadata.", err)
			continue
		}
		title := meta.Data.Results[0].IssueMeta.Title
		fmt.Println(title)
		outPath := filepath.Join(outFolder, sanitize(title)+".cbz")
		exists, err := fileExists(outPath)
		if err != nil {
			fmt.Println("Failed to check if comic already exists locally.", err)
			continue
		}
		if exists {
			fmt.Println("Comic already exists locally.")
			continue
		}
		assetMeta, err := getAssetMeta(id)
		if err != nil {
			fmt.Println("Failed to get asset metadata.", err)
			continue
		}
		pageCount, err := downloadPages(assetMeta)
		if err != nil {
			fmt.Println("Failed to download pages.", err)
			continue
		}
		err = createCbz(pageCount, outPath)
		if err != nil {
			fmt.Println("Failed to create CBZ.", err)
		}
	}
}
