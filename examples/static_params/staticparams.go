package main

import (
	// "fmt"
	sci "github.com/scipipe/scipipe"
)

func main() {
	// Init
	wf := sci.NewWorkflow("staticparams_wf", 4)

	fls := NewFileSender(wf, "filesender")

	abc := sci.NewProc(wf, "abc", "echo {p:a} {p:b} {p:c} > {o:out} # {i:in}")
	abc.ParamInPort("a").ConnectStr("a1")
	abc.ParamInPort("b").ConnectStr("b1")
	abc.ParamInPort("c").ConnectStr("c1")
	abc.SetPathCustom("out", func(task *sci.Task) string {
		return task.InPath("in")
	})

	prt := wf.NewProc("prt", "echo {i:in} >> log.txt")

	// Connect
	abc.In("in").Connect(fls.Out())
	prt.In("in").Connect(abc.Out("out"))

	// Run
	wf.Run()
}

// --------------------------------
//         A Custom task
// --------------------------------

type FileSender struct {
	sci.BaseProcess
}

func NewFileSender(wf *sci.Workflow, name string) *FileSender {
	p := &FileSender{
		BaseProcess: sci.NewBaseProcess(wf, name),
	}
	p.InitOutPort(p, "out")
	wf.AddProc(p)
	return p
}

func (p *FileSender) Out() *sci.OutPort { return p.OutPort("out") }

func (p *FileSender) Run() {
	defer p.Out().Close()
	for _, fn := range []string{"file1.txt", "file2.txt", "file3.txt"} {
		p.Out().Send(sci.NewFileIP(fn))
	}
}
