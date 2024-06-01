package db

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ronchi-oss/bib/conf"
	"gopkg.in/yaml.v3"
)

type Note struct {
	ID        int       `yaml:"id,omitempty"`
	CreatedAt time.Time `yaml:"created_at"`
	Pinned    bool      `yaml:"pinned"`
	Title     string    `yaml:"title,omitempty"`
	Body      []byte    `yaml:"body,omitempty"`
}

type NaturalSort []os.DirEntry

func (d NaturalSort) Len() int      { return len(d) }
func (d NaturalSort) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d NaturalSort) Less(i, j int) bool {
	id1, err := strconv.Atoi(d[i].Name())
	if err != nil {
		panic(err)
	}
	id2, err := strconv.Atoi(d[j].Name())
	if err != nil {
		panic(err)
	}
	return id1 < id2
}

func GetNotes(targetDir string) ([]*Note, error) {
	files, err := os.ReadDir(fmt.Sprintf("%s/src", targetDir))
	if err != nil {
		return nil, err
	}
	sort.Sort(NaturalSort(files))
	var notes []*Note
	for _, file := range files {
		id, err := strconv.Atoi(file.Name())
		if err != nil {
			return nil, err
		}
		note, err := GetNote(targetDir, id)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func AppendNote(targetDir string, title string) (int, error) {
	if err := os.MkdirAll(fmt.Sprintf("%s/src", targetDir), 0750); err != nil {
		return 0, err
	}
	var newID int
	files, err := os.ReadDir(fmt.Sprintf("%s/src", targetDir))
	if err != nil {
		return 0, err
	}
	sort.Sort(NaturalSort(files))
	if len(files) == 0 {
		newID = 1
	} else {
		largestID, err := strconv.Atoi(files[len(files)-1].Name())
		if err != nil {
			return 0, err
		}
		newID = largestID + 1
	}
	d, err := conf.YAMLEncode(&Note{
		ID:        newID,
		CreatedAt: time.Now(),
		Pinned:    false,
	})
	if err != nil {
		return 0, err
	}
	if err := os.Mkdir(fmt.Sprintf("%s/src/%d", targetDir, newID), 0750); err != nil {
		return 0, err
	}
	f1, err := os.OpenFile(fmt.Sprintf("%s/src/%d/metadata.yml", targetDir, newID), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer f1.Close()
	if _, err := f1.Write(d); err != nil {
		return 0, err
	}
	f2, err := os.OpenFile(fmt.Sprintf("%s/src/%d/README.md", targetDir, newID), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer f2.Close()
	if _, err := f2.WriteString(fmt.Sprintf("# %s\n", title)); err != nil {
		return 0, err
	}
	return newID, nil
}

func WriteNoteMetadata(targetDir string, note *Note) error {
	d, err := conf.YAMLEncode(note)
	if err != nil {
		return fmt.Errorf("cannot marshal note: %v", err)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s/src/%d/metadata.yml", targetDir, note.ID), d, 0644); err != nil {
		return fmt.Errorf("cannot write metadata file: %v", err)
	}
	return nil
}

func GetNote(targetDir string, id int) (*Note, error) {
	yamlFile, err := os.ReadFile(fmt.Sprintf("%s/src/%d/metadata.yml", targetDir, id))
	if err != nil {
		return nil, fmt.Errorf("cannot read note database file (id = %d): %v", id, err)
	}
	var note Note
	if err := yaml.Unmarshal(yamlFile, &note); err != nil {
		return nil, fmt.Errorf("cannot unmarshal data: %v", err)
	}
	mdFile, err := os.Open(fmt.Sprintf("%s/src/%d/README.md", targetDir, id))
	if err != nil {
		return nil, fmt.Errorf("cannot open Markdown file (id = %d): %v", id, err)
	}
	defer mdFile.Close()
	var (
		scanner     = bufio.NewScanner(mdFile)
		mdFirstLine string
	)
	for scanner.Scan() {
		mdFirstLine = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cannot scan Markdown file (id = %d): %v", id, err)
	}
	note.ID = id
	note.Title = strings.TrimLeft(mdFirstLine, "# ")
	return &note, nil
}

func TogglePin(targetDir string, id int) error {
	note, err := GetNote(targetDir, id)
	if err != nil {
		return err
	}
	note.Pinned = !note.Pinned
	return WriteNoteMetadata(targetDir, note)
}
