package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

//文件日志输出器

type FileOutputerOptions struct {
	filename      string //文件名
	lastSplitHour int    //上次分割的时间
}

type FileOutputer struct {
	file       *os.File //系统日志
	accessFile *os.File //访问日志
	option     *FileOutputerOptions
}

func NewFileOutputer(filename string) (Outputer, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("filename: %v\n", filename)
	option := &FileOutputerOptions{
		filename: filename,
	}
	log := &FileOutputer{
		option: option,
	}

	err = log.init()
	return log, err
}

//生成当前实际文件名和正常文件名
//实际文件名: hello.log.2022121823
//正常文件名: hello.log
func (f *FileOutputer) getCurFilename() (curFilename, originFilename string) {
	now := time.Now()
	curFilename = fmt.Sprintf("%s.%04d%02d%02d%02d", f.option.filename, now.Year(), now.Month(), now.Day(), now.Hour())
	originFilename = f.option.filename
	return
}

//生成当前实际访问日志文件名和正常访问日志文件名
//实际文件名: hello.log.access.2022121823
//正常文件名: hello.log.access
func (f *FileOutputer) getCurAccessFilename() (curAccessFilename, originAccessFilename string) {
	now := time.Now()
	curAccessFilename = fmt.Sprintf("%s.access.%04d%02d%02d%02d", f.option.filename, now.Year(), now.Month(), now.Day(), now.Hour())
	originAccessFilename = fmt.Sprintf("%s.access", f.option.filename)
	return
}

//打开文件, 指向当前小时的文件
//如果传入"hello.log", 那就直接指向"hello.log.2022121823"这个文件
func (f *FileOutputer) initFile(filename, originFilename string) (file *os.File, err error) {

	//获取文件目录
	// fmt.Println(filepath.Dir(filename))
	dir := filepath.Dir(filename)
	fileInfo, _ := os.Stat(filepath.Dir(filename))
	if fileInfo == nil {
		os.MkdirAll(dir, 0755)
	}

	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		err = fmt.Errorf("open file %s failed, err : %v", filename, err)
		return
	}

	//创建软连接. 用旧文件名指向新文件名
	os.Symlink(filename, originFilename)
	return
}

//初始化文件日志
func (f *FileOutputer) init() (err error) {
	//1. 获取当前文件名和正常文件名
	curFilename, originFilename := f.getCurFilename()
	f.file, err = f.initFile(curFilename, originFilename)
	if err != nil {
		return
	}

	accessFilename, originAccessFilename := f.getCurAccessFilename()
	f.accessFile, err = f.initFile(accessFilename, originAccessFilename)
	if err != nil {
		return
	}

	f.option.lastSplitHour = time.Now().Hour()
	return
}

//检查日志是否需要切分
func (f *FileOutputer) checkSplitFile(curTime time.Time) {
	hour := curTime.Hour()
	//如果当前小时数 = 上次切分的小时数, 无需切分
	if hour == f.option.lastSplitHour {
		return
	}

	//如果超过了上次切分的小时数, 那么调用init初始化
	f.init()
}

//实现接口写日志
func (f *FileOutputer) Write(data *LogData) {
	f.checkSplitFile(data.curTime)
	file := f.file
	if data.level == LogLevelAccess {
		file = f.accessFile
	}

	file.Write(data.Bytes())
}

func (f *FileOutputer) Close() {
	f.file.Close()
	f.accessFile.Close()
}
