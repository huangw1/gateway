/**
 * @Author: huangw1
 * @Date: 2019/7/10 20:11
 */

package logging

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}
