package mpv

import (
	"errors"
	"fmt"
	"strconv"
)

// Client is a more comfortable higher level interface
// to LLClient. It can use any LLClient implementation.
type Client struct {
	LLClient
}

// NewClient creates a new highlevel client based on a lowlevel client.
func NewClient(llclient LLClient) *Client {
	return &Client{
		llclient,
	}
}

// Mode options for Loadfile
const (
	LoadFileModeReplace    = "replace"
	LoadFileModeAppend     = "append"
	LoadFileModeAppendPlay = "append-play" // Starts if nothing is playing
)

// Loadfile loads a file, it either replaces the currently playing file `LOAD_REPLACE`,
// appends to the current playlist `LOAD_APPEND` or appends to playlist and plays if
// nothing is playing right now `LOAD_APPEND_PLAY`
func (c *Client) Loadfile(path string, mode string) error {
	_, err := c.Exec("loadfile", path, mode)
	return err
}

// Mode options for Seek
const (
	SeekModeRelative = "relative"
	SeekModeAbsolute = "absolute"
)

// Seek seeks to a position in the current file.
func (c *Client) Seek(n int, mode string) error {
	_, err := c.Exec("seek", strconv.Itoa(n), mode)
	return err
}

// PlaylistNext plays the next playlistitem or NOP if no item is available.
func (c *Client) PlaylistNext() error {
	_, err := c.Exec("playlist-next", "weak")
	return err
}

// PlaylistPrevious plays the previous playlistitem or NOP if no item is available.
func (c *Client) PlaylistPrevious() error {
	_, err := c.Exec("playlist-prev", "weak")
	return err
}

// Mode options for LoadList
const (
	LoadListModeReplace = "replace"
	LoadListModeAppend  = "append"
)

// LoadList loads a playlist from path. It can either replace the current playlist `LOADLIST_REPLACE`
// or append to the current playlist `LOADLIST_APPEND`.
func (c *Client) LoadList(path string, mode string) error {
	_, err := c.Exec("loadlist", path, mode)
	return err
}

// GetProperty reads a property by name and returns the data as a string.
func (c *Client) GetProperty(name string) (string, error) {
	res, err := c.Exec("get_property", name)
	if res == nil {
		return "", err
	}
	return fmt.Sprintf("%#v", res.Data), err
}

// SetProperty sets the value of a property.
func (c *Client) SetProperty(name string, value interface{}) error {
	_, err := c.Exec("set_property", name, value)
	return err
}

// ErrInvalidType is returned if the response data does not match the methods return type.
// Use GetProperty or find matching type in mpv docs.
var ErrInvalidType = errors.New("Invalid type")

// GetFloatProperty reads a float property and returns the data as a float64.
func (c *Client) GetFloatProperty(name string) (float64, error) {
	res, err := c.Exec("get_property", name)
	if res == nil {
		return 0, err
	}
	if val, found := res.Data.(float64); found {
		return val, err
	}
	return 0, ErrInvalidType
}

// GetBoolProperty reads a bool property and returns the data as a boolean.
func (c *Client) GetBoolProperty(name string) (bool, error) {
	res, err := c.Exec("get_property", name)
	if res == nil {
		return false, err
	}
	if val, found := res.Data.(bool); found {
		return val, err
	}
	return false, ErrInvalidType
}

// Filename returns the currently playing filename
func (c *Client) Filename() (string, error) {
	return c.GetProperty("filename")
}

// Path returns the currently playing path
func (c *Client) Path() (string, error) {
	return c.GetProperty("path")
}

// Pause returns true if the player is paused
func (c *Client) Pause() (bool, error) {
	return c.GetBoolProperty("pause")
}

// SetPause pauses or unpauses the player
func (c *Client) SetPause(pause bool) error {
	return c.SetProperty("pause", pause)
}

// Idle returns true if the player is idle
func (c *Client) Idle() (bool, error) {
	return c.GetBoolProperty("idle")
}

// Mute returns true if the player is muted.
func (c *Client) Mute() (bool, error) {
	return c.GetBoolProperty("mute")
}

// SetMute mutes or unmutes the player.
func (c *Client) SetMute(mute bool) error {
	return c.SetProperty("mute", mute)
}

// Fullscreen returns true if the player is in fullscreen mode.
func (c *Client) Fullscreen() (bool, error) {
	return c.GetBoolProperty("fullscreen")
}

// SetFullscreen activates/deactivates the fullscreen mode.
func (c *Client) SetFullscreen(v bool) error {
	return c.SetProperty("fullscreen", v)
}

// Volume returns the current volume level.
func (c *Client) Volume() (float64, error) {
	return c.GetFloatProperty("volume")
}

// Speed returns the current playback speed.
func (c *Client) Speed() (float64, error) {
	return c.GetFloatProperty("speed")
}

// Duration returns the duration of the currently playing file.
func (c *Client) Duration() (float64, error) {
	return c.GetFloatProperty("duration")
}

// Position returns the current playback position in seconds.
func (c *Client) Position() (float64, error) {
	return c.GetFloatProperty("time-pos")
}

// PercentPosition returns the current playback position in percent.
func (c *Client) PercentPosition() (float64, error) {
	return c.GetFloatProperty("percent-pos")
}
