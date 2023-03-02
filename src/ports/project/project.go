package project

import (
	"io/ioutil"
	"log"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const dirname string = ".modelhelper"

type defaultProject struct {
	path   string
	Config models.ProjectConfig
}

// return ProjectService
func NewProjectConfigService() modelhelper.ProjectConfigService {
	return &defaultProject{}
}
func (p *defaultProject) Load() (*models.ProjectConfig, error) {
	path := DefaultLocation()
	return loadProjectFromFile(path)

}
func (p *defaultProject) LoadFromFile(path string) (*models.ProjectConfig, error) {
	return loadProjectFromFile(path)

}
func (p *defaultProject) New() (*models.ProjectConfig, error) {
	return nil, nil
}
func (p *defaultProject) Save() error {

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

func (p *defaultProject) Exists() bool {

	pathInfo, err := os.Stat(p.path)

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

// func LoadProjects(path ...string) []defaultProject {
// 	l := []defaultProject{}

//		for _, p := range path {
//			project, _ := loadProjectFromFile(p)
//			l = append(l, *project)
//		}
//		return l
//	}
// func Load(path string) (*defaultProject, error) {
// 	if len(path) > 0 {
// 		pathInfo, err := os.Stat(path)
// 		if os.IsNotExist(err) || pathInfo.IsDir() {
// 			// log.Fatal("Project does not exits")
// 			return nil, err
// 		}

// 	} else {
// 		p, err := os.Getwd()
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		path = filepath.Join(p, dirname, "project.yaml")
// 	}

// 	f, err := loadProjectFromFile(path)

// 	return f, err

// }

func loadProjectFromFile(fileName string) (*models.ProjectConfig, error) {
	var p *models.ProjectConfig

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

// FindRelatedProjects gets a list of all the related projects by following the path from DefaultDir() to the volumeroot
// The returning list is in correct importance from least to most (the project nearest to DefaultDir())
func (p *defaultProject) FindReleatedProjects(startPath string) []string {
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

func (p *defaultProject) FindNearestProjectDir() (string, bool) {
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

// func JoinProject(joinType string, projects ...project) project {
// 	switch joinType {
// 	case "merge":
// 		return mergeProject(projects...)

// 	case "smart":
// 		return smartMergeProject(projects...)

// 	case "replace":
// 		return replaceProject(projects...)
// 	default:
// 		return smartMergeProject(projects...)

// 	}
// }
// func mergeProject(projects ...project) project {
// 	current := project{}

// 	conProv := []project{}
// 	for _, p := range projects {
// 		conProv = append(conProv, p)
// 	}

// 	current.Connections = source.JoinConnections("merge", conProv...)

// 	for _, proj := range projects {

// 		current.Code = proj.Code
// 	}

// 	return current
// }

// func smartMergeProject(projects ...project) project {
// 	current := project{}
// 	current.Options = make(map[string]string)
// 	conProv := []modelhelper.ConnectionProvider{}
// 	for _, p := range projects {
// 		conProv = append(conProv, &p)
// 	}

// 	current.Connections = source.JoinConnections("smart", conProv...)

// 	for _, proj := range projects {
// 		current.Name = mergeString(current.Name, proj.Name)
// 		// if len(proj.Name) > 0 {
// 		// 	current.Name = proj.Name
// 		// }
// 		current.Description = mergeString(current.Description, proj.Description)
// 		current.DefaultKey = mergeString(current.DefaultKey, proj.DefaultKey)
// 		current.DefaultSource = mergeString(current.DefaultSource, proj.DefaultSource)
// 		current.OwnerName = mergeString(current.OwnerName, proj.OwnerName)
// 		current.Language = mergeString(current.Language, proj.Language)

// 		for optKey, optVal := range proj.Options {
// 			current.Options[optKey] = optVal
// 		}

// 		current.Code = proj.Code
// 	}

// 	return current
// }

func mergeString(current string, target string) string {
	if len(target) > 0 {
		return target
	}

	return current
}
func replaceProject(projects ...defaultProject) defaultProject {
	current := defaultProject{}

	l := len(projects)

	if l > 0 {
		return projects[l-1]
	}

	return current
}
