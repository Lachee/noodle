package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gobwas/glob"
)

//arrayFlags is an array of string flags.
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "ArrayFlags"
}
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

//Prepare the variables
var (
	cmd        string
	filter     string
	globFilter glob.Glob
	dirs       arrayFlags
	args       arrayFlags
	current    *exec.Cmd
	lastTime   time.Time
)

func main() {

	//Parse the flag
	flag.StringVar(&cmd, "cmd", "go build", "The command that will be executed when a change has been discovered")
	flag.StringVar(&filter, "filter", "*.go", "Filters the files that are modified")
	flag.Var(&args, "args", "Arguments for the command")
	flag.Var(&dirs, "dir", "Folders to listen for changes.")
	flag.Parse()

	globFilter = glob.MustCompile(filter)

	//Setup the watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	//Create the watcher callback
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				now := time.Now()
				diff := now.Sub(lastTime)
				if diff.Milliseconds() >= 1000 {
					if event.Op&fsnotify.Write == fsnotify.Write {
						absFile, _ := filepath.Abs(event.Name)
						absFile = strings.Replace(absFile, "\\", "/", -1)

						if globFilter.Match(absFile) {
							log.Println("Starting Build")
							cmd, err := runCommand(cmd, args...)
							if err == nil {
								lastTime = now
								current = cmd
								current.Wait()
								log.Println("Build Complete")
							} else {
								log.Println("Build Failed", err)
							}

							//Clear the current out
							current = nil
						}
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("error:", err)
			}
		}
	}()

	//Watch the folder and wait for exit
	for _, f := range dirs {
		log.Println("Watching: ", f)
		err = watcher.Add(f)
		if err != nil {
			log.Fatal(err)
		}
	}
	<-done
}

// runCommand runs the command with given name and arguments. It copies the
// logs to standard output
func runCommand(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return cmd, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cmd, err
	}

	if err := cmd.Start(); err != nil {
		return cmd, err
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	return cmd, nil
}
