// Code generated from XMLParser.g4 by ANTLR 4.12.0. DO NOT EDIT.

package parser // XMLParser

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type XMLParser struct {
	*antlr.BaseParser
}

var xmlparserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func xmlparserParserInit() {
	staticData := &xmlparserParserStaticData
	staticData.literalNames = []string{
		"", "", "", "", "", "", "", "'<'", "", "", "'>'", "", "'/>'", "'/'",
		"'='",
	}
	staticData.symbolicNames = []string{
		"", "COMMENT", "CDATA", "DTD", "EntityRef", "CharRef", "SEA_WS", "OPEN",
		"XMLDeclOpen", "TEXT", "CLOSE", "SPECIAL_CLOSE", "SLASH_CLOSE", "SLASH",
		"EQUALS", "STRING", "Name", "S", "PI",
	}
	staticData.ruleNames = []string{
		"document", "prolog", "content", "element", "reference", "attribute",
		"chardata", "misc",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 18, 98, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 1, 0, 3, 0, 18, 8, 0, 1, 0, 5, 0,
		21, 8, 0, 10, 0, 12, 0, 24, 9, 0, 1, 0, 1, 0, 5, 0, 28, 8, 0, 10, 0, 12,
		0, 31, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 5, 1, 37, 8, 1, 10, 1, 12, 1, 40,
		9, 1, 1, 1, 1, 1, 1, 2, 3, 2, 45, 8, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3,
		2, 52, 8, 2, 1, 2, 3, 2, 55, 8, 2, 5, 2, 57, 8, 2, 10, 2, 12, 2, 60, 9,
		2, 1, 3, 1, 3, 1, 3, 5, 3, 65, 8, 3, 10, 3, 12, 3, 68, 9, 3, 1, 3, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 5, 3, 80, 8, 3, 10, 3,
		12, 3, 83, 9, 3, 1, 3, 3, 3, 86, 8, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1,
		5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 0, 0, 8, 0, 2, 4, 6, 8, 10, 12, 14, 0,
		3, 1, 0, 4, 5, 2, 0, 6, 6, 9, 9, 3, 0, 1, 1, 6, 6, 18, 18, 103, 0, 17,
		1, 0, 0, 0, 2, 34, 1, 0, 0, 0, 4, 44, 1, 0, 0, 0, 6, 85, 1, 0, 0, 0, 8,
		87, 1, 0, 0, 0, 10, 89, 1, 0, 0, 0, 12, 93, 1, 0, 0, 0, 14, 95, 1, 0, 0,
		0, 16, 18, 3, 2, 1, 0, 17, 16, 1, 0, 0, 0, 17, 18, 1, 0, 0, 0, 18, 22,
		1, 0, 0, 0, 19, 21, 3, 14, 7, 0, 20, 19, 1, 0, 0, 0, 21, 24, 1, 0, 0, 0,
		22, 20, 1, 0, 0, 0, 22, 23, 1, 0, 0, 0, 23, 25, 1, 0, 0, 0, 24, 22, 1,
		0, 0, 0, 25, 29, 3, 6, 3, 0, 26, 28, 3, 14, 7, 0, 27, 26, 1, 0, 0, 0, 28,
		31, 1, 0, 0, 0, 29, 27, 1, 0, 0, 0, 29, 30, 1, 0, 0, 0, 30, 32, 1, 0, 0,
		0, 31, 29, 1, 0, 0, 0, 32, 33, 5, 0, 0, 1, 33, 1, 1, 0, 0, 0, 34, 38, 5,
		8, 0, 0, 35, 37, 3, 10, 5, 0, 36, 35, 1, 0, 0, 0, 37, 40, 1, 0, 0, 0, 38,
		36, 1, 0, 0, 0, 38, 39, 1, 0, 0, 0, 39, 41, 1, 0, 0, 0, 40, 38, 1, 0, 0,
		0, 41, 42, 5, 11, 0, 0, 42, 3, 1, 0, 0, 0, 43, 45, 3, 12, 6, 0, 44, 43,
		1, 0, 0, 0, 44, 45, 1, 0, 0, 0, 45, 58, 1, 0, 0, 0, 46, 52, 3, 6, 3, 0,
		47, 52, 3, 8, 4, 0, 48, 52, 5, 2, 0, 0, 49, 52, 5, 18, 0, 0, 50, 52, 5,
		1, 0, 0, 51, 46, 1, 0, 0, 0, 51, 47, 1, 0, 0, 0, 51, 48, 1, 0, 0, 0, 51,
		49, 1, 0, 0, 0, 51, 50, 1, 0, 0, 0, 52, 54, 1, 0, 0, 0, 53, 55, 3, 12,
		6, 0, 54, 53, 1, 0, 0, 0, 54, 55, 1, 0, 0, 0, 55, 57, 1, 0, 0, 0, 56, 51,
		1, 0, 0, 0, 57, 60, 1, 0, 0, 0, 58, 56, 1, 0, 0, 0, 58, 59, 1, 0, 0, 0,
		59, 5, 1, 0, 0, 0, 60, 58, 1, 0, 0, 0, 61, 62, 5, 7, 0, 0, 62, 66, 5, 16,
		0, 0, 63, 65, 3, 10, 5, 0, 64, 63, 1, 0, 0, 0, 65, 68, 1, 0, 0, 0, 66,
		64, 1, 0, 0, 0, 66, 67, 1, 0, 0, 0, 67, 69, 1, 0, 0, 0, 68, 66, 1, 0, 0,
		0, 69, 70, 5, 10, 0, 0, 70, 71, 3, 4, 2, 0, 71, 72, 5, 7, 0, 0, 72, 73,
		5, 13, 0, 0, 73, 74, 5, 16, 0, 0, 74, 75, 5, 10, 0, 0, 75, 86, 1, 0, 0,
		0, 76, 77, 5, 7, 0, 0, 77, 81, 5, 16, 0, 0, 78, 80, 3, 10, 5, 0, 79, 78,
		1, 0, 0, 0, 80, 83, 1, 0, 0, 0, 81, 79, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0,
		82, 84, 1, 0, 0, 0, 83, 81, 1, 0, 0, 0, 84, 86, 5, 12, 0, 0, 85, 61, 1,
		0, 0, 0, 85, 76, 1, 0, 0, 0, 86, 7, 1, 0, 0, 0, 87, 88, 7, 0, 0, 0, 88,
		9, 1, 0, 0, 0, 89, 90, 5, 16, 0, 0, 90, 91, 5, 14, 0, 0, 91, 92, 5, 15,
		0, 0, 92, 11, 1, 0, 0, 0, 93, 94, 7, 1, 0, 0, 94, 13, 1, 0, 0, 0, 95, 96,
		7, 2, 0, 0, 96, 15, 1, 0, 0, 0, 11, 17, 22, 29, 38, 44, 51, 54, 58, 66,
		81, 85,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// XMLParserInit initializes any static state used to implement XMLParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewXMLParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func XMLParserInit() {
	staticData := &xmlparserParserStaticData
	staticData.once.Do(xmlparserParserInit)
}

// NewXMLParser produces a new parser instance for the optional input antlr.TokenStream.
func NewXMLParser(input antlr.TokenStream) *XMLParser {
	XMLParserInit()
	this := new(XMLParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &xmlparserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "XMLParser.g4"

	return this
}

// XMLParser tokens.
const (
	XMLParserEOF           = antlr.TokenEOF
	XMLParserCOMMENT       = 1
	XMLParserCDATA         = 2
	XMLParserDTD           = 3
	XMLParserEntityRef     = 4
	XMLParserCharRef       = 5
	XMLParserSEA_WS        = 6
	XMLParserOPEN          = 7
	XMLParserXMLDeclOpen   = 8
	XMLParserTEXT          = 9
	XMLParserCLOSE         = 10
	XMLParserSPECIAL_CLOSE = 11
	XMLParserSLASH_CLOSE   = 12
	XMLParserSLASH         = 13
	XMLParserEQUALS        = 14
	XMLParserSTRING        = 15
	XMLParserName          = 16
	XMLParserS             = 17
	XMLParserPI            = 18
)

// XMLParser rules.
const (
	XMLParserRULE_document  = 0
	XMLParserRULE_prolog    = 1
	XMLParserRULE_content   = 2
	XMLParserRULE_element   = 3
	XMLParserRULE_reference = 4
	XMLParserRULE_attribute = 5
	XMLParserRULE_chardata  = 6
	XMLParserRULE_misc      = 7
)

// IDocumentContext is an interface to support dynamic dispatch.
type IDocumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Element() IElementContext
	EOF() antlr.TerminalNode
	Prolog() IPrologContext
	AllMisc() []IMiscContext
	Misc(i int) IMiscContext

	// IsDocumentContext differentiates from other interfaces.
	IsDocumentContext()
}

type DocumentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDocumentContext() *DocumentContext {
	var p = new(DocumentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = XMLParserRULE_document
	return p
}

func (*DocumentContext) IsDocumentContext() {}

func NewDocumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DocumentContext {
	var p = new(DocumentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = XMLParserRULE_document

	return p
}

func (s *DocumentContext) GetParser() antlr.Parser { return s.parser }

func (s *DocumentContext) Element() IElementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IElementContext)
}

func (s *DocumentContext) EOF() antlr.TerminalNode {
	return s.GetToken(XMLParserEOF, 0)
}

func (s *DocumentContext) Prolog() IPrologContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrologContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrologContext)
}

