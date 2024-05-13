package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

type ResponseWriter struct {
	gin.ResponseWriter
	StatusCode    int
	ContextLength int
}

func NewResponseWriter(w gin.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, w.Status(), w.Size()}
}

func (lrw *ResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ContextLength = lrw.Size()
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *ResponseWriter) LogQueryParams(c *gin.Context, startTime time.Time) {
	uri := c.Request.RequestURI
	method := c.Request.Method
	log.Info().
		Str("uri", uri).
		Str("method", method).
		Int("status", lrw.StatusCode).
		Int("contextLength", lrw.ContextLength).
		Dur("execTime", time.Since(startTime)).
		Msg("logQueryParams(): stats")
}

func (lrw *ResponseWriter) SendEncodedBody(code int, bytesResponse []byte) {
	lrw.Header().Add("Content-Encoding", "gzip")
	lrw.Header().Add("Content-Type", "application/json")
	lrw.WriteHeader(code)
	_, err := lrw.Write(bytesResponse)
	if err != nil {
		log.Error().
			Err(err).
			Int("code", code).
			Msg("SendEncodedBody(): write error")
	}
}
