package judgetools

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

//调用os.MkdirAll递归创建文件夹
func CreateFile(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		os.Chmod(filePath, 0777)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func FastCreateFile(path string) {
	createErr := CreateFile(path)
	if createErr != nil {
		fmt.Println("create file error: ", createErr)
		panic(createErr)
	}
}

func DelFile(path string) error {
	cmd := exec.Command("rm", "-rf", path)
	err := cmd.Run()
	if err != nil {
		fmt.Println("del file: " + path + " error: ", err)
		return err
	}
	return nil
}

func FastDelFile(path string) {
	delErr := DelFile(path)
	if delErr != nil {
		fmt.Println("create file error: ", delErr)
		panic(delErr)
	}
}




func Compare(spath, dpath string) bool {
	sinfo, err := os.Lstat(spath)
	if err != nil {
		return false
	}
	dinfo, err := os.Lstat(dpath)
	if err != nil {
		return false
	}
	if sinfo.Size() != dinfo.Size() {
		return false
	}
	return comparefile(spath, dpath)
}

func comparefile(spath, dpath string) bool {
	sFile, err := os.Open(spath)
	if err != nil {
		return false
	}
	dFile, err := os.Open(dpath)
	if err != nil {
		return false
	}
	fmt.Println("compare eeeee : ====")
	fmt.Println(spath)
	fmt.Println(dpath)
	fmt.Println(sFile)
	fmt.Println(dFile)
	b := comparebyte(sFile, dFile)
	sFile.Close()
	dFile.Close()
	return b
}

//下面可以代替md5比较.
func comparebyte(sfile *os.File, dfile *os.File) bool {
	var sbyte []byte = make([]byte, 512)
	var dbyte []byte = make([]byte, 512)
	var serr, derr error
	for {
		_, serr = sfile.Read(sbyte)
		_, derr = dfile.Read(dbyte)
		if serr != nil || derr != nil {
			if serr != derr {
				return false
			}
			if serr == io.EOF {
				break
			}
		}
		if bytes.Equal(sbyte, dbyte) {
			continue
		}
		return false
	}
	return true
}
