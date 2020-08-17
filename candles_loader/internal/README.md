# Internal sources

* config - Structs and reader for configuration
  * configuration - Configuration structure
  * environment - Functions for reading environment variables into configuration
  * reader - Functions for reading configuration from file
* date - Tools for working with dates types
  * generator - Generator of dates sequencies
  * generator_test - Unit-tests for generator
  * tools - Different tools for working with date types
  * tools_test - Unit-tests for date's tools
* db - Data base connector
  * checker - Functions for check info in data base
  * configuration - Struct for data base connection configuration
  * deleter - Functions for deleting information from data base
  * uploader - Functions for upload information into data base
* globalrank - Connector for getting global rank from Forbes
  * data - Struct for reading global Forbes companies rank
  * reader - Functions for reading global Forbes info from files
* tinkoff - Connector to Tinkoff Investing API
  * markets - Functions for getting information about markets from Tinkoff API
