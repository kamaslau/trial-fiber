package drivers

// Initiate Driver Components first, fail fast so we can break early
func InitDrivers() {
	ConnectDB()    // Exit on failure
	ConnectCache() // Warning on failure
	ConnectTS()    // Warning on failure
	ConnectMQ()    // Warning on failure
}
