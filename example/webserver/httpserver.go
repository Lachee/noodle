// +build !js

package main

import (
	"flag"
	"io"
	"log"
	"net/http"
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

var (
	cmd               string
	filter            string
	resourceDirectory string
	globFilter        glob.Glob
	dirs              arrayFlags
	args              arrayFlags
	current           *exec.Cmd
	lastTime          time.Time
	watcher           *fsnotify.Watcher
	client            *wsclient
)

func watchFiles() {
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
                            client.Broadcast(payloadEvent{ Event: "BuildSuccess", Asset: absFile })
                            log.Println("Build Complete")
						} else {
                            client.Broadcast(payloadEvent{ Event: "BuildFailure", Asset: absFile })
							log.Println("Build Failed", err)
						}

						//Clear the current out
						current = nil
                    } else {                        
                        client.Broadcast(payloadEvent{ Event: "AssetUpdated", Asset: absFile })
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
}

func main() {
	//Parse the flag
	flag.StringVar(&cmd, "cmd", "go build", "The command that will be executed when a change has been discovered")
	flag.StringVar(&filter, "filter", "*.go", "Filters the files that are modified")
	flag.StringVar(&resourceDirectory, "resources", "./resources/", "Resource Directory")
	flag.Var(&args, "args", "Arguments for the command")
	flag.Var(&dirs, "dir", "Folders to listen for changes.")
	flag.Parse()

	//Setup the file watcher
	globFilter = glob.MustCompile(filter)

	fileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer fileWatcher.Close()
	watcher = fileWatcher

	for _, f := range dirs {
		log.Println("Watching: ", f)
		err = watcher.Add(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	go watchFiles()

	//SErve the files
	baseFileServe := http.FileServer(http.Dir("./"))
	http.Handle("/", http.StripPrefix("/", baseFileServe))

	resourceFileServe := http.FileServer(http.Dir(resourceDirectory))
	http.Handle("/resources/", http.StripPrefix("/resources/", resourceFileServe))

	//Listens
	client = &wsclient{}
	http.HandleFunc("/listen", client.handle)

	http.ListenAndServe(":8090", nil)
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
