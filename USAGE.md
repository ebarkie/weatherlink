# weatherlink

```go
import "github.com/ebarkie/weatherlink"
```

Package weatherlink implements the Davis Instruments serial, USB, and TCP/IP
communication protocol.

## Usage

```go
const (
	GetDmps cmd = iota
	GetEEPROM
	GetHiLows
	GetLoops
	LampsOff
	LampsOn
	Stop
	SyncConsTime
)
```
Commands.

```go
var (
	Trace = log.New(ioutil.Discard, "[TRCE]", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	Debug = log.New(ioutil.Discard, "[DBUG]", log.LstdFlags|log.Lshortfile)
	Info  = log.New(ioutil.Discard, "[INFO]", log.LstdFlags)
	Warn  = log.New(ioutil.Discard, "[WARN]", log.LstdFlags|log.Lshortfile)
	Error = log.New(ioutil.Discard, "[ERRO]", log.LstdFlags|log.Lshortfile)
)
```
Loggers.

```go
var (
	ConsTimeSyncFreq = 24 * time.Hour
)
```
Tunables.

```go
var (
	ErrCmdFailed = errors.New("command failed")
)
```
Errors.

```go
var Sdump = func(i ...interface{}) (s string) {
	return fmt.Sprintf(strings.Repeat("%+v\n", len(i)), i...)
}
```
Sdump returns a variable as a string. It includes field names and pointers, if
applicable.

#### func  StdIdle

```go
func StdIdle(c *Conn, ec chan<- interface{}) (err error)
```
StdIdle is the standard idler which reads loop packets and new archive records
when they're available.

#### type Conn

```go
type Conn struct {
	LastDmp   time.Time // Time of the last downloaded archive record
	NewArcRec bool      // Indicates a new archive record is available

	Q chan cmd // Command queue
}
```

Conn holds the weatherlink connnection context.

#### func  Dial

```go
func Dial(addr string) (c Conn, err error)
```
Dial establishes the weatherlink connection.

#### func (Conn) Close

```go
func (c Conn) Close() error
```
Close closes the weatherlink connection.

#### func (Conn) GetConsTime

```go
func (c Conn) GetConsTime() (t time.Time, err error)
```
GetConsTime gets the console time.

#### func (Conn) GetDmps

```go
func (c Conn) GetDmps(ec chan<- interface{}, lastRec time.Time) (newLastRec time.Time, err error)
```
GetDmps downloads all archive records *after* lastRec and sends them to the
event channel ordered from oldest to newest. It returns the time of the last
record it read.

If lastRec does not match an existing archive timestamp (which is the case if
left uninitialized) then all records in memory are returned.

#### func (Conn) GetEEPROM

```go
func (c Conn) GetEEPROM(ec chan<- interface{}) error
```
GetEEPROM retrieves the entire EEPROM configuration.

#### func (Conn) GetHiLows

```go
func (c Conn) GetHiLows(ec chan<- interface{}) error
```
GetHiLows retrieves the record high and lows.

#### func (*Conn) GetLoops

```go
func (c *Conn) GetLoops(ec chan<- interface{}) (err error)
```
GetLoops starts a stream of loop packets and sends them to the event channel. It
exits when numLoops is hit, an archive record was written, or a command is
pending.

#### func (Conn) SetLamps

```go
func (c Conn) SetLamps(on bool) (err error)
```
SetLamps sets the console lamps state.

#### func (*Conn) Start

```go
func (c *Conn) Start(idle Idler) <-chan interface{}
```
Start starts the command broker. If no commands are pending it runs the idler.

#### func (Conn) Stop

```go
func (c Conn) Stop()
```
Stop stops the command broker.

#### func (Conn) SyncConsTime

```go
func (c Conn) SyncConsTime() (err error)
```
SyncConsTime synchronizes the console time with the local system time if the
offset exceeds 10 seconds.

#### type Idler

```go
type Idler func(*Conn, chan<- interface{}) error
```

Idler is the idle function the command broker executes when there are no pending
commands in the queue.
