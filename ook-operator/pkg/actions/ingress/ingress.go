
package ingress

import (
    "net/http"
    "os"
    "os/exec"

    "github.com/kupenstack/kupenstack/ook-operator/settings"
    pkg "github.com/kupenstack/kupenstack/ook-operator/pkg/actions"
)



func Apply(w http.ResponseWriter, r *http.Request) {
    log := settings.Log.WithValues("action", "apply-ingress")

    err := pkg.PrepareOOKValues(r, []string{"ingressKubeSystem.yaml", "ingressComponent.yaml"})
    if err != nil {
        log.Error(err, "Failed to prepare OOK chart values.")
        http.Error(w, http.StatusText(http.StatusInternalServerError),
                   http.StatusInternalServerError)
        return
    }

    cmd := exec.Command(settings.ActionsDir+"ingress/apply")
    cmd.Stdout = os.Stdout
    err = cmd.Start()
    if err != nil{
        log.Error(err, "Failed to apply changes")
        http.Error(w, http.StatusText(http.StatusInternalServerError),
                   http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}






