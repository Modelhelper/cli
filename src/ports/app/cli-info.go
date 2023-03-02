package app

import (
	"fmt"
	"math/rand"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/ports/config"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gookit/color"
)

type appInfo struct {
	name    string
	version string
	isBeta  bool
}

// Version implements modelhelper.AppInfoService
func (a *appInfo) Version() string {
	return a.version
}

// About implements modelhelper.AppInfoService
func (a *appInfo) About() string {
	infoElement := `
  Code
  ModelHelper CLI is a Command Line Interface tool to generate code based on an input source
  like a database table, REST api endpoint, a GraphQL endpoint or a proto file.

  Templates
  You can create your own templates based on Golang template ... each template is specified in a
  yaml- file and placed in a folder structure.

  Data
  This CLI can also help you understand database tables and perform some database tasks
  It works with MS SQL and Postgres. 

  Other input sources
  An input source can be either a database table or a set of tables. But it can also be a REST endpoint or graphql
  endpoint
  
  Application
  ------------
  Name:           mh 
  Version:        %v
  Location:       %v
  Environment:    %v
  Architecture:   %v
  Compiler:       %v
  Language:       go (version: %v)
  
  Config
  ------------
  Location:       %v
  `
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	// user := "Hans-Petter Eitvet" // nødvendig? - lagres i så fall i config
	gos := runtime.GOOS
	gar := runtime.GOARCH
	gv := runtime.Version()
	gc := runtime.Compiler

	cl := config.Location()

	return fmt.Sprintf(infoElement, a.version, exPath, gos, gar, gc, gv, cl)
}

// Logo implements modelhelper.AppInfoService
func (a *appInfo) Logo() string {
	var logo = color.Green.Sprintf(`
888b     d888               888          888 888    888          888                           
8888b   d8888               888          888 888    888          888                           
88888b.d88888               888          888 888    888          888                           
888Y88888P888  .d88b.   .d88888  .d88b.  888 8888888888  .d88b.  888 88888b.   .d88b.  888d888 
888 Y888P 888 d88""88b d88" 888 d8P  Y8b 888 888    888 d8P  Y8b 888 888 "88b d8P  Y8b 888P"   
888  Y8P  888 888  888 888  888 88888888 888 888    888 88888888 888 888  888 88888888 888     
888   "   888 Y88..88P Y88b 888 Y8b.     888 888    888 Y8b.     888 888 d88P Y8b.     888     
888       888  "Y88P"   "Y88888  "Y8888  888 888    888  "Y8888  888 88888P"   "Y8888  888     
                                                                     888                       
                                                                     888                       
                                                                     888   CLI v%v             
`, a.version)
	return logo
}

func NewCliInfo(name, version string, isBeta bool) modelhelper.AppInfoService {
	return &appInfo{name, version, isBeta}
}

// Slogan implements modelhelper.AppInfoService
func (*appInfo) Slogan() string {
	out := `
%s
`
	msg := fmt.Sprintf("'ModelHelper'\033[90m the\033[0m \033[32m%s\033[0m \033[90mhelper...\033[0m", RandAdjective())
	slogan := fmt.Sprintf(randBand(), msg)

	return fmt.Sprintf(out, slogan)
}

func positivityList() []string {
	return []string{
		"accomplished",
		"accurate",
		"adaptable",
		"adept",
		"adventurous",
		"affectionate",
		"agreeable",
		"alluring",
		"amazing",
		"ambitious",
		"amiable",
		"ample",
		"approachable",
		"articulate",
		"awesome",
		"blithesome",
		"bountiful",
		"brave",
		"bright",
		"brilliant",
		"capable",
		"captivating",
		"charismatic",
		"charming",
		"coherent",
		"colorful",
		"competitive",
		"concise",
		"confident",
		"considerate",
		"cool",
		"courageous",
		"creative",
		"credible",
		"cuddly",
		"cultivated",
		"cushy",
		"darling",
		"dashing",
		"dazzling",
		"decent",
		"decorous",
		"dedicated",
		"deliberate",
		"delightful",
		"demonstrative",
		"dependable",
		"determined",
		"devoted",
		"diligent",
		"diplomatic",
		"disarming",
		"distinguished",
		"dynamic",
		"eager",
		"educated",
		"efficient",
		"effortless",
		"electric",
		"elegant",
		"enchanting",
		"enduring",
		"energetic",
		"engaging",
		"enormous",
		"enriching",
		"enthusiastic",
		"excellent",
		"expert",
		"exuberant",
		"fabulous",
		"faithful",
		"fancy",
		"fantastic",
		"far-sighted",
		"fascinating",
		"fast",
		"faultless",
		"favorable",
		"favorite",
		"fearless",
		"flamboyant",
		"flexible",
		"focused",
		"forgiving",
		"fortuitous",
		"frank",
		"friendly",
		"fruitful",
		"fulfilling",
		"funny",
		"futuristic",
		"generous",
		"gentle",
		"giving",
		"gleaming",
		"gleeful",
		"glimmering",
		"glistening",
		"glittering",
		"glowing",
		"good-humored",
		"good-looking",
		"goodhearted",
		"gorgeous",
		"graceful",
		"greathearted",
		"gregarious",
		"hard-working",
		"hardworking",
		"harmonious",
		"helpful",
		"heroic",
		"high-powered",
		"honest",
		"hopeful",
		"humble",
		"humorous",
		"idealistic",
		"imaginative",
		"immediate",
		"impeccable",
		"incredible",
		"indefatigable",
		"independent",
		"innocent",
		"innovative",
		"inquisitive",
		"insightful",
		"jazzy",
		"jiggish",
		"kind",
		"kind",
		"knowable",
		"knowledgeable",
		"likable",
		"lionhearted",
		"lovely",
		"loving",
		"loyal",
		"luminous",
		"lustrous",
		"magnificent",
		"magnificentv",
		"marvelous",
		"marvelous",
		"mirthful",
		"motivational",
		"nice",
		"open-minded",
		"optimistic",
		"organized",
		"outstanding",
		"passionate",
		"patient",
		"perfect",
		"persistent",
		"personable",
		"philosophical",
		"plucky",
		"polite",
		"powerful",
		"powerful",
		"practical",
		"productive",
		"proficient",
		"propitious",
		"qualified",
		"rational",
		"ravishing",
		"relaxed",
		"remarkable",
		"resourceful",
		"responsible",
		"romantic",
		"rousing",
		"self-confident",
		"sensible",
		"shimmering",
		"sincere",
		"sleek",
		"sparkling",
		"spectacular",
		"splendid",
		"stellar",
		"stunning",
		"stupendous",
		"super",
		"technological",
		"thoughtful",
		"twinkling",
		"unique",
		"upbeat",
		"vibrant",
		"vivacious",
		"vivid",
		"warmhearted",
		"willing",
		"wondrous",
	}
}

