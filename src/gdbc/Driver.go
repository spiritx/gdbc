package gdbc

type Driver interface {
	/*
	Attempts to make a database connection to the given URL.
	The driver should return "false" if it realizes it is the wrong kind of driver to connect to the given URL.
	This will be common, as when the GDBC driver manager is asked to connect to a given URL
	it passes the URL to each loaded driver in turn.

	The properties argument can be used to pass arbitrary string tag/value pairs as connection arguments.
	Normally at least "user" and "password" properties should be included in the Properties object.

	Note: If a property is specified as part of the url and is also specified in the Properties object,
	it is implementation-defined as to which value will take precedence.
	For maximum portability, an application should only specify a property once.

	Parameters:
	url - the URL of the database to which to connect
	info - a list of arbitrary string tag/value pairs as connection arguments.
	Normally at least a "user" and "password" property should be included.
	Returns:
	a Connection object that represents a connection to the URL
	 */
	Connect(url string, info map[string]string) (connection Connection, ok bool)
}




