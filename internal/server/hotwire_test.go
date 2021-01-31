package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

const message = `<turbo-stream action="replace" target="load">
    <template>
        <span id="load">04:20:13: 1.9</span>
    </template>
</turbo-stream>`

const example = `event: message
id: 6
data: <turbo-stream action="replace" target="load">
data:     <template>
data:         <span id="load">04:20:13: 1.9</span>
data:     </template>
data: </turbo-stream>

`

func Test_writeMessageWithoutNewline(t *testing.T) {
	assert := require.New(t)
	buf := new(bytes.Buffer)

	err := writeMessage(buf, 6, "message", message)
	assert.NoError(err)
	assert.Equal(example, buf.String())
}
