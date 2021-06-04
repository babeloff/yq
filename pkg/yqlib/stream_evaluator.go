package yqlib

import (
	"bytes"
	"container/list"
	"io"
	"os"

	jprops "github.com/magiconair/properties"
	yaml "gopkg.in/yaml.v3"
)

// A yaml expression evaluator that runs the expression multiple times for each given yaml document.
// Uses less memory than loading all documents and running the expression once, but this cannot process
// cross document expressions.
type StreamEvaluator interface {
	EvaluateYaml(filename string, reader io.Reader, node *ExpressionNode, printer Printer) error
	EvaluateProps(filename string, reader io.Reader, node *ExpressionNode, printer Printer) error
	EvaluateFiles(expression string, filenames []string, printer Printer) error
	EvaluateNew(expression string, printer Printer) error
}

type streamEvaluator struct {
	inputType     InputTypeEnum
	treeNavigator DataTreeNavigator
	treeCreator   ExpressionParser
	fileIndex     int
}

func NewStreamEvaluator(inType InputTypeEnum) StreamEvaluator {
	return &streamEvaluator{
		inputType:     inType,
		treeNavigator: NewDataTreeNavigator(),
		treeCreator:   NewExpressionParser(),
	}
}

func (s *streamEvaluator) EvaluateNew(expression string, printer Printer) error {
	node, err := s.treeCreator.ParseExpression(expression)
	if err != nil {
		return err
	}
	candidateNode := &CandidateNode{
		Document:  0,
		Filename:  "",
		Node:      &yaml.Node{Tag: "!!null", Kind: yaml.ScalarNode},
		FileIndex: 0,
	}
	inputList := list.New()
	inputList.PushBack(candidateNode)

	result, errorParsing := s.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputList}, node)
	if errorParsing != nil {
		return errorParsing
	}
	return printer.PrintResults(result.MatchingNodes)
}

func (s *streamEvaluator) EvaluateFiles(expression string, filenames []string, printer Printer) error {

	node, err := s.treeCreator.ParseExpression(expression)
	if err != nil {
		return err
	}

	for _, filename := range filenames {
		reader, err := readStream(filename)
		if err != nil {
			return err
		}
		switch s.inputType {
		case FromYaml, FromJson:
			err = s.EvaluateYaml(filename, reader, node, printer)
		case FromProps:
			err = s.EvaluateProps(filename, reader, node, printer)
		}

		if err != nil {
			return err
		}

		switch reader := reader.(type) {
		case *os.File:
			safelyCloseFile(reader)
		}
	}
	return nil
}

func (s *streamEvaluator) EvaluateYaml(filename string, reader io.Reader, node *ExpressionNode, printer Printer) error {

	var currentIndex uint

	decoder := yaml.NewDecoder(reader)
	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			s.fileIndex = s.fileIndex + 1
			return nil
		} else if errorReading != nil {
			return errorReading
		}
		candidateNode := &CandidateNode{
			Document:  currentIndex,
			Filename:  filename,
			Node:      &dataBucket,
			FileIndex: s.fileIndex,
		}
		inputList := list.New()
		inputList.PushBack(candidateNode)

		result, errorParsing := s.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputList}, node)
		if errorParsing != nil {
			return errorParsing
		}
		err := printer.PrintResults(result.MatchingNodes)
		if err != nil {
			return err
		}
		currentIndex = currentIndex + 1
	}
}

func (s *streamEvaluator) EvaluateProps(filename string, reader io.Reader, node *ExpressionNode, printer Printer) error {
	var currentIndex uint
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	dataBucket, errorReading := jprops.Load(buf.Bytes(), jprops.UTF8)
	if errorReading == nil {
		return errorReading
	}

	for _, key := range dataBucket.Keys() {
		value, _ := dataBucket.Get(key)
		comment := dataBucket.GetComment(key)
		dataNode := yaml.Node{
			Kind:        0,
			Style:       0,
			Tag:         key,
			Value:       value,
			Anchor:      "",
			Alias:       nil,
			Content:     nil,
			HeadComment: comment,
			LineComment: "",
			FootComment: "",
			Line:        0,
			Column:      0,
		}
		candidateNode := &CandidateNode{
			Document:  currentIndex,
			Filename:  filename,
			Node:      &dataNode,
			FileIndex: s.fileIndex,
		}
		inputList := list.New()
		inputList.PushBack(candidateNode)

		result, errorParsing := s.treeNavigator.GetMatchingNodes(Context{MatchingNodes: inputList}, node)
		if errorParsing != nil {
			return errorParsing
		}
		err := printer.PrintResults(result.MatchingNodes)
		if err != nil {
			return err
		}
		currentIndex = currentIndex + 1
	}
	return nil
}
