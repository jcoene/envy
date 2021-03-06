package envy

var c *Client

// Create a new default client with a given name and set of machines.
func NewDefaultClient(name string, machines []string) {
	c = NewClient(name, machines)
}

func checkDefaultClient() {
	if c == nil {
		panic("no default client set")
	}
}

// Set default key/value pairs using the default client.
func SetDefaults(kv map[string]string) {
	checkDefaultClient()
	c.SetDefaults(kv)
}

// Load all keys using the default client.
func LoadAll() error {
	checkDefaultClient()
	return c.LoadAll()
}

// Get a key using the default client.
func Get(k string) string {
	return c.Get(k)
}

// Get a key using the default client.
func GetSafe(k string) string {
	return c.GetSafe(k)
}

// Get all keys using the default client
func GetAll() map[string]string {
	return c.GetAll()
}

// Create a new lock using the default client
func NewLock(name string, ttl int) *Lock {
	checkDefaultClient()
	return c.NewLock(name, ttl)
}

// Create a new register using the default client
func NewRegister(name string) *Register {
	checkDefaultClient()
	return c.NewRegister(name)
}
