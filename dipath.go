package dipath

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// Project 함수는 경로를 받아서 프로젝트 이름을 반환한다.
func Project(path string) (string, error) {
	if path == "" {
		return path, errors.New("빈 문자열 입니다")
	}
	p := strings.Replace(path, "\\", "/", -1)
	regRule := `/show[/_](\S[^/]+)`
	if strings.HasPrefix(p, "/backup/") {
		regRule = `/backup/\d+?/(\S[^/]+)`
	}
	re, err := regexp.Compile(regRule)
	if err != nil {
		return "", err
	}

	results := re.FindStringSubmatch(p)
	if results == nil {
		return "", errors.New(path + " 경로에서 프로젝트 정보를 가지고 올 수 없습니다.")
	}
	return results[len(results)-1], nil
}

//Projectlist 함수는 프로젝트경로의 폴더를 문자열 리스트로 가지고 온다.
func Projectlist() []string {
	var dirlist []string
	projectpath := "/show"
	files, _ := ioutil.ReadDir(projectpath)
	for _, f := range files {
		fileInfo, _ := os.Lstat(projectpath + "/" + f.Name())
		if !strings.HasPrefix(f.Name(), ".") && fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			dirlist = append(dirlist, f.Name())
		} else if !strings.HasPrefix(f.Name(), ".") && fileInfo.IsDir() {
			dirlist = append(dirlist, f.Name())
		}
	}
	sort.Strings(dirlist)
	return dirlist
}

//TEMP 함수는 서버의 temp경로를 반환한다.
func TEMP() string {
	switch runtime.GOOS {
	case "windows":
		return "\\\\10.0.200.100\\show_TEMP\\tmp\\"
	case "linux":
		return "/show/TEMP/tmp/"
	case "darwin":
		return "/show/TEMP/tmp/"
	default:
		return "/show/TEMP/tmp/"
	}
}

// Win2lin 함수는 윈도우즈 경로를 리눅스 경로로 바꾼다. 만약, 변환되지 않으면 패스를 그대로 출력한다.
func Win2lin(path string) string {
	if strings.HasPrefix(path, "W:\\") {
		return "/show/" + strings.Replace(path[3:], "\\", "/", len(path[3:]))
	} else if strings.HasPrefix(path, "/show/") {
		return path
	} else if strings.HasPrefix(path, "/lustre") { // lustre, lustre2, lustre3, lustre4 로 시작할 때..
		return path
	} else if strings.HasPrefix(path, "\\\\10.0.200.100\\show_") {
		return "/show/" + strings.Replace(path[20:], "\\", "/", len(path[20:]))
	} else if strings.HasPrefix(path, "\\\\10.0.200.100\\lustre_Digitalidea_source\\") {
		return "/lustre2/Digitalidea_source/" + strings.Replace(path[41:], "\\", "/", len(path[41:]))
	} else {
		return path
	}
}

//Lin2win 함수는 리눅스 경로를 윈도우즈 경로로 바꾼다.
func Lin2win(path string) string {
	if strings.HasPrefix(path, "/lustre2/Digitalidea_source/flib") { //flib
		return "\\\\10.0.200.100\\lustre_Digitalidea_source\\flib" + strings.Replace(path[32:], "/", "\\", len(path[32:]))
	} else if strings.HasPrefix(path, "/lustre/Digitalidea_source/flib") { //flib
		return "\\\\10.0.200.100\\lustre_Digitalidea_source\\flib" + strings.Replace(path[31:], "/", "\\", len(path[31:]))
	} else if strings.HasPrefix(path, "/show") {
		return "\\\\10.0.200.100\\show_" + strings.Replace(path[6:], "/", "\\", len(path[6:]))
	} else if strings.HasPrefix(path, "/lustre/show") {
		return "\\\\10.0.200.100\\show_" + strings.Replace(path[13:], "/", "\\", len(path[13:]))
	} else if strings.HasPrefix(path, "/lustre2/show") {
		return "\\\\10.0.200.100\\show_" + strings.Replace(path[14:], "/", "\\", len(path[14:]))
	} else if strings.HasPrefix(path, "/lustre3/show") {
		return "\\\\10.0.200.100\\show_" + strings.Replace(path[14:], "/", "\\", len(path[14:]))
	} else if strings.HasPrefix(path, "/lustre4/show") {
		return "\\\\10.0.200.100\\show_" + strings.Replace(path[14:], "/", "\\", len(path[14:]))
	} else {
		return path
	}
}

