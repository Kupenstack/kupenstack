/*
Copyright 2021 The Kupenstack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
)

// RecordEventf is a helper function to create an event for the given object.
func RecordEventf(eventRecorder record.EventRecorder, object metav1.Object, scheme *runtime.Scheme, eventtype, reason, messageFmt string, args ...interface{}) error {
	runtimeObject, ok := object.(runtime.Object)
	if !ok {
		return fmt.Errorf("object (%T) is not a runtime.Object", object)
	}

	objectRef, err := reference.GetReference(scheme, runtimeObject)
	if err != nil {
		return fmt.Errorf("Unable to get reference to owner")
	}

	eventRecorder.Eventf(objectRef, eventtype, reason, messageFmt, args...)

	return nil
}
