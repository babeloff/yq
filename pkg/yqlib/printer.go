package yqlib

import (
	"bufio"
	"container/list"
	"github.com/mikefarah/yq/v4/cmd"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type Printer interface {
	PrintResults(matchingNodes *list.List) error
	PrintedAnything() bool
}

type resultsPrinter struct {
	outputToProps      bool
	outputType         cmd.OutputTypeEnum
	unwrapScalar       bool
	colorsEnabled      bool
	indent             int
	printDocSeparators bool
	writer             io.Writer
	firstTimePrinting  bool
	previousDocIndex   uint
	previousFileIndex  int
	printedMatches     bool
	treeNavigator      DataTreeNavigator
}

func NewPrinter(writer io.Writer, outType cmd.OutputTypeEnum, unwrapScalar bool, colorsEnabled bool, indent int, printDocSeparators bool) Printer {
	return &resultsPrinter{
		writer:             writer,
		outputType:         outType,
		unwrapScalar:       unwrapScalar,
		colorsEnabled:      colorsEnabled,
		indent:             indent,
		printDocSeparators: outType == cmd.ToJson && printDocSeparators,
		firstTimePrinting:  true,
		treeNavigator:      NewDataTreeNavigator(),
	}
}

func (p *resultsPrinter) PrintedAnything() bool {
	return p.printedMatches
}

func (p *resultsPrinter) printNode(node *yaml.Node, writer io.Writer) error {
	p.printedMatches = p.printedMatches || (node.Tag != "!!null" &&
		(node.Tag != "!!bool" || node.Value != "false"))

	var encoder Encoder
	if node.Kind == yaml.ScalarNode && p.unwrapScalar && p.outputType != cmd.ToJson {
		return p.writeString(writer, node.Value+"\n")
	}
	if p.outputType == cmd.ToJson {
		encoder = NewJsonEncoder(writer, p.indent)
	} else {
		encoder = NewYamlEncoder(writer, p.indent, p.colorsEnabled)
	}
	return encoder.Encode(node)
}

func (p *resultsPrinter) writeString(writer io.Writer, txt string) error {
	_, errorWriting := writer.Write([]byte(txt))
	return errorWriting
}

func (p *resultsPrinter) safelyFlush(writer *bufio.Writer) {
	if err := writer.Flush(); err != nil {
		log.Error("Error flushing writer!")
		log.Error(err.Error())
	}
}

func (p *resultsPrinter) PrintResults(matchingNodes *list.List) error {
	log.Debug("PrintResults for %v matches", matchingNodes.Len())
	if p.outputType == cmd.ToJson {
		explodeOp := Operation{OperationType: explodeOpType}
		explodeNode := ExpressionNode{Operation: &explodeOp}
		context, err := p.treeNavigator.GetMatchingNodes(Context{MatchingNodes: matchingNodes}, &explodeNode)
		if err != nil {
			return err
		}
		matchingNodes = context.MatchingNodes
	}

	bufferedWriter := bufio.NewWriter(p.writer)
	defer p.safelyFlush(bufferedWriter)

	if matchingNodes.Len() == 0 {
		log.Debug("no matching results, nothing to print")
		return nil
	}
	if p.firstTimePrinting {
		node := matchingNodes.Front().Value.(*CandidateNode)
		p.previousDocIndex = node.Document
		p.previousFileIndex = node.FileIndex
		p.firstTimePrinting = false
	}

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		mappedDoc := el.Value.(*CandidateNode)
		log.Debug("-- print sep logic: p.firstTimePrinting: %v, previousDocIndex: %v, mappedDoc.Document: %v, printDocSeparators: %v", p.firstTimePrinting, p.previousDocIndex, mappedDoc.Document, p.printDocSeparators)
		if (p.previousDocIndex != mappedDoc.Document || p.previousFileIndex != mappedDoc.FileIndex) && p.printDocSeparators {
			log.Debug("-- writing doc sep")
			if err := p.writeString(bufferedWriter, "---\n"); err != nil {
				return err
			}
		}

		if err := p.printNode(mappedDoc.Node, bufferedWriter); err != nil {
			return err
		}

		p.previousDocIndex = mappedDoc.Document
	}

	return nil
}
