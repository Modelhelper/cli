package project

import (
	"io/ioutil"
	"log"
	"modelhelper/cli/code"
	"modelhelper/cli/source"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const dirname string = ".modelhelper"

type Project struct {
	Version       string                       `yaml:"version"`
	Name          string                       `yaml:"name"`
	Language      string                       `yaml:"language"`
	Description   string                       `yaml:"description"`
	DefaultSource string                       `yaml:"defaultSource,omitempty"`
	DefaultKey    string                       `yaml:"defaultKey,omitempty"`
	Connections   map[string]source.Connection `yaml:"connections,omitempty"`
	Code          map[string]code.Code         `yaml:"code,omitempty"`
	OwnerName     string                       `yaml:"ownerName,omitempty"`
	Options       map[string]string            `yaml:"options,omitempty"`
	Custom        interface{}                  `yaml:"custom,omitempty"`
	Header        string                       `yaml:"header,omitempty"`
}

func (p *Project) Save() error {

	d, err := yaml.Marshal(&p)

	if err != nil {

		return err
	}

	// path := filepath.Join(DefaultLocation(), "config.yaml")
	path := DefaultLocation()

	if !Exists(path) {
		CreateDir(dirname)
	}

	err = ioutil.WriteFile(path, d, 0777)

	return err
}

func CreateDir(name string) {
	err := os.Mkdir(name, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
func DefaultDir() string {
	p, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return filepath.Join(p, dirname)
}
func DefaultLocation() string {
	// p, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// }
	return filepath.Join(DefaultDir(), "project.yaml")
}

func (P *Project) Exists(path string) bool {

	pathInfo, err := os.Stat(path)

	if os.IsNotExist(err) || pathInfo.IsDir() {
		return false
	}

	return true
}

func Exists(path string) bool {

	pathInfo, err := os.Stat(path)

	if os.IsNotExist(err) || pathInfo.IsDir() {
		return false
	}

	return true
}

func LoadProjects(path ...string) []Project {
	l := []Project{}

	for _, p := range path {
		project, _ := loadProjectFromFile(p)
		l = append(l, *project)
	}
	return l
}
func Load(path string) (*Project, error) {
	if len(path) > 0 {
		pathInfo, err := os.Stat(path)
		if os.IsNotExist(err) || pathInfo.IsDir() {
			// log.Fatal("Project does not exits")
			return nil, err
		}

	} else {
		p, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		path = filepath.Join(p, dirname, "project.yaml")
	}

	f, err := loadProjectFromFile(path)

	return f, err

}

func loadProjectFromFile(fileName string) (*Project, error) {
	var p *Project

	dat, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
		return nil, e
	}

	err := yaml.Unmarshal(dat, &p)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, err
	}

	return p, nil
}

func (p *Project) GetConnections() (*map[string]source.Connection, error) {
	return &p.Connections, nil
}

//FindRelatedProjects gets a list of all the related projects by following the path from DefaultDir() to the volumeroot
//The returning list is in correct importance from least to most (the project nearest to DefaultDir())
func FindReleatedProjects(startPath string) []string {
	basePath, _ := filepath.Split(startPath)
	list := []string{}

	dirs := strings.Split(startPath, string(os.PathSeparator))
	if len(dirs) > 0 {
		dirs[0] = filepath.VolumeName(startPath) + string(os.PathSeparator)
	}
	for i := 0; i <= len(dirs); i++ {
		basePath = filepath.Join(dirs[0:i]...)
		if len(basePath) == 0 {
			continue
		}

		if dirs[i] == dirname {
			fp := filepath.Join(basePath, dirname, "project.yaml")
			list = append(list, fp)
			break
		}

		files, err := os.ReadDir(basePath)
		// files, err := ioutil.ReadDir(basePath)
		if err != nil {
			log.Fatal(err)
			break
		}
		for _, f := range files {

			if f.IsDir() && f.Name() == dirname {
				fp := filepath.Join(basePath, f.Name(), "project.yaml")
				list = append(list, fp)
			}

		}

	}
	return list

}

func FindNearestProjectDir() (string, bool) {
	basePath := "./"
	root := DefaultDir()
	updir := "../"

	for i := 1; i < 6; i++ {
		files, err := ioutil.ReadDir(basePath)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {

			if f.IsDir() && f.Name() == dirname {
				fp := filepath.Join(basePath, f.Name(), "project.yaml")
				return fp, true
			}
		}

		relp, _ := filepath.Rel(root, basePath)
		basePath = filepath.Join(updir, relp)

		updir += "../"
	}

	return "", false
}

func JoinProject(joinType string, projects ...Project) Project {
	switch joinType {
	case "merge":
		return mergeProject(projects...)

	case "smart":
		return smartMergeProject(projects...)

	case "replace":
		return replaceProject(projects...)
	default:
		return smartMergeProject(projects...)

	}
}
func mergeProject(projects ...Project) Project {
	current := Project{}

	conProv := []source.ConnectionProvider{}
	for _, p := range projects {
		conProv = append(conProv, &p)
	}

	current.Connections = source.JoinConnections("merge", conProv...)

	for _, proj := range projects {

		current.Code = proj.Code
	}

	return current
}

func smartMergeProject(projects ...Project) Project {
	current := Project{}
	current.Options = make(map[string]string)
	conProv := []source.ConnectionProvider{}
	for _, p := range projects {
		conProv = append(conProv, &p)
	}

	current.Connections = source.JoinConnections("smart", conProv...)

	for _, proj := range projects {
		current.Name = mergeString(current.Name, proj.Name)
		// if len(proj.Name) > 0 {
		// 	current.Name = proj.Name
		// }
		current.Description = mergeString(current.Description, proj.Description)
		current.DefaultKey = mergeString(current.DefaultKey, proj.DefaultKey)
		current.DefaultSource = mergeString(current.DefaultSource, proj.DefaultSource)
		current.OwnerName = mergeString(current.OwnerName, proj.OwnerName)
		current.Language = mergeString(current.Language, proj.Language)

		for optKey, optVal := range proj.Options {
			current.Options[optKey] = optVal
		}

		current.Code = proj.Code
	}

	return current
}

func mergeString(current string, target string) string {
	if len(target) > 0 {
		return target
	}

	return current
}
func replaceProject(projects ...Project) Project {
	current := Project{}

	l := len(projects)

	if l > 0 {
		return projects[l-1]
	}

	return current
}
