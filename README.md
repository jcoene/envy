# Envy

Envy provides application configuration using Etcd. It also supports simple defaults so that you don't need etcd in your development environment.

It's an easy to use drop-in replacement for shipping configuration files or environment variables during deployment.

Pull requests welcome.

## Namespaces

An instance of the Envy client is initialized with a namespace, which is used to determine the path for configuration keys.

Given the namespace "myapp", the path for the key "foo" is:

*/myapp/config/foo*

## License

Envy is released under the MIT License. See the included LICENSE file for details.
