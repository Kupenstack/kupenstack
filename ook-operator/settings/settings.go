
// Holds global variables. These can be overridden with ENV.
package settings

import (
	"github.com/go-logr/logr"
)


// Port to server at
var Port string

// Default configuration dir
var DefaultsDir string

// Final Configuration dir. Files in this are the pulled when any
// openstack-helm automation script/plugins are executed.
var ConfigDir string

// Directory containing all executables to automate openstack-helm
var ActionsDir string

var Log logr.Logger