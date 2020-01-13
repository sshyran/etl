// Code generated for package schema by go-bindata DO NOT EDIT. (@generated)
// sources:
// descriptions/NDTResultRow.yaml
// descriptions/PTTest.yaml
// descriptions/README.md
// descriptions/TCPRow.yaml
// descriptions/toplevel.yaml
package schema

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _ndtresultrowYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x54\x4b\x6b\xe3\x3c\x14\xdd\xe7\x57\x1c\xba\xea\x07\xa9\x3f\x08\xcc\xa6\xbb\xe2\xc0\xd0\xa1\xc9\x84\xd8\x49\xd7\x8a\x74\x13\x0b\xf4\x30\x7a\xb8\xf4\xdf\x0f\x92\xed\xc9\x24\x6e\xda\x45\xb7\xf7\x5a\xe7\x71\xcf\xc1\x3f\x65\xa8\x1a\xeb\x42\x69\xb5\x96\xe1\x71\x06\x2c\xc9\x73\x27\xdb\x20\xad\x79\xc4\xe5\x1a\xd2\x23\x34\x94\xa6\xe0\xfd\xe4\xde\xa7\x35\x8e\xd6\xe9\xff\x60\x8f\x79\xed\xa2\x31\xd2\x9c\x66\x00\xe0\xc9\x75\xe4\xc0\xad\xa0\x62\xb6\x27\xe7\x13\xec\x35\xcb\x30\x1f\xe1\xfd\xbb\x3e\x58\x25\x39\xba\x61\x7e\x2f\x8f\x60\xe6\xfd\x9a\x60\x00\xcf\x3c\x3d\xc1\xac\xb4\x26\x38\xab\x26\x0c\x2b\x0a\x4c\xb0\xc0\x92\x50\xd4\xe5\x06\xdc\x1a\x43\x3c\x2d\x3d\x82\xcd\xa8\xeb\x65\xfd\x23\xcd\x13\x00\x78\xc3\x8c\x21\x55\xe0\x49\xa9\x4c\x90\xb7\x9a\x98\x8f\x8e\x34\x99\xe0\xd1\xb0\x8e\xc0\x26\x2f\x46\x0d\xc5\x6e\xf7\xbc\x9c\x08\x49\xc3\x2c\xe2\xe6\xb3\x8d\xb3\xc1\xf2\x0f\x3c\x8c\x0b\x44\x4f\x22\x63\x54\x8b\x12\xcc\x08\x94\x8b\xea\x42\x5a\x81\x3d\x53\x91\x3c\xa4\xe1\x2a\x0a\xc2\x6b\x35\xcf\x26\x5e\xab\x6a\x9e\x5f\x6c\x5e\x9e\x9e\xd7\x67\xce\x15\x79\xcf\x4e\x74\x93\xfa\xd9\x08\xd9\x49\x11\x99\x82\xee\x3f\xf5\x60\x8e\xe0\xc9\x04\xbc\xc9\xd0\xe4\x03\x5e\xa1\x8c\x2a\x32\xf3\xa8\xe4\x57\xf5\x7b\x3d\x47\xfd\xb2\x3f\x93\x97\x4a\x92\x09\x63\x42\x13\xee\x7e\xfd\xe0\xa8\xb5\x2e\x90\x80\x1e\xa3\x64\x1e\x86\x69\xfa\xbf\x4b\x2c\x68\x99\x74\x3e\x35\x60\x51\x7d\x9e\xfe\x80\x17\xec\x43\xd5\x77\xf3\x22\xd4\x96\x5c\xea\x32\x09\x44\x9f\x2a\x16\x1a\x3a\xc7\xdf\x0e\xce\xe6\xa0\xe2\x54\xe0\x6e\xd7\x2a\xcb\xc4\x5d\x91\x48\xbf\x88\x3b\x45\x74\xd9\xbb\xfe\xd5\x8a\x98\xa9\x1b\x67\xe3\xa9\x69\x63\x58\x1d\x5a\x3f\xc1\x18\x64\x72\xa6\x78\x54\x2c\x9d\x80\x75\xe4\xd8\x89\x10\x33\x3f\x1c\x0b\xa9\xfb\xd5\xa2\xfc\xdc\x79\x0f\x94\x9c\xf7\x37\xf8\x86\xf3\xa5\x7d\x33\x83\xf7\x6a\x51\x7e\xe1\x3d\xd5\xf4\xda\x7b\x7a\xf5\x2d\xef\x62\x50\x30\xb8\xcf\x78\xd2\x6c\xeb\x7a\x82\x51\x37\x04\x2d\x8d\xd4\x51\x63\x5b\xd7\xb0\x87\xfc\xdf\x10\x10\xd1\x0d\x46\xcf\x68\xff\x9c\xa4\x07\xed\x4f\xb5\x1d\xda\xf7\xa1\xc8\x44\xc0\xfb\x8b\xfe\x6d\xe9\x0d\x99\x7f\x02\x00\x00\xff\xff\x64\xad\xf8\x58\x6f\x05\x00\x00")

func ndtresultrowYamlBytes() ([]byte, error) {
	return bindataRead(
		_ndtresultrowYaml,
		"NDTResultRow.yaml",
	)
}

func ndtresultrowYaml() (*asset, error) {
	bytes, err := ndtresultrowYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "NDTResultRow.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pttestYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func pttestYamlBytes() ([]byte, error) {
	return bindataRead(
		_pttestYaml,
		"PTTest.yaml",
	)
}