func (s *DocumentContext) AllMisc() []IMiscContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IMiscContext); ok {
			len++
		}
	}

	tst := make([]IMiscContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IMiscContext); ok {
			tst[i] = t.(IMiscContext)
			i++
		}
	}

	return tst
}

func (s *DocumentContext) Misc(i int) IMiscContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMiscContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMiscContext)
}

func (s *DocumentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DocumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DocumentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.EnterDocument(s)
	}
}

func (s *DocumentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.ExitDocument(s)
	}
}

func (p *XMLParser) Document() (localctx IDocumentContext) {
	this := p
	_ = this

	localctx = NewDocumentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, XMLParserRULE_document)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(17)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == XMLParserXMLDeclOpen {
		{
			p.SetState(16)
			p.Prolog()
		}

	}
	p.SetState(22)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&262210) != 0 {
		{
			p.SetState(19)
			p.Misc()
		}

		p.SetState(24)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(25)
		p.Element()
	}
	p.SetState(29)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&262210) != 0 {
		{
			p.SetState(26)
			p.Misc()
		}

		p.SetState(31)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(32)
		p.Match(XMLParserEOF)
	}

	return localctx
}

// IPrologContext is an interface to support dynamic dispatch.
type IPrologContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	XMLDeclOpen() antlr.TerminalNode
	SPECIAL_CLOSE() antlr.TerminalNode
	AllAttribute() []IAttributeContext
	Attribute(i int) IAttributeContext

	// IsPrologContext differentiates from other interfaces.
	IsPrologContext()
}

type PrologContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrologContext() *PrologContext {
	var p = new(PrologContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = XMLParserRULE_prolog
	return p
}

func (*PrologContext) IsPrologContext() {}

func NewPrologContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrologContext {
	var p = new(PrologContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = XMLParserRULE_prolog

	return p
}

func (s *PrologContext) GetParser() antlr.Parser { return s.parser }

func (s *PrologContext) XMLDeclOpen() antlr.TerminalNode {
	return s.GetToken(XMLParserXMLDeclOpen, 0)
}

func (s *PrologContext) SPECIAL_CLOSE() antlr.TerminalNode {
	return s.GetToken(XMLParserSPECIAL_CLOSE, 0)
}

func (s *PrologContext) AllAttribute() []IAttributeContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAttributeContext); ok {
			len++
		}
	}

	tst := make([]IAttributeContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAttributeContext); ok {
			tst[i] = t.(IAttributeContext)
			i++
		}
	}

	return tst
}

func (s *PrologContext) Attribute(i int) IAttributeContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAttributeContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAttributeContext)
}

func (s *PrologContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrologContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrologContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.EnterProlog(s)
	}
}

func (s *PrologContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.ExitProlog(s)
	}
}

func (p *XMLParser) Prolog() (localctx IPrologContext) {
	this := p
	_ = this

	localctx = NewPrologContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, XMLParserRULE_prolog)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(34)
		p.Match(XMLParserXMLDeclOpen)
	}
	p.SetState(38)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == XMLParserName {
		{
			p.SetState(35)
			p.Attribute()
		}

		p.SetState(40)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(41)
		p.Match(XMLParserSPECIAL_CLOSE)
	}

	return localctx
}

// IContentContext is an interface to support dynamic dispatch.
type IContentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllChardata() []IChardataContext
	Chardata(i int) IChardataContext
	AllElement() []IElementContext
	Element(i int) IElementContext
	AllReference() []IReferenceContext
	Reference(i int) IReferenceContext
	AllCDATA() []antlr.TerminalNode
	CDATA(i int) antlr.TerminalNode
	AllPI() []antlr.TerminalNode
	PI(i int) antlr.TerminalNode
	AllCOMMENT() []antlr.TerminalNode
	COMMENT(i int) antlr.TerminalNode

	// IsContentContext differentiates from other interfaces.
	IsContentContext()
}

type ContentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyContentContext() *ContentContext {
	var p = new(ContentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = XMLParserRULE_content
	return p
}

func (*ContentContext) IsContentContext() {}

func NewContentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContentContext {
	var p = new(ContentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = XMLParserRULE_content

	return p
}

func (s *ContentContext) GetParser() antlr.Parser { return s.parser }

func (s *ContentContext) AllChardata() []IChardataContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IChardataContext); ok {
			len++
		}
	}

	tst := make([]IChardataContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IChardataContext); ok {
			tst[i] = t.(IChardataContext)
			i++
		}
	}

	return tst
}

func (s *ContentContext) Chardata(i int) IChardataContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IChardataContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IChardataContext)
}

func (s *ContentContext) AllElement() []IElementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IElementContext); ok {
			len++
		}
	}

	tst := make([]IElementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IElementContext); ok {
			tst[i] = t.(IElementContext)
			i++
		}
	}

	return tst
}

func (s *ContentContext) Element(i int) IElementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IElementContext)
}

func (s *ContentContext) AllReference() []IReferenceContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IReferenceContext); ok {
			len++
		}
	}

	tst := make([]IReferenceContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IReferenceContext); ok {
			tst[i] = t.(IReferenceContext)
			i++
		}
	}

	return tst
}

func (s *ContentContext) Reference(i int) IReferenceContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReferenceContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IReferenceContext)
}

func (s *ContentContext) AllCDATA() []antlr.TerminalNode {
	return s.GetTokens(XMLParserCDATA)
}

func (s *ContentContext) CDATA(i int) antlr.TerminalNode {
	return s.GetToken(XMLParserCDATA, i)
}

func (s *ContentContext) AllPI() []antlr.TerminalNode {
	return s.GetTokens(XMLParserPI)
}

func (s *ContentContext) PI(i int) antlr.TerminalNode {
	return s.GetToken(XMLParserPI, i)
}

func (s *ContentContext) AllCOMMENT() []antlr.TerminalNode {
	return s.GetTokens(XMLParserCOMMENT)
}

func (s *ContentContext) COMMENT(i int) antlr.TerminalNode {
	return s.GetToken(XMLParserCOMMENT, i)
}

func (s *ContentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ContentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.EnterContent(s)
	}
}

func (s *ContentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.ExitContent(s)
	}
}

func (p *XMLParser) Content() (localctx IContentContext) {
	this := p
	_ = this

	localctx = NewContentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, XMLParserRULE_content)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(44)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == XMLParserSEA_WS || _la == XMLParserTEXT {
		{
			p.SetState(43)
			p.Chardata()
		}

	}
	p.SetState(58)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(51)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case XMLParserOPEN:
				{
					p.SetState(46)
					p.Element()
				}

			case XMLParserEntityRef, XMLParserCharRef:
				{
					p.SetState(47)
					p.Reference()
				}

			case XMLParserCDATA:
				{
					p.SetState(48)
					p.Match(XMLParserCDATA)
				}

			case XMLParserPI:
				{
					p.SetState(49)
					p.Match(XMLParserPI)
				}

			case XMLParserCOMMENT:
				{
					p.SetState(50)
					p.Match(XMLParserCOMMENT)
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}
			p.SetState(54)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == XMLParserSEA_WS || _la == XMLParserTEXT {
				{
					p.SetState(53)
					p.Chardata()
				}

			}

		}
		p.SetState(60)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())
	}

	return localctx
}

// IElementContext is an interface to support dynamic dispatch.
type IElementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllOPEN() []antlr.TerminalNode
	OPEN(i int) antlr.TerminalNode
	AllName() []antlr.TerminalNode
	Name(i int) antlr.TerminalNode
	AllCLOSE() []antlr.TerminalNode
	CLOSE(i int) antlr.TerminalNode
	Content() IContentContext
	SLASH() antlr.TerminalNode
	AllAttribute() []IAttributeContext
	Attribute(i int) IAttributeContext
	SLASH_CLOSE() antlr.TerminalNode

	// IsElementContext differentiates from other interfaces.
	IsElementContext()
}

type ElementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyElementContext() *ElementContext {
	var p = new(ElementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = XMLParserRULE_element
	return p
}

func (*ElementContext) IsElementContext() {}

func NewElementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElementContext {
	var p = new(ElementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = XMLParserRULE_element

	return p
}

func (s *ElementContext) GetParser() antlr.Parser { return s.parser }

func (s *ElementContext) AllOPEN() []antlr.TerminalNode {
	return s.GetTokens(XMLParserOPEN)
}

func (s *ElementContext) OPEN(i int) antlr.TerminalNode {
	return s.GetToken(XMLParserOPEN, i)
}

func (s *ElementContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(XMLParserName)
}

func (s *ElementContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(XMLParserName, i)
}

func (s *ElementContext) AllCLOSE() []antlr.TerminalNode {
	return s.GetTokens(XMLParserCLOSE)
}

func (s *ElementContext) CLOSE(i int) antlr.TerminalNode {
	return s.GetToken(XMLParserCLOSE, i)
}

func (s *ElementContext) Content() IContentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IContentContext)
}

func (s *ElementContext) SLASH() antlr.TerminalNode {
	return s.GetToken(XMLParserSLASH, 0)
}

func (s *ElementContext) AllAttribute() []IAttributeContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAttributeContext); ok {
			len++
		}
	}

	tst := make([]IAttributeContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAttributeContext); ok {
			tst[i] = t.(IAttributeContext)
			i++
		}
	}

	return tst
}

func (s *ElementContext) Attribute(i int) IAttributeContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAttributeContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAttributeContext)
}

func (s *ElementContext) SLASH_CLOSE() antlr.TerminalNode {
	return s.GetToken(XMLParserSLASH_CLOSE, 0)
}

func (s *ElementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ElementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ElementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.EnterElement(s)
	}
}

func (s *ElementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.ExitElement(s)
	}
}

func (p *XMLParser) Element() (localctx IElementContext) {
	this := p
	_ = this

	localctx = NewElementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, XMLParserRULE_element)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(85)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(61)
			p.Match(XMLParserOPEN)
		}
		{
			p.SetState(62)
			p.Match(XMLParserName)
		}
		p.SetState(66)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == XMLParserName {
			{
				p.SetState(63)
				p.Attribute()
			}

			p.SetState(68)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(69)
			p.Match(XMLParserCLOSE)
		}
		{
			p.SetState(70)
			p.Content()
		}
		{
			p.SetState(71)
			p.Match(XMLParserOPEN)
		}
		{
			p.SetState(72)
			p.Match(XMLParserSLASH)
		}
		{
			p.SetState(73)
			p.Match(XMLParserName)
		}
		{
			p.SetState(74)
			p.Match(XMLParserCLOSE)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(76)
			p.Match(XMLParserOPEN)
		}
		{
			p.SetState(77)
			p.Match(XMLParserName)
		}
		p.SetState(81)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == XMLParserName {
			{
				p.SetState(78)
				p.Attribute()
			}

			p.SetState(83)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(84)
			p.Match(XMLParserSLASH_CLOSE)
		}

	}

	return localctx
}

// IReferenceContext is an interface to support dynamic dispatch.
type IReferenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EntityRef() antlr.TerminalNode
	CharRef() antlr.TerminalNode

	// IsReferenceContext differentiates from other interfaces.
	IsReferenceContext()
}

type ReferenceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReferenceContext() *ReferenceContext {
	var p = new(ReferenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = XMLParserRULE_reference
	return p
}

func (*ReferenceContext) IsReferenceContext() {}

func NewReferenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReferenceContext {
	var p = new(ReferenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = XMLParserRULE_reference

	return p
}

func (s *ReferenceContext) GetParser() antlr.Parser { return s.parser }

func (s *ReferenceContext) EntityRef() antlr.TerminalNode {
	return s.GetToken(XMLParserEntityRef, 0)
}

func (s *ReferenceContext) CharRef() antlr.TerminalNode {
	return s.GetToken(XMLParserCharRef, 0)
}

func (s *ReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReferenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.EnterReference(s)
	}
}

func (s *ReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.ExitReference(s)
	}
}

func (p *XMLParser) Reference() (localctx IReferenceContext) {
	this := p
	_ = this

	localctx = NewReferenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, XMLParserRULE_reference)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(87)
		_la = p.GetTokenStream().LA(1)

		if !(_la == XMLParserEntityRef || _la == XMLParserCharRef) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IAttributeContext is an interface to support dynamic dispatch.
type IAttributeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Name() antlr.TerminalNode
	EQUALS() antlr.TerminalNode
	STRING() antlr.TerminalNode

	// IsAttributeContext differentiates from other interfaces.
	IsAttributeContext()
}

type AttributeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttributeContext() *AttributeContext {
	var p = new(AttributeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = XMLParserRULE_attribute
	return p
}

func (*AttributeContext) IsAttributeContext() {}

func NewAttributeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributeContext {
	var p = new(AttributeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = XMLParserRULE_attribute

	return p
}

func (s *AttributeContext) GetParser() antlr.Parser { return s.parser }

func (s *AttributeContext) Name() antlr.TerminalNode {
	return s.GetToken(XMLParserName, 0)
}

func (s *AttributeContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(XMLParserEQUALS, 0)
}

func (s *AttributeContext) STRING() antlr.TerminalNode {
	return s.GetToken(XMLParserSTRING, 0)
}

func (s *AttributeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AttributeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AttributeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.EnterAttribute(s)
	}
}

func (s *AttributeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.ExitAttribute(s)
	}
}

func (p *XMLParser) Attribute() (localctx IAttributeContext) {
	this := p
	_ = this

	localctx = NewAttributeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, XMLParserRULE_attribute)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(89)
		p.Match(XMLParserName)
	}
	{
		p.SetState(90)
		p.Match(XMLParserEQUALS)
	}
	{
		p.SetState(91)
		p.Match(XMLParserSTRING)
	}

	return localctx
}

// IChardataContext is an interface to support dynamic dispatch.
type IChardataContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TEXT() antlr.TerminalNode
	SEA_WS() antlr.TerminalNode

	// IsChardataContext differentiates from other interfaces.
	IsChardataContext()
}

type ChardataContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyChardataContext() *ChardataContext {
	var p = new(ChardataContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = XMLParserRULE_chardata
	return p
}

func (*ChardataContext) IsChardataContext() {}

func NewChardataContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ChardataContext {
	var p = new(ChardataContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = XMLParserRULE_chardata

	return p
}

func (s *ChardataContext) GetParser() antlr.Parser { return s.parser }

func (s *ChardataContext) TEXT() antlr.TerminalNode {
	return s.GetToken(XMLParserTEXT, 0)
}

func (s *ChardataContext) SEA_WS() antlr.TerminalNode {
	return s.GetToken(XMLParserSEA_WS, 0)
}

func (s *ChardataContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ChardataContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ChardataContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.EnterChardata(s)
	}
}

func (s *ChardataContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.ExitChardata(s)
	}
}

func (p *XMLParser) Chardata() (localctx IChardataContext) {
	this := p
	_ = this

	localctx = NewChardataContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, XMLParserRULE_chardata)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(93)
		_la = p.GetTokenStream().LA(1)

		if !(_la == XMLParserSEA_WS || _la == XMLParserTEXT) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IMiscContext is an interface to support dynamic dispatch.
type IMiscContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COMMENT() antlr.TerminalNode
	PI() antlr.TerminalNode
	SEA_WS() antlr.TerminalNode

	// IsMiscContext differentiates from other interfaces.
	IsMiscContext()
}

type MiscContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMiscContext() *MiscContext {
	var p = new(MiscContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = XMLParserRULE_misc
	return p
}

func (*MiscContext) IsMiscContext() {}

func NewMiscContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MiscContext {
	var p = new(MiscContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = XMLParserRULE_misc

	return p
}

func (s *MiscContext) GetParser() antlr.Parser { return s.parser }

func (s *MiscContext) COMMENT() antlr.TerminalNode {
	return s.GetToken(XMLParserCOMMENT, 0)
}

func (s *MiscContext) PI() antlr.TerminalNode {
	return s.GetToken(XMLParserPI, 0)
}

func (s *MiscContext) SEA_WS() antlr.TerminalNode {
	return s.GetToken(XMLParserSEA_WS, 0)
}

func (s *MiscContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MiscContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MiscContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.EnterMisc(s)
	}
}

func (s *MiscContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(XMLParserListener); ok {
		listenerT.ExitMisc(s)
	}
}

func (p *XMLParser) Misc() (localctx IMiscContext) {
	this := p
	_ = this

	localctx = NewMiscContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, XMLParserRULE_misc)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(95)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&262210) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}
