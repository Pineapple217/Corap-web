// Code generated by templ@v0.2.408 DO NOT EDIT.

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"reflect"
	"time"
)

func DevicePlot(plotType string, datas []float32, timestamps []time.Time) templ.Component {
	return templ.ComponentFunc(func(templ_7745c5c3_Ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		templ_7745c5c3_Ctx = templ.InitializeContext(templ_7745c5c3_Ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(templ_7745c5c3_Ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		templ_7745c5c3_Ctx = templ.ClearChildren(templ_7745c5c3_Ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString("line-" + plotType))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" style=\"width: 800px; height: 450px;\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = OnLoad(lineGraph(plotType, datas, timestamps), plotType, datas, timestamps).Render(templ_7745c5c3_Ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func lineGraph(plotType string, datas []float32, timestamps []time.Time) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_lineGraph_6c02`,
		Function: `function __templ_lineGraph_6c02(plotType, datas, timestamps){var trace1 = {
        x: timestamps,
        y: datas,
        type: "scatter",
    };
    Plotly.newPlot("line-" + plotType, [trace1]);}`,
		Call: templ.SafeScript(`__templ_lineGraph_6c02`, plotType, datas, timestamps),
	}
}

func DeviceHeatmap(plotType string, datas []float32, timestamps []time.Time) templ.Component {
	return templ.ComponentFunc(func(templ_7745c5c3_Ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		templ_7745c5c3_Ctx = templ.InitializeContext(templ_7745c5c3_Ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(templ_7745c5c3_Ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		templ_7745c5c3_Ctx = templ.ClearChildren(templ_7745c5c3_Ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString("heatmap-" + plotType))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" style=\"width: 800px; height: 450px;\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if d := ChunkifyByWeek(datas, int(timestamps[0].Weekday())); true {
			templ_7745c5c3_Err = OnLoad(graph(plotType, d, timestamps), plotType, d, timestamps).Render(templ_7745c5c3_Ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func graph(plotType string, datas interface{}, timestamps []time.Time) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_graph_97a8`,
		Function: `function __templ_graph_97a8(plotType, datas, timestamps){let o = datas
    var trace1 = [{
        z: o,
        x: ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'],
		y: timestamps.reduce((acc, value, index) => (index % 7 === 6) ? [...acc, value] : acc, []),
		hovertemplate: '%{z:.2f}',
        type: "heatmap",
		colorscale: 'Portland',
    }];
	var layout = {
		yaxis: { autorange: "reversed" }
	}
    Plotly.newPlot("heatmap-" + plotType, trace1, layout);}`,
		Call: templ.SafeScript(`__templ_graph_97a8`, plotType, datas, timestamps),
	}
}

func ChunkifyByWeek(input []float32, startDay int) (result interface{}) {
	chunkSize := 7
	nilSlice := make([]float32, chunkSize-startDay)
	newSlice := append(nilSlice, input...)

	inputValue := reflect.ValueOf(newSlice)

	if inputValue.Kind() != reflect.Slice {
		panic("Input must be a slice")
	}

	length := inputValue.Len()

	// Calculate the number of chunks needed
	numChunks := ((length + chunkSize - 1) / chunkSize) - 1

	// Create a slice to hold the result
	resultSlice := reflect.MakeSlice(reflect.SliceOf(inputValue.Type()), numChunks, numChunks)

	for i := 0; i < numChunks; i++ {
		// Calculate start and end indices for each chunk
		start := i * chunkSize
		end := (i + 1) * chunkSize

		if i == 0 {
			start += 7 - startDay
		}
		if end > length {
			end = length
		}

		subarray := inputValue.Slice(start, end)
		reversed := reverseImmutableSlice(subarray)
		resultSlice.Index(i).Set(reversed)
	}

	return resultSlice.Interface()
}

func reverseImmutableSlice(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Slice {
		fmt.Println("Not a slice")
		return reflect.Value{}
	}

	length := v.Len()
	reversed := reflect.MakeSlice(v.Type(), length, length)

	for i := 0; i < length; i++ {
		reversed.Index(i).Set(v.Index(length - i - 1))
	}

	return reversed
}
