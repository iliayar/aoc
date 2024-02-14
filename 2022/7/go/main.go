package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Entry struct {
	isDir   bool
	parent  *Entry
	entries map[string]*Entry
	size    int
}

func CreateDir(parent *Entry) *Entry {
	return &Entry{
		isDir:   true,
		parent:  parent,
		entries: make(map[string]*Entry),
	}
}

func CreateFile(size int) *Entry {
	return &Entry{
		isDir: false,
		size:  size,
	}
}

type State struct {
	pwd  *Entry
	root *Entry
}

func (state *State) up() {
	if state.pwd.parent == nil {
		panic("Trying get parent of root")
	}

	state.pwd = state.pwd.parent
}

func (state *State) down(dir string) {
	newDir, ok := state.pwd.entries[dir]

	if !ok {
		newDir = CreateDir(state.pwd)
		state.pwd.entries[dir] = newDir
	}

	if !newDir.isDir {
		panic(fmt.Sprintf("%s is not a directory", dir))
	}

	state.pwd = newDir
}

func (state *State) cd(path string) {
	if path == ".." {
		state.up()
	} else if path == "/" {
		state.pwd = state.root
	} else {
		state.down(path)
	}
}

func (state *State) addFile(filename string, size int) {
	existingEntry, ok := state.pwd.entries[filename]
	if ok {
		if existingEntry.isDir {
			panic(fmt.Sprintf("%s directory already exists", filename))
		}
		if existingEntry.size != size {
			panic(fmt.Sprintf("%s fiel already exists and %d != %d", filename, existingEntry.size, size))
		}
		return
	}

	state.pwd.entries[filename] = CreateFile(size)
}

func (state *State) addDirectory(dirname string) {
	existingEntry, ok := state.pwd.entries[dirname]
	if ok {
		if !existingEntry.isDir {
			panic(fmt.Sprintf("%s file already exists", dirname))
		}
		return
	}

	state.pwd.entries[dirname] = CreateDir(state.pwd)
}

type Reader struct {
	state *State
	lines []string
}

func (reader *Reader) peakLine() string {
	return reader.lines[0]
}

func (reader *Reader) nextLine() string {
	curLine := reader.lines[0]
	reader.lines = reader.lines[1:]
	return curLine
}

func (reader *Reader) nextCmd() {
	curLine := reader.nextLine()

	cmd, isCmd := strings.CutPrefix(curLine, "$ ")

	if !isCmd {
		panic(fmt.Sprintf("Expected command line, got %s", curLine))
	}

	args := strings.Split(cmd, " ")

	cmd, args = args[0], args[1:]

	reader.handleCmd(cmd, args)
}

func (reader *Reader) handleCmd(cmd string, args []string) {
	switch cmd {
	case "cd":
		reader.state.cd(args[0])
	case "ls":
		for !reader.isEof() && !strings.HasPrefix(reader.peakLine(), "$ ") {
			fileLine := reader.nextLine()
			fields := strings.Split(fileLine, " ")

			t, name := fields[0], fields[1]

			switch t {
			case "dir":
				reader.state.addDirectory(name)
			default:
				size, err := strconv.Atoi(t)
				if err != nil {
					panic(fmt.Sprintf("Failed to parse number %s", t))
				}

				reader.state.addFile(name, size)
			}
		}
	}
}

func (reader *Reader) isEof() bool {
	return len(reader.lines) == 0
}

func countSize(entry *Entry, cb func(*Entry, int)) int {
	if !entry.isDir {
        cb(entry, entry.size)
        return entry.size
	}

    size := 0

    for _, child := range entry.entries {
        childSize := countSize(child, cb)
        size += childSize
    }

    cb(entry, size)
    return size
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	contentStr := string(content)

	root := CreateDir(nil)
	state := State{
		pwd:  root,
		root: root,
	}

	lines := strings.Split(strings.Trim(contentStr, "\r\n "), "\n")
	reader := Reader{
		state: &state,
		lines: lines,
	}

	for !reader.isEof() {
		reader.nextCmd()
	}

    const TOTAL = 70000000
    const NEED = 30000000

    used := countSize(root, func(entry *Entry, size int) {})
    toDelete := NEED - (TOTAL - used)

    res := math.MaxInt
    countSize(root, func(entry *Entry, size int) {
        if entry.isDir && size >= toDelete {
            if res > size {
                res = size
            }
        }
    })

    fmt.Println(res)
}
