package logging

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

func Print(a ...interface{}) {
	fmt.Print(a...)
}

func Println(a ...interface{}) {
	fmt.Println(a...)
}

func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func Important(msg string) {
	print(color.Bold, msg)
}

func Infof(format string, args ...interface{}) {
	printf(color.FgCyan, format, args)
}

func Info(msg string) {
	println(color.FgCyan, msg)
}

func Successf(format string, args ...interface{}) {
	printf(color.FgGreen, format, args...)
}

func Success(msg string) {
	println(color.FgGreen, msg)
}

func Warning(msg string) {
	print(color.FgYellow, msg)
}

func Warningln(msg string) {
	println(color.FgYellow, msg)
}

func Warningf(format string, args ...interface{}) {
	printf(color.FgYellow, format, args...)
}

func Failure(msg string) {
	println(color.FgRed, msg)
}

func Failuref(format string, args ...interface{}) {
	printf(color.FgRed, format, args...)
}

func Error(a ...interface{}) {
	logrus.Error(a...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatal(a ...interface{}) {
	logrus.Fatal(a...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func printf(colorCode color.Attribute, format string, args ...interface{}) {
	_, err := color.New(colorCode).Printf(format, args...)
	if err != nil {
		logrus.Errorf("failed to print: %v", err)
	}
}

func println(colorCode color.Attribute, msg string) {
	_, err := color.New(colorCode).Println(msg)
	if err != nil {
		logrus.Errorf("failed to print: %v", err)
	}
}

func print(colorCode color.Attribute, msg string) {
	_, err := color.New(colorCode).Print(msg)
	if err != nil {
		logrus.Errorf("failed to print: %v", err)
	}
}
