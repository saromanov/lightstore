# Client for lightstore

## API

### NewClient(addr string)
Initialization of client

### Set(key, value string)

### Get(key)

### SetMap(map[string]string)
Set several key-value pairs

### Stat()
Return statistics

### CreatePage(pagename string)
Create new page

### SetToPage(pagename, key, value string)
Write key-value data to specific pagename

### GetFromPage(pagename, key string)
Get key-value data from specific pagename
