package dotnames

import (
	"fmt"
	"io"

	"github.com/serho/osrm-backend/integration/osrmfiles/osrmtype"
	"github.com/serho/osrm-backend/integration/osrmfiles/osrmtype/nametable"

	"github.com/serho/osrm-backend/integration/osrmfiles/fingerprint"
	"github.com/golang/glog"
)

// Contents represents `.osrm.names` file structure.
type Contents struct {
	Fingerprint fingerprint.Fingerprint
	nametable.IndexedData

	// for internal implementation
	writers  map[string]io.Writer
	filePath string
}

// New creates an empty Contents for `.osrm.names`.
func New(file string) *Contents {
	c := Contents{}

	c.filePath = file

	// init writers
	c.writers = map[string]io.Writer{}
	c.writers["osrm_fingerprint.meta"] = &c.Fingerprint
	c.writers["/common/names/blocks.meta"] = &c.IndexedData.BlocksMeta
	c.writers["/common/names/values.meta"] = &c.IndexedData.ValuesMeta
	c.writers["/common/names/blocks"] = &c.IndexedData.BlocksBuffer
	c.writers["/common/names/values"] = &c.IndexedData.ValuesBuffer

	return &c
}

// PrintSummary prints summary and head lines of contents.
func (c *Contents) PrintSummary(head int) {
	glog.Infof("Loaded from %s\n", c.filePath)
	glog.Infof("  %s\n", &c.Fingerprint)

	glog.Infof("  nametable.IndexedData blocks meta %d\n", c.IndexedData.BlocksMeta)
	glog.Infof("  nametable.IndexedData values meta %d\n", c.IndexedData.ValuesMeta)
	for i := 0; i < head; i++ {
		nameID := osrmtype.NameID(i * nametable.NameIDOffset)
		if s, err := c.IndexedData.GetNamesForID(nameID); err != nil {
			glog.Infof("    index %d, NameID %d: get names failed, err: %v", i, nameID, err)
		} else {
			glog.Infof("    index %d, NameID %d: %#v", i, nameID, s)
		}
	}
}

// Validate checks whether the contents valid or not.
func (c *Contents) Validate() error {
	if !c.Fingerprint.IsValid() {
		return fmt.Errorf("invalid fingerprint %v", c.Fingerprint)
	}
	if err := c.IndexedData.Validate(); err != nil {
		return err
	}

	return nil
}

// PostProcess post process the conents once contents loaded if necessary.
func (c *Contents) PostProcess() error {
	return c.IndexedData.Assemble()
}

// FindWriter find io.Writer for the specified name.
func (c *Contents) FindWriter(name string) (io.Writer, bool) {
	w, b := c.writers[name]
	return w, b
}

// FilePath returns the file path that stores the contents.
func (c *Contents) FilePath() string {
	return c.filePath
}
