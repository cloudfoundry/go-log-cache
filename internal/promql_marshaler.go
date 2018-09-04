package internal

import (
	"encoding/json"
	"fmt"
	"io"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type PromqlMarshaler struct {
	fallback runtime.Marshaler
}

func NewPromqlMarshaler(fallback runtime.Marshaler) *PromqlMarshaler {
	return &PromqlMarshaler{
		fallback: fallback,
	}
}

func (m *PromqlMarshaler) Marshal(v interface{}) ([]byte, error) {
	switch q := v.(type) {
	case *logcache_v1.PromQL_InstantQueryResult:
		return json.Marshal(m.assembleInstantQueryResult(q))
	case *logcache_v1.PromQL_RangeQueryResult:
		return json.Marshal(m.assembleRangeQueryResult(q))
	default:
		return m.fallback.Marshal(v)
	}
}

type queryResult struct {
	Status    string     `json:"status"`
	Data      resultData `json:"data"`
	ErrorType string     `json:"errorType"`
	Error     string     `json:"error"`
}

type queryResultUnmarshal struct {
	Status    string              `json:"status"`
	Data      resultDataUnmarshal `json:"data"`
	ErrorType string              `json:"errorType"`
	Error     string              `json:"error"`
}

type resultData struct {
	ResultType string        `json:"resultType"`
	Result     []interface{} `json:"result,omitempty"`
}

type resultDataUnmarshal struct {
	ResultType string          `json:"resultType"`
	Result     json.RawMessage `json:"result,omitempty"`
}

type sample struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type series struct {
	Metric map[string]string `json:"metric"`
	Values [][]interface{}   `json:"values"`
}

func (m *PromqlMarshaler) assembleInstantQueryResult(v *logcache_v1.PromQL_InstantQueryResult) *queryResult {
	var data resultData
	switch v.GetResult().(type) {
	case *logcache_v1.PromQL_InstantQueryResult_Scalar:
		data = assembleScalarResultData(v.GetScalar())
	case *logcache_v1.PromQL_InstantQueryResult_Vector:
		data = assembleVectorResultData(v.GetVector())
	case *logcache_v1.PromQL_InstantQueryResult_Matrix:
		data = assembleMatrixResultData(v.GetMatrix())
	}

	return &queryResult{
		Status: "success",
		Data:   data,
	}
}

func (m *PromqlMarshaler) assembleRangeQueryResult(v *logcache_v1.PromQL_RangeQueryResult) *queryResult {
	var data resultData
	switch v.GetResult().(type) {
	case *logcache_v1.PromQL_RangeQueryResult_Matrix:
		data = assembleMatrixResultData(v.GetMatrix())
	}

	return &queryResult{
		Status: "success",
		Data:   data,
	}
}

func assembleScalarResultData(v *logcache_v1.PromQL_Scalar) resultData {
	return resultData{
		ResultType: "scalar",
		Result:     []interface{}{v.GetTime(), v.GetValue()},
	}
}

func assembleVectorResultData(v *logcache_v1.PromQL_Vector) resultData {
	var samples []interface{}

	for _, s := range v.GetSamples() {
		p := s.GetPoint()
		samples = append(samples, sample{
			Metric: s.GetMetric(),
			Value:  []interface{}{p.GetTime(), p.GetValue()},
		})
	}

	return resultData{
		ResultType: "vector",
		Result:     samples,
	}
}

func assembleMatrixResultData(v *logcache_v1.PromQL_Matrix) resultData {
	var result []interface{}

	for _, s := range v.GetSeries() {
		var values [][]interface{}
		for _, p := range s.GetPoints() {
			values = append(values, []interface{}{p.GetTime(), p.GetValue()})
		}
		result = append(result, series{
			Metric: s.GetMetric(),
			Values: values,
		})
	}

	return resultData{
		ResultType: "matrix",
		Result:     result,
	}
}

func (m *PromqlMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	fallbackEncoder := m.fallback.NewEncoder(w)
	jsonEncoder := json.NewEncoder(w)

	return runtime.EncoderFunc(func(v interface{}) error {
		switch q := v.(type) {
		case *logcache_v1.PromQL_InstantQueryResult:
			return jsonEncoder.Encode(m.assembleInstantQueryResult(q))
		case *logcache_v1.PromQL_RangeQueryResult:
			return jsonEncoder.Encode(m.assembleRangeQueryResult(q))
		default:
			return fallbackEncoder.Encode(v)
		}
	})
}

// The special marshaling for PromQL results is currently only implemented
// for encoding.
func (m *PromqlMarshaler) Unmarshal(data []byte, v interface{}) error {
	var result queryResultUnmarshal
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	switch q := v.(type) {
	case *logcache_v1.PromQL_InstantQueryResult:
		r, err := m.disassembleInstantQueryResult(result)
		if err != nil {
			return err
		}
		*q = *r
	case *logcache_v1.PromQL_RangeQueryResult:
		r, err := m.disassembleRangeQueryResult(result)
		if err != nil {
			return err
		}
		*q = *r
	default:
		return nil
	}

	return nil
	// return m.fallback.Unmarshal(data, v)
}

func (m *PromqlMarshaler) disassembleInstantQueryResult(q queryResultUnmarshal) (*logcache_v1.PromQL_InstantQueryResult, error) {
	switch q.Data.ResultType {
	case "scalar":
		r, err := unmarshalScalarResultData(q.Data.Result)
		if err != nil {
			return nil, err
		}

		return &logcache_v1.PromQL_InstantQueryResult{
			Result: &logcache_v1.PromQL_InstantQueryResult_Scalar{
				Scalar: r,
			},
		}, nil
	case "vector":
		r, err := unmarshalVectorResultData(q.Data.Result)
		if err != nil {
			return nil, err
		}

		return &logcache_v1.PromQL_InstantQueryResult{
			Result: &logcache_v1.PromQL_InstantQueryResult_Vector{
				Vector: r,
			},
		}, nil
	case "matrix":
		r, err := unmarshalMatrixResultData(q.Data.Result)
		if err != nil {
			return nil, err
		}

		return &logcache_v1.PromQL_InstantQueryResult{
			Result: &logcache_v1.PromQL_InstantQueryResult_Matrix{
				Matrix: r,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown instant query resultType '%s'", q.Data.ResultType)
	}
}

func (m *PromqlMarshaler) disassembleRangeQueryResult(q queryResultUnmarshal) (*logcache_v1.PromQL_RangeQueryResult, error) {
	switch q.Data.ResultType {
	case "matrix":
		r, err := unmarshalMatrixResultData(q.Data.Result)
		if err != nil {
			return nil, err
		}

		return &logcache_v1.PromQL_RangeQueryResult{
			Result: &logcache_v1.PromQL_RangeQueryResult_Matrix{
				Matrix: r,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown range query resultType '%s'", q.Data.ResultType)
	}
}

func unmarshalScalarResultData(data []byte) (*logcache_v1.PromQL_Scalar, error) {
	var value []float64
	err := json.Unmarshal(data, &value)
	if err != nil {
		return nil, err
	}

	return &logcache_v1.PromQL_Scalar{
		Time:  int64(value[0]),
		Value: value[1],
	}, nil
}

func unmarshalVectorResultData(data []byte) (*logcache_v1.PromQL_Vector, error) {
	var values []sample
	err := json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}

	var samples []*logcache_v1.PromQL_Sample
	for _, value := range values {
		samples = append(samples, &logcache_v1.PromQL_Sample{
			Metric: value.Metric,
			Point: &logcache_v1.PromQL_Point{
				Time:  int64(value.Value[0].(float64)),
				Value: value.Value[1].(float64),
			},
		})
	}
	return &logcache_v1.PromQL_Vector{
		Samples: samples,
	}, nil
}

func unmarshalMatrixResultData(data []byte) (*logcache_v1.PromQL_Matrix, error) {
	var values []series
	err := json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}

	var serieses []*logcache_v1.PromQL_Series
	for _, value := range values {
		var points []*logcache_v1.PromQL_Point
		for _, point := range value.Values {
			points = append(points, &logcache_v1.PromQL_Point{
				Time:  int64(point[0].(float64)),
				Value: point[1].(float64),
			})
		}
		serieses = append(serieses, &logcache_v1.PromQL_Series{
			Metric: value.Metric,
			Points: points,
		})
	}

	return &logcache_v1.PromQL_Matrix{
		Series: serieses,
	}, nil
}

func (m *PromqlMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return m.fallback.NewDecoder(r)
}

func (m *PromqlMarshaler) ContentType() string {
	return `application/json`
}
