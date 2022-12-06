package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/common-nighthawk/go-figure"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/html/charset"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

var globalpath string
var outputfileneeed bool = false
var outputpath string
var showme bool
var sliceofslices [][]string

// Task represents the <Task> element in the XML file
type Task struct {
	XMLName          xml.Name         `xml:"Task"`
	Version          string           `xml:"version,attr"`
	RegistrationInfo RegistrationInfo `xml:"RegistrationInfo"`
	Triggers         Triggers         `xml:"Triggers"`
	Principals       Principals       `xml:"Principals"`
	Settings         Settings         `xml:"Settings"`
	Actions          Actions          `xml:"Actions"`
}

// RegistrationInfo represents the <RegistrationInfo> element in the XML file
type RegistrationInfo struct {
	Version     string `xml:"Version"`
	Description string `xml:"Description"`
	URI         string `xml:"URI"`
}

// Triggers represents the <Triggers> element in the XML file
type Triggers struct {
	LogonTrigger    LogonTrigger    `xml:"LogonTrigger"`
	CalendarTrigger CalendarTrigger `xml:"CalendarTrigger"`
}

// LogonTrigger represents the <LogonTrigger> element in the XML file
type LogonTrigger struct {
	Enabled bool `xml:"Enabled"`
}

// CalendarTrigger represents the <CalendarTrigger> element in the XML file
type CalendarTrigger struct {
	StartBoundary string        `xml:"StartBoundary"`
	ScheduleByDay ScheduleByDay `xml:"ScheduleByDay"`
}

// ScheduleByDay represents the <ScheduleByDay> element in the XML file
type ScheduleByDay struct {
	DaysInterval int `xml:"DaysInterval"`
}

type Principals struct {
	Principal Principal `xml:"Principal"`
}

// Principal represents the <Principal> element in the XML file
type Principal struct {
	ID       string `xml:"id,attr"`
	UserID   string `xml:"UserId"`
	RunLevel string `xml:"RunLevel"`
}

// Settings represents the <Settings> element in the XML file
type Settings struct {
	MultipleInstancesPolicy    string `xml:"MultipleInstancesPolicy"`
	DisallowStartIfOnBatteries bool   `xml:"DisallowStartIfOnBatteries"`
	StartWhenAvailable         bool   `xml:"StartWhenAvailable"`
	RunOnlyIfNetworkAvailable  bool   `xml:"RunOnlyIfNetworkAvailable"`
	Enabled                    bool   `xml:"Enabled"`
	RunOnlyIfIdle              bool   `xml:"RunOnlyIfIdle"`
	WakeToRun                  bool   `xml:"WakeToRun"`
	ExecutionTimeLimit         string `xml:"ExecutionTimeLimit"`
}

type Actions struct {
	Context string `xml:"Context,attr"`
	Exec    Exec   `xml:"Exec"`
}

type Exec struct {
	Command   string `xml:"Command"`
	Arguments string `xml:"Arguments"`
}

func getalltaskfiles(path string) {

	// Create a slice to hold the file paths
	var files []string

	// Recursively search for files
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path is a regular file
		if info.Mode().IsRegular() {
			files = append(files, path)
			parseTask(path)
			if err != nil {
				ct.Foreground(ct.Yellow, true)
				fmt.Println("[-] Failed to parse task:", err)
				ct.ResetColor()
			}
		}

		return nil
	})
	if err != nil {
		ct.Foreground(ct.Red, true)
		fmt.Println("[!] Failed to search for files:", err)
		ct.ResetColor()
		return
	}

}

func parseTask(path string) {
	f16, err := os.Open(path)
	if err != nil {
		return
	}

	task := Task{}
	err = DecodeUtf16XML(f16, &task)
	if err != nil {
		return
	} else {

		newslice := []string{task.Version, task.RegistrationInfo.Description, task.RegistrationInfo.URI, task.Triggers.CalendarTrigger.StartBoundary, strconv.Itoa(task.Triggers.CalendarTrigger.ScheduleByDay.DaysInterval), task.Principals.Principal.UserID, task.Principals.Principal.RunLevel, task.Settings.MultipleInstancesPolicy, strconv.FormatBool(task.Settings.DisallowStartIfOnBatteries), strconv.FormatBool(task.Settings.StartWhenAvailable), strconv.FormatBool(task.Settings.RunOnlyIfNetworkAvailable), strconv.FormatBool(task.Settings.Enabled), strconv.FormatBool(task.Settings.RunOnlyIfIdle), strconv.FormatBool(task.Settings.WakeToRun), task.Settings.ExecutionTimeLimit, task.Actions.Exec.Command, task.Actions.Exec.Arguments}
		sliceofslices = append(sliceofslices, newslice)
		//table := tablewriter.NewWriter(os.Stdout)
		//table.SetHeader([]string{"Version", "Description", "URI", "StartBoundary", "DaysInterval", "UserId", "RunLevel", "MultipleInstancesPolicy", "DisallowStartIfOnBatteries", "StartWhenAvailable", "RunOnlyIfNetworkAvailable", "Enabled", "RunOnlyIfIdle", "WakeToRun", "ExecutionTimeLimit", "Command", "Arguments"})
		//table.Append([]string{task.Version, task.RegistrationInfo.Description, task.RegistrationInfo.URI, task.Triggers.CalendarTrigger.StartBoundary, strconv.Itoa(task.Triggers.CalendarTrigger.ScheduleByDay.DaysInterval), task.Principals.Principal.UserID, task.Principals.Principal.RunLevel, task.Settings.MultipleInstancesPolicy, strconv.FormatBool(task.Settings.DisallowStartIfOnBatteries), strconv.FormatBool(task.Settings.StartWhenAvailable), strconv.FormatBool(task.Settings.RunOnlyIfNetworkAvailable), strconv.FormatBool(task.Settings.Enabled), strconv.FormatBool(task.Settings.RunOnlyIfIdle), strconv.FormatBool(task.Settings.WakeToRun), task.Settings.ExecutionTimeLimit, task.Actions.Exec.Command, task.Actions.Exec.Arguments})
		//table.Render()
	}
}