// RmProtocol 함수는 웹에서 파일을 드레그시 붙는 file:// 형태의 프로토콜 문자열을 제거한다.
func RmProtocol(path string) string {
	prefix := []string{"file://", "http://", "ftp://"}
	for _, p := range prefix {
		if strings.HasPrefix(path, p) {
			return path[len(p):]
		}
	}
	return path
}

//Seqnum 함수는 파일명을 받아서 시퀀스넘버를 반환한다. 만약 리턴할 시컨스넘버가 없으면 -1과 에러를 반환한다.
func Seqnum(path string) (int, error) {
	re, err := regexp.Compile("([0-9]+)\\.[a-zA-Z]+$")
	if err != nil {
		return -1, errors.New("정규 표현식이 잘못되었습니다")
	}

	//예를 들어 "SS_0010_comp_v01.0001.jpg"값이 들어오면
	//results리스트는 다음값을 가집니다. [0]:"0001.jpg", [1]:"0001"
	results := re.FindStringSubmatch(path)
	if results == nil {
		return -1, errors.New("시퀀스 파일이 아닙니다")
	}
	seq := results[1]
	seqNum, err := strconv.Atoi(seq)
	if err != nil {
		return -1, errors.New("시퀀스 파일이 아닙니다")
	}
	return seqNum, nil
}

// Vernum 함수는 파일을 받아서 파일 버젼과 서브버전을 반환한다. 만약 리턴할 버전과 서브버전이  없으면 -1과 에러를 반환한다.
func Vernum(path string) (int, int, error) {
	re, err := regexp.Compile(`_[vV]([0-9]+)(_[wW]([0-9]+))*`)
	if err != nil {
		return -1, -1, errors.New("레귤러 익스프레션이 잘못되었습니다")
	}

	//예를 들어 "S_0010_ani_v01_w02.mb"값이 들어오면
	//results리스트는 다음값을 가집니다. [0]:"v01_w02.mb", [1]:"01", [3]:"02"
	results := re.FindStringSubmatch(path)
	if results == nil {
		return -1, -1, errors.New("버전 정보를 가지고 올 수 없습니다")
	}
	verNum, err := strconv.Atoi(results[1])
	if err != nil {
		return -1, -1, errors.New("버전 정보를 가지고 올 수 없습니다")
	}
	//버전은 값이 있고 서브버전에 값이 없다면 -1을 반환
	subNum, err := strconv.Atoi(results[3])
	if err != nil {
		subNum = -1
	}

	return verNum, subNum, nil
}

// Ideapath 함수는 입력받은 경로가 소유자, 그룹이 0775권한을 가지도록 설정한다.
// 이 권한은 전 사원이 읽고 쓸 수 있는 권한을 가지게된다.
func Ideapath(path string) error {
	err := os.Chmod(path, 0775)
	if err != nil {
		return err
	}
	err = os.Chown(path, 500, 500) // 회사에서 사용하는 공용소유자, 공용그룹
	if err != nil {
		return err
	}
	return nil
}

// Safepath 함수는 입력받은 경로가 소유자,그룹,손님이 555권한을 가지도록 설정한다.
// 이 권한은 전 사원이 읽고 실행만 가능하다. 삼바서버에서 마우스 드레그사고를 방지한다.
// 회사는 주요 상위 경로를 이 권한으로 설정하고 폴더권한을 보호한다. 예) shot 상위폴더
func Safepath(path string) error {
	err := os.Chmod(path, 0555)
	if err != nil {
		return err
	}
	err = os.Chown(path, 500, 500) // 회사에서 사용하는 공용소유자, 공용그룹
	if err != nil {
		return err
	}
	return nil
}

// Seq 함수는 경로를 받아서 시퀀스를 반환한다.
func Seq(path string) (string, error) {
	if path == "" {
		return path, errors.New("빈 문자열 입니다")
	}
	p := strings.Replace(path, "\\", "/", -1)
	regRule := `/show[/_]\S+?/seq/(\S[^/]+)`
	if strings.HasPrefix(p, "/backup/") {
		regRule = `/backup/\d+?/\S+?/\S+?/seq/(\S[^/]+)`
	}
	re, err := regexp.Compile(regRule)
	if err != nil {
		return "", err
	}
	results := re.FindStringSubmatch(p)
	if results == nil {
		return "", errors.New(path + " 경로에서 시퀀스 정보를 가지고 올 수 없습니다.")
	}
	return results[len(results)-1], nil
}

