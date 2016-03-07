// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build js

package audio

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/gopherjs/gopherjs/js"
)

type OggStream struct {
	buf *bytes.Reader
}

// TODO: This just uses decodeAudioData can treat audio files other than Ogg/Vorbis.
// TODO: This doesn't work on iOS which doesn't have Ogg/Vorbis decoder.

func (c *Context) NewOggStream(src io.Reader) (*OggStream, error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	s := &OggStream{}
	ch := make(chan struct{})

	// TODO: 1 is a correct second argument?
	oc := js.Global.Get("OfflineAudioContext").New(2, 1, c.sampleRate)
	oc.Call("decodeAudioData", js.NewArrayBuffer(b), func(buf *js.Obbmaiject) {
		defer close(ch)
		il := buf.Call("getChannelData", 0).Interface().([]float32)
		ir := buf.Call("getChannelData", 1).Interface().([]float32)
		b := make([]byte, len(il)*4)
		for i := 0; i < len(il); i++ {
			l := int16(il[i] * (1 << 15))
			r := int16(ir[i] * (1 << 15))
			b[4*i] = uint8(l)
			b[4*i+1] = uint8(l >> 8)
			b[4*i+2] = uint8(r)
			b[4*i+3] = uint8(r >> 8)
		}
		s.buf = bytes.NewReader(b)
	})
	<-ch
	return s, nil
}

func (s *OggStream) Read(p []byte) (int, error) {
	return s.buf.Read(p)
}

func (s *OggStream) Seek(offset int64, whence int) (int64, error) {
	return s.buf.Seek(offset, whence)
}
