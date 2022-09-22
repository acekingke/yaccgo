package builder

var goObjectTemplateStr string = `
/*Generator Code , do not modify*/
// Code header part 
{{.CodeHeader}}

{{if .HttpParser}}
import(
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
)
{{end}}
// const part
{{.ConstPart}}
{{ if .NeedPacked }}
// Terminal Size
const NTERMINALS = {{.NTerminals}}
{{end}}

var IsTrace bool = false

type Context struct {
	StackSym []StateSym
	Stackpos int
}

type ValType struct {
	// Union part
	{{.UnionPart}}
}
type StateSym struct {
	Yystate int // state
	
 	//sym val
	YySymIndex int
	//other
	ValType
}

{{ if .NeedPacked }}
 // It is NeedPacked  
 {{.PackAnalyTable}}
 func (s *StateSym) Action(a int) int {
	if StatePackOffset[s.Yystate]+a  < 0 {
		 return ERROR_ACTION
	}
	if StatePackOffset[s.Yystate]+a >= len(StackPackCheck) || 
		StackPackCheck[StatePackOffset[s.Yystate]+a] != s.Yystate {
		 if a > NTERMINALS {
			 return StackPackGotoDef[a - NTERMINALS - 1]
		 }else {
			 return StackPackActDef[s.Yystate]
		 }
	}else{
		return StatePackAction[StatePackOffset[s.Yystate]+a]
	}
}
 {{else}} 
 // It is not packed
 var StateActionArray = [][]int{
	{{.AnalyTable}}
}
func (s *StateSym) Action(a int) int {
	return StateActionArray[s.Yystate][a]
}

 {{ end }}

 func TraceShift(s *StateSym) {
	if IsTrace {
	fmt.Printf("Shift %s, push state %d\n", TraceTranslate(s.YySymIndex), s.Yystate)
	}
}

// Reduce function
func (c *Context)ReduceFunc(reduceIndex int) *StateSym {
	dollarDolar := &StateSym{}
	topIndex := c.Stackpos - 1
	switch reduceIndex {
		{{.ReduceFunc}}
	}
	return dollarDolar
}

// Push StateSym
func (c *Context) PushStateSym(state *StateSym) {
	TraceShift(state)
	if c.Stackpos >= len(c.StackSym) {
		c.StackSym = append(c.StackSym, *state)
	} else {
		c.StackSym[c.Stackpos] = *state
	}
	c.Stackpos++
}

// Pop num StateSym
func (c *Context) PopStateSym(num int) {
	c.Stackpos -= num
}

func (c *Context) ParserInit() {
	c.StackSym = append(c.StackSym, StateSym{
		Yystate:    0,
		YySymIndex: 1, //$
	})

	c.Stackpos = 1
}
func MakeParserContext() *Context{
	c := &Context{}
	c.ParserInit()
	return c
}
func (c *Context) Parser(input string) *ValType {
	var currentPos int = 0
	var val ValType
	lookAhead := fetchLookAhead(input, &val, &currentPos)
	for {

		if c.Stackpos == 0 {
			break
		}
		if c.Stackpos > len(c.StackSym) {
			break
		}
		s := &c.StackSym[c.Stackpos-1]
		a := s.Action(lookAhead)

		if a == ERROR_ACTION {
			panic(fmt.Sprintf("Grammar error near pos %d", currentPos) + ":" + TraceTranslate(lookAhead))
		} else if a == ACCEPT_ACTION {
			return &s.ValType
		} else {
			if a > 0 {
				// shift
				c.PushStateSym(&StateSym{
					Yystate:    a,
					YySymIndex: lookAhead,
					ValType:    val,
				})
				lookAhead = fetchLookAhead(input, &val, &currentPos)
			} else {
				reduceIndex := -a
				SymTy := c.ReduceFunc(reduceIndex)
				s := &c.StackSym[c.Stackpos-1]
				gotoState := s.Action(SymTy.YySymIndex)
				SymTy.Yystate = gotoState
				TraceReduce(reduceIndex, gotoState, TraceTranslate(lookAhead))
				c.PushStateSym(SymTy)
			}
		}
	}
	return nil
}

func fetchLookAhead(input string, val *ValType, pos *int) int {
	token := GetToken(input, val, pos)
	return translate(token)
}
func translate(c int) int {
	var conv int = 0
	switch c {
		{{.Translate}}
	}
	return conv
}

// Trace function for translate
func TraceTranslate(c int) string {
	var conv string = ""
	switch c {
{{.TranslateTrace}}
	}
	return conv
}
// Trace function for reduce
func TraceReduce(reduceIndex, s int, look string) {
	if IsTrace {
		switch reduceIndex {
{{.ReduceTrace}}
		}
	}
}
{{if .HttpParser}}
// Trace function for Ping
func TracePingFun(rest string) {
	var result map[string]interface{} = make(map[string]interface{})
	var stateStack []string
	var symStack []string
	var valueStack []string
	for i := 0; i < StackPointer; i++ {
		v := StateSymStack[i]
		stateStack = append(stateStack, fmt.Sprintf("%d", v.Yystate))
		symStack = append(symStack, TraceTranslate(v.YySymIndex))
		valueStack = append(valueStack, fmt.Sprintf("%v", v.ValType.val))
	}
	result["states"] = stateStack
	result["symbols"] = symStack
	result["values"] = valueStack
	result["rest"] = rest

	js, _ := json.Marshal(result)
	ochan <- string(js)
	<-schan
}

type PingType struct {
	Input string` + "`" + `json:"input"` + "`" + `
}

var schan chan string = make(chan string)
var ochan chan string = make(chan string)
var finished bool = true

// Ping Function
func handlerPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if finished {
		go ParserFun()
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	var ping PingType

	json.Unmarshal(reqBody, &ping)
	fmt.Println(ping.Input)
	schan <- ping.Input
	w.Write([]byte(<-ochan))
	fmt.Println(time.Now(), r.Method, r.RequestURI, r.UserAgent())
}

func ParserFun() {
	finished = false
	ParserInit()
	input := <-schan
	_ = Parser(input)
	finished = true
}

func main() {
	http.HandleFunc("/ping", handlerPing)

	fmt.Println("ping listening on 0.0.0.0, port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting ping server: ", err)
	}
}

{{end}}
// Code Last part
{{.CodeLast}}`
