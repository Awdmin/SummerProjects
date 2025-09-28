package main


import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"os"
	"time"
	
	"github.com/gocolly/colly"
	"github.com/rivo/tview"
	"github.com/joho/godotenv"
)


type Lecture struct {
	Name		string
	Prof		string
	Classroom	string
	Type		string	
	Day			int
	UNI			string
	Start		int
	End			int
}

const FLOAT_PARSE_ERROR string = "Error with parsing float"
const INT_PARSE_ERROR string = "Error with parsing int"

var weekOffset int = 0

var DayToInt = map[string]int{
	"ponedeljek":	0,
	"torek": 		1,
	"sreda": 		2,
	"četrtek": 		3,
	"petek":		4,
}

var IntToDay = map[int]string{
	0:	"ponedeljek",
	1:	"torek",
	2:	"sreda",
	3:	"četrtek",
	4:	"petek",
}

var Lectures []Lecture

func checkLectureRow(row int, l Lecture) bool {
	if row+7 >= l.Start && row+7 <= l.End {
		return true
	}
	return false
}


func getCSSProp(style, prop string) string {
	for _, rule := range strings.Split(style, ";") {
		rule = strings.TrimSpace(rule)
		if strings.HasPrefix(rule, prop+":") {
			return strings.TrimSpace(strings.TrimPrefix(rule, prop+":"))
		}
	}
	return ""
}

func parseFRIObj(e *colly.HTMLElement) {

	lines := strings.Split(e.Text, "\n")

	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines[idx] = line
	}

	l4split := strings.Split(lines[4], "(")
	name := l4split[0]
	typ := strings.Split(l4split[1], "_")[1]

	teacher := lines[5]
	classroom := lines[2]

	time := lines[1]
	timesplit := strings.Split(time, " ")

	dayStr := timesplit[0]
	day := DayToInt[dayStr]
	
	startStr := timesplit[1]
	endStr := timesplit[3]

	start, err := strconv.Atoi(strings.Split(startStr, ":")[0])
	if err != nil {
		fmt.Println(INT_PARSE_ERROR, err)
		return 
	}

	end, err := strconv.Atoi(strings.Split(endStr, ":")[0])
	if err != nil {
		fmt.Println(INT_PARSE_ERROR, err)
		return 
	}

	l := Lecture{
		Name:		name,
		Prof:		teacher,	
		Classroom:	classroom,
		Type:		typ,	
		UNI:		"FRI",
		Day:		day,
		Start:		start,
		End:		end,
	}

	Lectures = append(Lectures, l)

}


func parseFMFObj(e *colly.HTMLElement) {
	name := e.ChildText("div.entry div.main-box span a.subject")
	typ := e.ChildText("div.entry div.main-box span span.entry-type")
	teacher := e.ChildText("div.entry div.teacher a")
	classroom := e.ChildText("div.entry div.classroom.classroom-box a")
	
	style := e.Attr("style")
	topPos := getCSSProp(style, "top")
	leftPos := getCSSProp(style, "left")
	height := getCSSProp(style, "height")

	topPos = strings.TrimSuffix(topPos, "%")
	leftPos = strings.TrimSuffix(leftPos, "%")
	height = strings.TrimSuffix(height, "%")

	top, err := strconv.ParseFloat(topPos, 64)
	if err != nil {
		fmt.Println(FLOAT_PARSE_ERROR, err)
		return 
	}

	hgt, err := strconv.ParseFloat(height, 64)
	if err != nil {
		fmt.Println(FLOAT_PARSE_ERROR, err)
		return
	}

	lft, err := strconv.ParseFloat(leftPos, 64)
	if err != nil {
		fmt.Println(FLOAT_PARSE_ERROR, err)
		return 
	}

	day := int(lft / 20.0)
	duration := int(hgt / 7.69)
	start := int(top / 7.69) + 7
	end := start + duration

	l := Lecture{
		Name:		name,
		Prof:		teacher,	
		Classroom:	classroom,
		Type:		typ,	
		UNI:		"FMF",
		Day:		day,
		Start:		start,
		End:		end,
	}

	Lectures = append(Lectures, l)
}


func fetchData(offset int, c *colly.Collector) {
	friURI := os.Getenv("FRI_URI")
	fmfURI := os.Getenv("FMF_URI")

	today := time.Now()
	date := today.AddDate(0, 0, 7*offset)
	formated := date.Format("2006-01-02")


	err := c.Visit(friURI)
	if err != nil {
		log.Fatal(err)
	}

		fmt.Println(formated)


	err = c.Visit(fmfURI+formated)
	if err != nil {
		log.Fatal(err)
	}

}


func fillTable(table *tview.Table) {
	for c := 0; c < 6; c++ {
		for r := 0; r < 12; r++ {
			table.SetCell(r, c, tview.NewTableCell(""))
		}
	}

	for t := 7; t < 20; t++ {
		table.SetCell(t-6, 0, tview.NewTableCell(fmt.Sprintf("%d:00", t)))
	}

	for d := 1; d < 6; d++ {
		table.SetCell(0, d, tview.NewTableCell(IntToDay[d-1]))
	}
	for _, lecture := range Lectures {
		for i := lecture.Start; i < lecture.End; i ++ { 
			color := "red"
			if lecture.UNI == "FMF" {
				color = "green"
			}
			content := fmt.Sprintf("%s [%s]%s", 
				lecture.Name, 
				color,
				lecture.UNI)

			cell := tview.NewTableCell(content).
				SetExpansion(1)

			table.SetCell(i-6, lecture.Day+1, cell)
		}
	}

}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env variables")
		return
	}

	if len(os.Args) > 1 {
		num := os.Args[1]
		n, err := strconv.Atoi(num)
		if err != nil {
			fmt.Println("Invalid argument (must be an int)", num) 
			return
		}
		weekOffset = n
	}

	fmfDomain := os.Getenv("FMF_DOMAIN")
	friDomain := os.Getenv("FRI_DOMAIN")
	c := colly.NewCollector(
		colly.AllowedDomains(fmfDomain, friDomain),
	)

	c.OnHTML("div.description div.entry-hover", parseFRIObj)

	c.OnHTML("div.entry-absolute-box", parseFMFObj)

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error: ", err)
	})
	

	fetchData(weekOffset, c)

	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)
	table.SetSelectable(true, true)
	
	
	fillTable(table)

	details := tview.NewTextView()
	details.SetDynamicColors(true)
	details.SetTextAlign(tview.AlignCenter)
	details.SetWrap(true)
	details.SetBorder(true)
	details.SetTitle("Details")
	
	table.SetSelectionChangedFunc(func(row, col int) {
		for _, lec := range Lectures {
			if col == lec.Day+1 && checkLectureRow(row, lec) {
				color := "green"
				if lec.UNI == "FRI" {
					color = "red"
				}
				content := fmt.Sprintf("\n\n[blue]%s\n[%s]%s\n[white]%s\nType: [orange]%s\n[white]%s",
					lec.Name,
					color,
					lec.UNI,
					lec.Classroom,
					lec.Type,
					lec.Prof,
				)

				details.SetText(content)
				break
			}
		}
	})
	
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(table, 0, 2, true).
		AddItem(details, 0, 1, false)


	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
