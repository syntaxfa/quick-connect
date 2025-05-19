package traceotela

import "google.golang.org/grpc/metadata"

type metadataCarrier struct {
	metadata.MD
}

func (c metadataCarrier) Get(key string) string {
	vs := c.MD.Get(key)
	if len(vs) == 0 {
		return ""
	}

	return vs[0]
}

func (c metadataCarrier) Set(key, value string) {
	c.MD.Set(key, value)
}

func (c metadataCarrier) Keys() []string {
	out := make([]string, 0, len(c.MD))

	for k := range c.MD {
		out = append(out, k)
	}

	return out
}
