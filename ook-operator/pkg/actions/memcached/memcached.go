
package memcached

import (
    "net/http"
    "os"
    "os/exec"

    "github.com/kupenstack/kupenstack/ook-operator/settings"
)



func Apply(w http.ResponseWriter, r *http.Request) {
    log := settings.Log.WithValues("action", "apply-memcached")

    cmd := exec.Command(settings.ActionsDir+"memcached/apply")
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if err != nil{
        log.Error(err, "Failed to apply changes")
        http.Error(w, http.StatusText(http.StatusInternalServerError),
                   http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}






