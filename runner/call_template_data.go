package runner

import (
	"bytes"
	"encoding/json"
	"strconv"
	"text/template"
	"time"

	"github.com/jhump/protoreflect/desc"
)

// call template data
type callTemplateData struct {
	RequestNumber      int64  // unique incrememnted request number for each request
	RequestNumberInBytes string // unique incrememnted request number for each request in unicode bytes
	FullyQualifiedName string // fully-qualified name of the method call
	MethodName         string // shorter call method name
	ServiceName        string // the service name
	InputName          string // name of the input message type
	OutputName         string // name of the output message type
	IsClientStreaming  bool   // whether this call is client streaming
	IsServerStreaming  bool   // whether this call is server streaming
	Timestamp          string // timestamp of the call in RFC3339 format
	TimestampUnix      int64  // timestamp of the call as unix time
}

// newCallTemplateData returns new call template data
func newCallTemplateData(mtd *desc.MethodDescriptor, reqNum int64) *callTemplateData {
	now := time.Now()

	return &callTemplateData{
		RequestNumber:      reqNum,
		RequestNumberInBytes: string([]byte(strconv.FormatInt(reqNum, 10))),
		FullyQualifiedName: mtd.GetFullyQualifiedName(),
		MethodName:         mtd.GetName(),
		ServiceName:        mtd.GetService().GetName(),
		InputName:          mtd.GetInputType().GetName(),
		OutputName:         mtd.GetOutputType().GetName(),
		IsClientStreaming:  mtd.IsClientStreaming(),
		IsServerStreaming:  mtd.IsServerStreaming(),
		Timestamp:          now.Format(time.RFC3339),
		TimestampUnix:      now.Unix(),
	}
}

func (td *callTemplateData) execute(data string) (*bytes.Buffer, error) {
	t := template.Must(template.New("call_template_data").Parse(data))
	var tpl bytes.Buffer
	err := t.Execute(&tpl, td)
	return &tpl, err
}

func (td *callTemplateData) executeData(data string) ([]byte, error) {
	if len(data) > 0 {
		input := []byte(data)
		tpl, err := td.execute(data)
		if err == nil {
			input = tpl.Bytes()
		}

		return input, nil
	}

	return []byte{}, nil
}

func (td *callTemplateData) executeMetadata(metadata string) (*map[string]string, error) {
	var mdMap map[string]string

	if len(metadata) > 0 {
		input := []byte(metadata)
		tpl, err := td.execute(metadata)
		if err == nil {
			input = tpl.Bytes()
		}

		err = json.Unmarshal(input, &mdMap)
		if err != nil {
			return nil, err
		}
	}

	return &mdMap, nil
}