// Shot 함수는 경로를 받아서 샷을 반환한다.
func Shot(path string) (string, error) {
	if path == "" {
		return path, errors.New("빈 문자열 입니다")
	}
	p := strings.Replace(path, "\\", "/", -1)
	regRule := `/show[/_]\S+?/seq/\S+?/\S+?_(\S[^/]+)`
	if strings.HasPrefix(p, "/backup/") {
		regRule = `/backup/\d+?/\S+?/\S+?/seq/\S+?/\S+?_(\S[^/]+)`
	}
	re, err := regexp.Compile(regRule)
	if err != nil {
		return "", err
	}
	results := re.FindStringSubmatch(p)
	if results == nil {
		return "", errors.New(path + " 경로에서 샷 정보를 가지고 올 수 없습니다.")
	}
	return results[len(results)-1], nil
}

// Task 함수는 경로를 받아서 Task를 반환한다.
func Task(path string) (string, error) {
	if path == "" {
		return path, errors.New("빈 문자열 입니다")
	}
	p := strings.Replace(path, "\\", "/", -1)
	regRule := `/show[/_]\S+?/seq/\S+?/\S+?_\S+?/(\S[^/]+)`
	if strings.HasPrefix(p, "/backup/") {
		regRule = `/backup/\d+?/\S+?/\S+?/seq/\S+?/\S+?_\S+?/(\S[^/]+)`
	}
	re, err := regexp.Compile(regRule)
	if err != nil {
		return "", err
	}
	results := re.FindStringSubmatch(p)
	if results == nil {
		return "", errors.New(path + " 경로에서 Task 정보를 가지고 올 수 없습니다.")
	}
	return results[len(results)-1], nil
}

// Element 함수는 경로를 받아서 Element를 반환한다.
// 회사는 아직 Element를 도입중이다. 이 함수는 아직 느슨하게 체크한다.
func Element(path string) (string, error) {
	if path == "" {
		return path, errors.New("빈 문자열 입니다")
	}
	p := strings.Replace(path, "\\", "/", -1)
	regRule := `/show[/_]\S+?/seq/\S+?/\S+?_\S+?/\S[^/]+/\S+?/([a-z0-9]+[^/])`
	if strings.HasPrefix(p, "/backup/") {
		regRule = `/backup/\d+?/\S+?/\S+?/seq/\S+?/\S+?_\S+?/\S[^/]+/\S+?/([a-z0-9]+[^/])`
	}
	re, err := regexp.Compile(regRule)
	if err != nil {
		return "", err
	}
	results := re.FindStringSubmatch(p)
	if results == nil {
		return "", errors.New(path + " 경로에서 Element 정보를 가지고 올 수 없습니다.")
	}
	return results[len(results)-1], nil
}

// Exist 함수는 경로가 존재하는지 체크한다.
// Go에서는 파이썬처럼 간단하게 작성하기 힘들기 때문에 편의를 위해 추가한다.
func Exist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// Seqnum2Sharp 함수는 경로와 파일명을 받아서 시퀀스부분을 #문자열로 바꾸고 시퀀스의 숫자를 int로 바꾼다.
// "test.0002.jpg" -> "test.####.jpg", 2, nil
func Seqnum2Sharp(filename string) (string, int, error) {
	re, err := regexp.Compile("([0-9]+)(\\.[a-zA-Z]+$)")
	// 이 정보를 통해서 파일명을 구하는 방식으로 바꾼다.
	if err != nil {
		return filename, -1, errors.New("정규 표현식이 잘못되었습니다")
	}
	results := re.FindStringSubmatch(filename)
	if results == nil {
		return filename, -1, errors.New("경로가 시퀀스 형식이 아닙니다")
	}
	seq := results[1]
	ext := results[2]
	header := filename[:strings.LastIndex(filename, seq+ext)]
	seqNum, err := strconv.Atoi(seq)
	if err != nil {
		return filename, -1, err
	}
	return header + strings.Repeat("#", len(seq)) + ext, seqNum, nil
}

// Sharp2Seqnum 함수는 경로의 #문자를 숫자(n)로 치환하는 함수이다.
func Sharp2Seqnum(path string, n int) (string, error) {
	sharpNum := strings.Count(path, "#")
	if sharpNum == 0 {
		return path, nil
	}
	strNum := strconv.Itoa(n)
	if sharpNum < len(strNum) {
		return "", fmt.Errorf("%s에 %d 숫자를 담을 수 없습니다", strings.Repeat("#", sharpNum), n)
	}
	listfile := strings.Split(path, strings.Repeat("#", sharpNum))
	head := listfile[0]               // #문자의 앞쪽
	tail := listfile[len(listfile)-1] // #문자의 뒤쪽
	return head + fmt.Sprintf("%0"+strconv.Itoa(sharpNum)+"d", n) + tail, nil
}
