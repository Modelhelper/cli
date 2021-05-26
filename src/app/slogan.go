package app

import (
	"fmt"
	"math/rand"
	"time"
)

func Slogan() string {
	out := `
%s
`
	msg := fmt.Sprintf("'ModelHelper' the %s helper...", RandAdjective())
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
		`(>'-')> <('_'<) ^('_')\- \m/(-_-)\m/ <( '-')> \_( .')> < ( ._.) -
		%s
		`,
		",.-~*´¨¯¨`*·~-.¸-( %s )-,.-~*´¨¯¨`*·~-.¸",
		"¸¸♬·¯·♩¸¸♪·¯·♫¸¸ %s ¸¸♬·¯·♩¸¸♪·¯·♫¸¸",
		"-=iii=<()  ♪·¯·♫¸  %s",
		"-=iii=<()  ♪·¯·♫¸  %s  ¸♫·¯·♪  ()>=iii=-",
		"(¯`·._.·(¯`·._.·(¯`·._.·  %s  ·._.·´¯)·._.·´¯)·._.·´¯)",
		",.-~*´¨¯¨`*·~-.¸-(_%s_)-,.-~*´¨¯¨`*·~-.¸",
		`┏(-_-)┛ ┗(-_-)┓ ┗(-_-)┛ ┏(-_-)┓
        %s
┏(-_-)┛ ┗(-_-)┓ ┗(-_-)┛ ┏(-_-)┓		
		`,
		`––•–√\/––√\/––•–– %s ––•–√\/––√\/––•––`,
	}

	rs := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rs)

	y := r.Intn(len(bands))

	return bands[y]
}
