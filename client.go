package envy

import (
	"fmt"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

type Client struct {
	ns string
	c  *etcd.Client
	v  map[string]string
}

// Create a new client with a given namespace and array of hosts
// If no machine set is given we use the default of 127.0.0.1:4001
func NewClient(ns string, machines []string) *Client {
	return &Client{
		ns: ns,
		c:  etcd.NewClient(machines),
		v:  make(map[string]string),
	}
}

// Set default key/value pairs to be used in case a key cannot be found in etcd.
// These will not override existing values set from another source.
func (c *Client) SetDefaults(kv map[string]string) (err error) {
	for k, v := range kv {
		if _, f := c.v[k]; !f {
			c.v[k] = v
		}
	}

	return
}

// Load all key-value pairs in our namespace from etcd.
func (c *Client) LoadAll() (err error) {
	var resp *etcd.Response

	if resp, err = c.c.Get(c.rootNode(), true, true); err != nil {
		return
	}

	err = c.parseNodes(resp.Node.Nodes)

	return
}

// Get a key
func (c *Client) Get(k string) string {
	return c.v[k]
}

// Get all keys
func (c *Client) GetAll() map[string]string {
	return c.v
}

func (c *Client) rootNode() string {
	return fmt.Sprintf("/%s/config", c.ns)
}

func (c *Client) parseNodes(ns etcd.Nodes) (err error) {
	for _, n := range ns {
		if n.Dir {
			if err = c.parseNodes(n.Nodes); err != nil {
				return
			}
		} else {
			c.v[c.extractKey(n.Key)] = n.Value
		}
	}

	return
}

func (c *Client) extractKey(path string) (key string) {
	return strings.Replace(path, c.rootNode()+"/", "", 1)
}