func pttestYaml() (*asset, error) {
	bytes, err := pttestYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "PTTest.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _readmeMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\xce\xb1\x4e\xc4\x30\x0c\xc6\xf1\x3d\x4f\xf1\x49\xcc\xb4\xcf\x00\x3a\x31\x71\x03\x82\x85\xd1\x4d\xdc\xd6\x22\x89\x7b\xb1\x2b\xdd\xbd\x3d\x4a\xe1\xa6\xdb\xfd\xfd\xfe\x7e\xc2\x9b\x70\x4e\x38\xb1\xc5\x26\x9b\x8b\x56\x0b\xe1\x6b\x15\x43\x92\xc6\xd1\xb5\xdd\x10\xb5\x3a\x49\x35\xbc\xca\xf2\xb1\x73\xbb\x3d\x8e\xf0\xfd\x72\x7e\xc7\x2c\x99\x0d\xb3\x36\x58\x5c\xb9\x50\x48\x3c\x4b\x95\x43\x85\x54\x2c\xe2\xeb\x3e\x0d\x51\xcb\x58\x9e\x33\x4d\x23\x7b\x1e\xff\x4e\x87\x5e\xe5\x23\xc5\xd5\xa1\x33\x7c\x65\xe3\x7f\xd2\x56\xdd\x73\xc2\xc4\x1d\x2f\xe4\xce\x09\x64\x98\x2e\xd7\xe1\xf3\x98\x9f\x34\x1a\xe6\xa6\xa5\xcf\xc2\x43\x67\xd1\x71\xba\x5c\xb1\x51\xfc\xa1\x85\x87\x10\xce\x7b\x76\xd9\xf2\xdd\x8f\x54\x3b\x4e\xdb\x96\x85\x13\x5c\x3b\x03\xa3\xc2\xb8\xbf\xf7\x1b\x00\x00\xff\xff\x2c\x9a\x49\x9a\x2a\x01\x00\x00")

func readmeMdBytes() ([]byte, error) {
	return bindataRead(
		_readmeMd,
		"README.md",
	)
}

func readmeMd() (*asset, error) {
	bytes, err := readmeMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "README.md", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tcprowYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func tcprowYamlBytes() ([]byte, error) {
	return bindataRead(
		_tcprowYaml,
		"TCPRow.yaml",
	)
}

func tcprowYaml() (*asset, error) {
	bytes, err := tcprowYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "TCPRow.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _toplevelYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\xd0\x4d\x6a\xc3\x40\x0c\x05\xe0\xbd\x4f\xf1\x4e\xe0\x03\x64\xdb\xd2\x52\xe8\x1f\x6d\xda\x6d\x50\x6d\x39\x16\x19\x8f\x8c\xa4\x34\xd7\x2f\x33\x84\x60\x1a\x67\x67\x1e\x9e\xef\x89\xf7\x4e\xe6\xfc\x94\x07\xdd\x34\xc0\x3d\x7b\x67\x32\x87\x68\xde\xe0\x85\x83\x7a\x0a\x02\xfd\xe8\x31\x30\xea\x09\x31\x32\xe6\xf2\xc0\x30\x9b\x76\xec\xce\x7d\x0d\x27\x26\x3f\x1a\x4f\x9c\xa3\x6d\x2e\x64\xbb\x25\x3f\x3c\x48\xe2\x57\x9a\xf8\xca\x7f\xbc\xfb\xc4\xd7\xc7\x33\x42\x2b\x41\xd6\x8d\xf2\xcb\xe8\x34\x07\x49\x96\xbc\xaf\x71\xb0\xc7\x4e\x7a\x0c\x6a\x88\x51\x1c\xa6\xa7\x65\x45\xfd\xda\xca\x8a\x5f\x42\xc4\x48\x71\xeb\xec\x1b\x98\x7d\xb3\x79\x11\xfe\x83\xe7\x1c\x3a\x2c\xc5\xda\xb0\xc6\x9e\x2f\xbf\x62\xde\x4c\xf6\x92\x29\x61\x90\xc4\x99\x26\x2e\xe0\x62\x40\x90\xe3\x64\x12\xc1\xb9\x6c\xd3\x8b\x1f\x40\xb9\x87\xe4\x52\xdb\x00\xa8\xd3\x9d\xf7\x6a\x9b\xa4\xfb\x5d\xac\x0d\x70\x29\x5a\xe2\x9d\xa6\xc4\x5d\xf9\x01\xe5\x91\x07\x4d\x73\xdb\xfc\x05\x00\x00\xff\xff\xde\xad\xbd\x13\x06\x02\x00\x00")

func toplevelYamlBytes() ([]byte, error) {
	return bindataRead(
		_toplevelYaml,
		"toplevel.yaml",
	)
}

func toplevelYaml() (*asset, error) {
	bytes, err := toplevelYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "toplevel.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"NDTResultRow.yaml": ndtresultrowYaml,
	"PTTest.yaml":       pttestYaml,
	"README.md":         readmeMd,
	"TCPRow.yaml":       tcprowYaml,
	"toplevel.yaml":     toplevelYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"NDTResultRow.yaml": &bintree{ndtresultrowYaml, map[string]*bintree{}},
	"PTTest.yaml":       &bintree{pttestYaml, map[string]*bintree{}},
	"README.md":         &bintree{readmeMd, map[string]*bintree{}},
	"TCPRow.yaml":       &bintree{tcprowYaml, map[string]*bintree{}},
	"toplevel.yaml":     &bintree{toplevelYaml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