func elementList() []string {
	return nil
}

func RandAdjective() string {
	slogans := positivityList()

	rs := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rs)

	i := r.Intn(len(slogans))

	return slogans[i]
}

func randBand() string {
	bands := []string{
		"\033[90m(>'-')> <('_'<) ^('_')\\- \\m/(-_-)\\m/ <( '-')> \\_( .')> < ( ._.)\033[0m\n%s",
		"\033[90m,.-~*´¨¯¨`*·~-.¸-(\033[0m %s \033[90m)-,.-~*´¨¯¨`*·~-.¸\033[0m",
		"\033[90m¸¸♬·¯·♩¸¸♪·¯·♫¸¸\033[0m %s \033[90m¸¸♬·¯·♩¸¸♪·¯·♫¸¸\033[0m",
		"\033[90m-=iii=<()  ♪·¯·♫¸\033[0m  %s",
		"\033[90m-=iii=<()  ♪·¯·♫¸\033[0m  %s  \033[90m¸♫·¯·♪  ()>=iii=-\033[0m",
		"\033[90m(¯`·._.·(¯`·._.·(¯`·._.·\033[0m  %s  \033[90m·._.·´¯)·._.·´¯)·._.·´¯)\033[0m",
		"\033[90m,.-~*´¨¯¨`*·~-.¸-(_\033[0m%s\033[90m_)-,.-~*´¨¯¨`*·~-.¸\033[0m",
		"\033[90m––•–√\\/––√\\/––•––\033[0m %s \033[90m––•–√\\/––√\\/––•––\033[0m",
		// 		`┏(-_-)┛ ┗(-_-)┓ ┗(-_-)┛ ┏(-_-)┓
		//         %s
		// ┏(-_-)┛ ┗(-_-)┓ ┗(-_-)┛ ┏(-_-)┓
		// 		`,
	}

	rs := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rs)

	y := r.Intn(len(bands))

	return bands[y]
}

func (a *appInfo) Welcome() string {
	usersName := ""
	betaMsg := ""

	usr, err := user.Current()
	if err == nil && usr != nil && len(usr.Name) > 0 {
		usersName = fmt.Sprintf(", %s", usr.Name)
	}

	if a.isBeta {
		betaMsg = fmt.Sprintf("\n%s\n", betaWarning())
	}

	return fmt.Sprintf(`
Welcome%s to ModelHelper CLI v.%s

Code
ModelHelper is a CLI tool to generate code based on input sources
like a database table

Templates
Templates are made with the Golang template language. Each template is specified in a
yaml- file and placed in a folder structure.

Data
Understand MS SQL tables and perform some database tasks.
%s
`, usersName, a.version, betaMsg)
	/*
	   Other input sources
	   An input source can be either a database table or a set of tables. But it can also be a REST endpoint or graphql
	   endpoint
	*/
}

func betaWarning() string {
	return color.Red.Sprint(`
██╗    ██╗ █████╗ ██████╗ ███╗   ██╗██╗███╗   ██╗ ██████╗ 
██║    ██║██╔══██╗██╔══██╗████╗  ██║██║████╗  ██║██╔════╝ 
██║ █╗ ██║███████║██████╔╝██╔██╗ ██║██║██╔██╗ ██║██║  ███╗
██║███╗██║██╔══██║██╔══██╗██║╚██╗██║██║██║╚██╗██║██║   ██║
╚███╔███╔╝██║  ██║██║  ██║██║ ╚████║██║██║ ╚████║╚██████╔╝
 ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝╚═╝  ╚═══╝ ╚═════╝ 
                                                          	
This version is still in beta, expect error and strange logic.
The codebase is also chanching rapidly (not every day, but almost).

TODO:
- The command api is not finalized and may change
- Missing a few templates
- Language definitions is not done
- Optimization

Use as is and feel free to throw in any issues you may encounter,
that will help me a lot to move the code to a final state.

(any nice words is also welcome :-)`)
}