func findTasksDir(basepath string) string {
	// Start the search at the root of the file system

	var tasksDir string

	// Recursively search for directories named "Tasks"
	err := filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path is a directory and is named "Tasks"
		if info.IsDir() && info.Name() == "Tasks" {
			tasksDir = path
			return nil
		}

		return nil
	})
	if err != nil {
		ct.Foreground(ct.Red, true)
		fmt.Println("[!] Failed to search for Tasks directories:", err)
		ct.ResetColor()
		return ""
	}

	return tasksDir
}

func main() {
	ct.Foreground(ct.Green, true)
	totesimportant := figure.NewFigure("Greg", "doom", true)

	totesimportant.Print()
	fmt.Println("\nThe Taskmaster - Windows Task XML parser\n")
	parser := argparse.NewParser("Okonma", "Go Wrapper for Forensic Tools")
	i := parser.String("i", "inputdir", &argparse.Options{Required: true, Help: "Input Directory"})
	o := parser.String("o", "outputdir", &argparse.Options{Required: false, Help: "(optional) param. Output Directory"})
	var showOnScreen *bool = parser.Flag("s", "show", &argparse.Options{Required: false, Help: "Dont Output to a file - show me the results onscreen"})
	ct.ResetColor()
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))

	}

	var inputset bool = false

	if *i != "" {
		globalpath = *i
		inputset = true
	}

	if *i == "" {
		os.Exit(1)
	}

	if *o != "" {
		outputpath = *o
		outputfileneeed = true
	}

	if *showOnScreen {
		showme = true
	}

	if inputset && showme == false && showme == false {
		ct.Background(ct.Yellow, true)
		fmt.Println("[-] You need to specific some action to take")
		ct.ResetColor()
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if inputset && (showme == true || outputfileneeed == true) {
		getalltaskfiles(*i)
		if showme == true {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Version", "Description", "URI", "StartBoundary", "DaysInterval", "UserId", "RunLevel", "MultipleInstancesPolicy", "DisallowStartIfOnBatteries", "StartWhenAvailable", "RunOnlyIfNetworkAvailable", "Enabled", "RunOnlyIfIdle", "WakeToRun", "ExecutionTimeLimit", "Command", "Arguments"})
			for _, value := range sliceofslices {
				table.Append(value)
			}
			table.Render()
		}
		if outputfileneeed == true {
			var newslice [][]string
			var tempslice []string
			tempslice = []string{"Version", "Description", "URI", "StartBoundary", "DaysInterval", "UserId", "RunLevel", "MultipleInstancesPolicy", "DisallowStartIfOnBatteries", "StartWhenAvailable", "RunOnlyIfNetworkAvailable", "Enabled", "RunOnlyIfIdle", "WakeToRun", "ExecutionTimeLimit", "Command", "Arguments"}
			newslice = append(newslice, tempslice)
			for _, value := range sliceofslices {
				newslice = append(newslice, value)
			}

			f, err := os.Create(outputpath)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			// Create a new CSV writer.
			w := csv.NewWriter(f)

			// Iterate over the slice of slices and write each row to the CSV writer.
			for _, v := range newslice {
				if err := w.Write(v); err != nil {
					panic(err)
				}

			}
			ct.Foreground(ct.Green, true)
			fmt.Println("[+] Output written to " + outputpath)
			ct.ResetColor()
			// Flush the writer to ensure that all data is written to the file.
			w.Flush()
		}

	}
}

func DecodeUtf16XML(r io.Reader, v interface{}) (err error) {
	//kudos to the below for really helping me out
	// https://www.tipitaka.org/romn/cscd/vin01m.mul.toc.xml
	// The Tipiṭaka XML is encoded in UTF-16
	// Google search: golang xml utf-16
	// https://stackoverflow.com/questions/6002619/unmarshal-an-iso-8859-1-xml-input-in-go
	// https://groups.google.com/forum/#!topic/golang-nuts/tXcECEKC2rs
	nr, err := charset.NewReader(r, "utf-16")
	if err != nil {
		return
	}
	decoder := xml.NewDecoder(nr)
	decoder.CharsetReader = BypassReader
	err = decoder.Decode(v)
	return
}
func BypassReader(label string, input io.Reader) (io.Reader, error) {
	return input, nil
}
