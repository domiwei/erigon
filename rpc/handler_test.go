// Copyright 2024 The Erigon Authors
// This file is part of Erigon.
//
// Erigon is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Erigon is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Erigon. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"reflect"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	"github.com/erigontech/erigon/common/log/v3"
	"github.com/erigontech/erigon/rpc/jsonstream"
)

func TestHandlerDoesNotDoubleWriteNull(t *testing.T) {

	tests := map[string]struct {
		params   []byte
		expected string
	}{
		"error_with_stream_write": {
			params:   []byte("[1]"),
			expected: `{"jsonrpc":"2.0","id":1,"result":null,"error":{"code":-32000,"message":"id 1"}}`,
		},
		"error_without_stream_write": {
			params:   []byte("[2]"),
			expected: `{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"id 2"}}`,
		},
		"no_error": {
			params:   []byte("[3]"),
			expected: `{"jsonrpc":"2.0","id":1,"result":{}}`,
		},
		"err_with_valid_json_empty_array": {
			params:   []byte("[4]"),
			expected: `{"jsonrpc":"2.0","id":1,"result":{"structLogs":[]},"error":{"code":-32000,"message":"id 4"}}`,
		},
		"err_with_valid_json_empty_object": {
			params:   []byte("[5]"),
			expected: `{"jsonrpc":"2.0","id":1,"result":{"structLogs":{}},"error":{"code":-32000,"message":"id 4"}}`,
		},
		"err_with_unclosed_result_object": {
			params:   []byte("[6]"),
			expected: `{"jsonrpc":"2.0","id":1,"result":{"structLogs":[]},"error":{"code":-32000,"message":"id 6"}}`,
		},
	}

	for name, testParams := range tests {
		t.Run(name, func(t *testing.T) {
			msg := jsonrpcMessage{
				Version: "2.0",
				ID:      []byte{49},
				Method:  "test_test",
				Params:  testParams.params,
				Error:   nil,
				Result:  nil,
			}

			dummyFunc := func(id int, stream jsonstream.Stream) error {
				if id == 1 {
					stream.WriteNil()
					return errors.New("id 1")
				}
				if id == 2 {
					return errors.New("id 2")
				}
				if id == 3 {
					stream.WriteEmptyObject()
					return nil
				}
				if id == 4 {
					stream.WriteObjectStart()
					stream.WriteObjectField("structLogs")
					stream.WriteEmptyArray()
					stream.WriteObjectEnd()
					return errors.New("id 4")
				}
				if id == 5 {
					stream.WriteObjectStart()
					stream.WriteObjectField("structLogs")
					stream.WriteEmptyObject()
					stream.WriteObjectEnd()
					return errors.New("id 4")
				}
				if id == 6 {
					stream.WriteObjectStart()
					stream.WriteObjectField("structLogs")
					stream.WriteEmptyArray()
					// intentionally leave the result object open: the tracer erroring out
					// mid-write must not leave the response's "result" object unclosed.
					return errors.New("id 6")
				}
				return nil
			}

			cb := &callback{
				fn:          reflect.ValueOf(dummyFunc),
				rcvr:        reflect.Value{},
				argTypes:    []reflect.Type{reflect.TypeFor[int]()},
				hasCtx:      false,
				errPos:      0,
				isSubscribe: false,
				streamable:  true,
			}

			args, err := parsePositionalArguments((msg).Params, cb.argTypes)
			if err != nil {
				t.Fatal(err)
			}

			var buf bytes.Buffer
			stream := jsonstream.New(jsoniter.NewStream(jsoniter.ConfigDefault, &buf, 4096))

			h := handler{}
			h.runMethod(context.Background(), &msg, cb, args, stream)

			stream.Flush()

			output := buf.String()
			assert.Equal(t, testParams.expected, output, "expected output should match")
		})
	}

}

func TestHandleBatchMixedWithResponseMessage(t *testing.T) {
	logger := log.New()
	server := newTestServer(logger)
	defer server.Stop()

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	go server.ServeCodec(NewCodec(serverConn), 0)

	batch := `[{"jsonrpc":"2.0","id":1,"method":"test_echo","params":["hello",10,{}]},` +
		`{"jsonrpc":"2.0","id":99,"result":"stale"}]` + "\n"

	clientConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if _, err := io.WriteString(clientConn, batch); err != nil {
		t.Fatalf("write: %v", err)
	}

	clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	reader := bufio.NewReader(clientConn)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("read: %v (deadlock if timeout)", err)
	}

	var msgs []json.RawMessage
	if err := json.Unmarshal([]byte(line), &msgs); err != nil {
		t.Fatalf("unmarshal batch response: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 response in batch, got %d", len(msgs))
	}
}

// TestRunMethodFlushHookNilFunc pins the invariant that runMethod must not panic when the
// gzip-streaming hook stored on the context is a typed nil func(), not just an untyped nil.
// The normal masking path (withoutGzipStreamingHook) stores an untyped nil so the type
// assertion fails outright, but runMethod's guard should not depend on callers always doing
// that correctly.
func TestRunMethodFlushHookNilFunc(t *testing.T) {
	msg := jsonrpcMessage{
		Version: "2.0",
		ID:      []byte{49},
		Method:  "test_test",
		Params:  []byte("[]"),
	}

	dummyFunc := func(stream jsonstream.Stream) error {
		stream.WriteEmptyObject()
		return nil
	}

	cb := &callback{
		fn:         reflect.ValueOf(dummyFunc),
		streamable: true,
	}

	args, err := parsePositionalArguments(msg.Params, cb.argTypes)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), httpFlusherContextKey{}, (func())(nil))

	var buf bytes.Buffer
	stream := jsonstream.New(jsoniter.NewStream(jsoniter.ConfigDefault, &buf, 4096))

	h := handler{}
	assert.NotPanics(t, func() {
		h.runMethod(ctx, &msg, cb, args, stream)
	})
}
