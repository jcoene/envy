package envy

import (
	"path"
	"strconv"
)

type Register struct {
	name string
	key  string
	c    *Client
}

func (c *Client) NewRegister(name string) *Register {
	return &Register{
		name: name,
		key:  path.Join(c.ns, "register", name),
		c:    c,
	}
}

func (r *Register) String() (string, error) {
	resp, err := r.c.c.Get(r.key, false, false)
	if err != nil {
		return "", err
	}

	return resp.Node.Value, nil
}

func (r *Register) Int64() (int64, error) {
	s, err := r.String()
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (r *Register) Set(v string) error {
	_, err := r.c.c.Set(r.key, v, 0)
	return err
}

func (r *Register) SetInt64(i int64) error {
	return r.Set(strconv.FormatInt(i, 10))
}
